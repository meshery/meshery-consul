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
	"path/filepath"
	"strings"
	"time"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-consul/consul/oam"
	"github.com/layer5io/meshery-consul/internal/config"
	configprovider "github.com/layer5io/meshkit/config/provider"

	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-consul/build"
	"github.com/layer5io/meshery-consul/consul"
	"github.com/layer5io/meshkit/logger"
)

var (
	serviceName = "consul-adapter"
	version     = "edge"
	gitsha      = "none"
)

func main() {
	log, err := logger.New(serviceName, logger.Options{Format: logger.SyslogLogFormat, DebugLevel: isDebug()})
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	cfg, err := config.New(configprovider.ViperKey)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	_ = cfg.GetObject(adapter.ServerKey, &service)

	kubeCfg, err := config.NewKubeconfigBuilder(configprovider.ViperKey)

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
	log.Info("Registering latest workload components for version ", version)
	version := build.LatestVersion
	gm := build.DefaultGenerationMethod

	// Prechecking to skip comp gen
	if os.Getenv("FORCE_DYNAMIC_REG") != "true" && oam.AvailableVersions[build.LatestAppVersion] {
		log.Info("Components available statically for version ", version, ". Skipping dynamic component registeration")
		return
	}
	//If a URL is passed from env variable, it will be used for component generation with default method being "using manifests"
	// In case a helm chart URL is passed, COMP_GEN_METHOD env variable should be set to Helm otherwise the component generation fails
	if os.Getenv("COMP_GEN_URL") != "" && (os.Getenv("COMP_GEN_METHOD") == "Helm" || os.Getenv("COMP_GEN_METHOD") == "Manifest") {
		build.OverrideURL = os.Getenv("COMP_GEN_URL")
		gm = os.Getenv("COMP_GEN_METHOD")
		log.Info("Registering workload components from url ", build.OverrideURL, " using ", gm, " method...")
		build.CRDnames = []string{"user passed configuration"}
	}
	// Register workloads

	for _, manifest := range build.CRDnames {
		log.Info("Registering for ", manifest)
		if err := adapter.CreateComponents(adapter.StaticCompConfig{
			URL:     build.GetDefaultURL(manifest, build.LatestVersion),
			Method:  gm,
			Path:    build.WorkloadPath,
			DirName: build.LatestAppVersion,
			Config:  build.NewConfig(build.LatestAppVersion),
		}); err != nil {
			log.Error(err)
			continue
		}
		log.Info(manifest, " registered")
	}

	//*The below log is checked in the workflows. If you change this log, reflect that change in the workflow where components are generated
	log.Info("Component creation completed for version ", version)

	//Now we will register in case
	log.Info("Registering workloads with Meshery Server for version ", version)
	originalPath := oam.WorkloadPath
	oam.WorkloadPath = filepath.Join(originalPath, version)
	defer resetWorkloadPath(originalPath)
	if err := oam.RegisterWorkloads(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
		return
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
func resetWorkloadPath(orig string) {
	oam.WorkloadPath = orig
}
