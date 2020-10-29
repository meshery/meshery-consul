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
	"time"

	"github.com/layer5io/gokit/errors"
	"github.com/layer5io/gokit/logger"
	"github.com/layer5io/gokit/utils"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-consul/consul"
	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/layer5io/meshery-consul/internal/operations"
)

const (
	serviceName    = "consul-adapter"
	configProvider = "viper"
)

var (
	kubeConfigPath = fmt.Sprintf("%s/.kube/meshery.config", utils.GetHome())
)

func main() {
	log, err := logger.New(serviceName)
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	cfg, err := config.New(configProvider, config.ServerDefaults, config.MeshSpecDefaults, config.MeshInstanceDefaults, config.ViperDefaults, operations.Operations)
	if err != nil {
		log.Err("Config Init Failed", err.Error())
		os.Exit(1)
	}
	cfg.SetKey("kube-config-path", kubeConfigPath)

	service := &grpc.Service{}
	_ = cfg.Server(&service)

	service.Handler = consul.New(cfg, log)
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()
	log.Info(fmt.Sprintf("%s starting on port: %s", service.Name, service.Port))
	err = grpc.Start(service, nil)
	if err != nil {
		log.Err(errors.ErrGrpcServer, "adapter crashed on startup")
		os.Exit(1)
	}
}
