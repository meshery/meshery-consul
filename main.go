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
	"time"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-consul/internal/config"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/utils/events"

	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-consul/consul"
	"github.com/layer5io/meshkit/logger"
)

var (
	serviceName = "consul-adapter"
	version     = "edge"
	gitsha      = "none"
	isDebug     = os.Getenv("DEBUG") == "true"
)

func init() {
	// Create the config path if it doesn't exists as the entire adapter
	// expects that directory to exists, which may or may not be true
	if err := os.MkdirAll(path.Join(config.RootPath(), "bin"), 0750); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main() {
	log, err := logger.New(serviceName, logger.Options{Format: logger.SyslogLogFormat, DebugLevel: isDebug})
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
	e := events.NewEventStreamer()
	service.Handler = consul.New(cfg, log, kubeCfg, e)
	service.EventStreamer = e
	service.StartedAt = time.Now()
	service.Version = version
	service.GitSHA = gitsha

	// Server Initialization
	log.Info("Adaptor Listening at port: ", service.Port)
	err = grpc.Start(service)
	if err != nil {
		log.Error(grpc.ErrGrpcServer(err))
		os.Exit(1)
	}
}
