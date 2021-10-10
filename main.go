// Copyright 2019 Layer5.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/layer5io/meshery-adapter-library/adapter"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-consul/consul/oam"
	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/layer5io/meshery-consul/internal/operations"

	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-consul/consul"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	serviceName = "consul-adapter"
	version     = "none"
	gitsha      = "none"
)

func main() {
	log, err := logger.New(serviceName, logger.Options{Format: logger.SyslogLogFormat, DebugLevel: isDebug()})
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	cfg, err := config.New(configprovider.Options{
		ServerConfig:   config.ServerDefaults,
		MeshSpec:       config.MeshSpecDefaults,
		ProviderConfig: config.ViperDefaults,
		Operations:     operations.Operations,
	},
	)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	_ = cfg.GetObject(adapter.ServerKey, &service)

	kubeCfg, err := config.New(configprovider.Options{
		ProviderConfig: config.KubeConfigDefaults,
	})

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// KUBECONFIG is required by the Helm library, and other tools that might be used by the adapter in the future.
	kubeconfig := path.Join(
		config.KubeConfigDefaults[configprovider.FilePath],
		fmt.Sprintf("%s.%s", config.KubeConfigDefaults[configprovider.FileName], config.KubeConfigDefaults[configprovider.FileType]))
	err = os.Setenv("KUBECONFIG", kubeconfig)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("KUBECONFIG: %s", kubeconfig))

	service.Handler = consul.New(cfg, log, kubeCfg)
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()
	service.Version = version
	service.GitSHA = gitsha

	go registerCapabilities(service.Port, log)        //Registering static capabilities
	go registerDynamicCapabilities(service.Port, log) //Registering latest capabilities periodically

	// Server Initialization
	log.Info("Adaptor Listening at port: ", service.Port)
	err = grpc.Start(service, nil)
	if err != nil {
		log.Error(grpc.ErrGrpcServer(err))
		os.Exit(1)
	}
}

func registerCapabilities(port string, log logger.Handler) {
	log.Info("Registering static workloads...")
	// Register workloads
	if err := oam.RegisterWorkloads(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Error(err)
	}

	// Register traits
	if err := oam.RegisterTraits(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Error(err)
	}
}

func registerDynamicCapabilities(port string, log logger.Handler) {
	registerWorkloads(port, log)
	//Start the ticker
	const reRegisterAfter = 24
	ticker := time.NewTicker(reRegisterAfter * time.Hour)
	for {
		<-ticker.C
		registerWorkloads(port, log)
	}
}

func registerWorkloads(port string, log logger.Handler) {
	crds, err := config.GetFileNames("https://api.github.com/repos/hashicorp/consul-k8s", "control-plane/config/crd/bases")
	if err != nil {
		log.Error(err)
		return
	}
	rel, err := config.GetLatestReleases(1)
	if err != nil {
		log.Info("Could not get latest version ", err.Error())
		return
	}
	appVersion := rel[0].TagName
	log.Info("Registering latest workload components for version ", appVersion)
	// Register workloads
	for _, manifest := range crds {
		log.Info("Registering for ", manifest)
		if err := adapter.RegisterWorkLoadsDynamically(mesheryServerAddress(), serviceAddress()+":"+port, &adapter.DynamicComponentsConfig{
			TimeoutInMinutes: 60,
			URL:              "https://raw.githubusercontent.com/hashicorp/consul-k8s/main/control-plane/config/crd/bases/" + manifest,
			GenerationMethod: adapter.Manifests,
			Config: manifests.Config{
				Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_CONSUL)],
				MeshVersion: appVersion,
				Filter: manifests.CrdFilter{
					RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
					NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
					VersionFilter: []string{"$[0]..spec.versions[0]"},
					GroupFilter:   []string{"$[0]..spec"},
					SpecFilter:    []string{"$[0]..openAPIV3Schema.properties.spec"},
					ItrFilter:     []string{"$[?(@.spec.names.kind"},
					ItrSpecFilter: []string{"$[?(@.spec.names.kind"},
					VField:        "name",
					GField:        "group",
				},
			},
			Operation: config.ConsulOperation,
		}); err != nil {
			log.Error(err)
			return
		}
		log.Info(manifest, " registered")
	}
	log.Info("Latest workload components successfully registered.")
}

func mesheryServerAddress() string {
	meshReg := os.Getenv("MESHERY_SERVER")

	if meshReg != "" {
		if strings.HasPrefix(meshReg, "http") {
			return meshReg
		}

		return "http://" + meshReg
	}

	return "http://localhost:9081"
}

func serviceAddress() string {
	svcAddr := os.Getenv("SERVICE_ADDR")

	if svcAddr != "" {
		return svcAddr
	}

	return "mesherylocal.layer5.io"
}
func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
