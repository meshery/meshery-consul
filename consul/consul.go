// Copyright 2020 Layer5, Inc.
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

package consul

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models"
	"github.com/layer5io/meshkit/utils/events"
	"gopkg.in/yaml.v2"
)

type Consul struct {
	adapter.Adapter
}

func New(config config.Handler, log logger.Handler, kubeConfig config.Handler, e *events.EventStreamer) adapter.Handler {
	return &Consul{
		adapter.Adapter{Config: config, Log: log, KubeconfigHandler: kubeConfig, EventStreamer: e},
	}
}

// CreateKubeconfigs creates and writes passed kubeconfig onto the filesystem
func (h *Consul) CreateKubeconfigs(kubeconfigs []string) error {
	var errs = make([]error, 0)
	for _, kubeconfig := range kubeconfigs {
		kconfig := models.Kubeconfig{}
		err := yaml.Unmarshal([]byte(kubeconfig), &kconfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// To have control over what exactly to take in on kubeconfig
		h.KubeconfigHandler.SetKey("kind", kconfig.Kind)
		h.KubeconfigHandler.SetKey("apiVersion", kconfig.APIVersion)
		h.KubeconfigHandler.SetKey("current-context", kconfig.CurrentContext)
		err = h.KubeconfigHandler.SetObject("preferences", kconfig.Preferences)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = h.KubeconfigHandler.SetObject("clusters", kconfig.Clusters)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = h.KubeconfigHandler.SetObject("users", kconfig.Users)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = h.KubeconfigHandler.SetObject("contexts", kconfig.Contexts)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return mergeErrors(errs)
}
