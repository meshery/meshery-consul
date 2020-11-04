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

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-consul/internal/operations"

	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-consul/consul"
	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
)

const (
	serviceName = "consul-adapter"
)

var (
	kubeConfigPath = fmt.Sprintf("%s/.kube/meshery.config", utils.GetHome())
)

func main() {
	log, err := logger.New(serviceName, logger.Options{Format: logger.JsonLogFormat, DebugLevel: false})
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	cfg, err := config.New(provider.Options{
		ServerConfig:   config.ServerDefaults,
		MeshSpec:       config.MeshSpecDefaults,
		MeshInstance:   config.MeshInstanceDefaults,
		ProviderConfig: config.ViperDefaults,
		Operations:     operations.Operations,
	},
	)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	cfg.SetKey(adapter.KubeconfigPathKey, kubeConfigPath)

	service := &grpc.Service{}
	_ = cfg.GetObject(adapter.ServerKey, &service)

	service.Handler = consul.New(cfg, log)
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()
	err = grpc.Start(service, nil)
	if err != nil {
		log.Error(grpc.ErrGrpcServer(err))
		os.Exit(1)
	}
}
