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

package config

import (
	"path"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/utils"
)

const (
	// OAM Metadata constants
	OAMAdapterNameMetadataKey       = "adapter.meshery.io/name"
	OAMComponentCategoryMetadataKey = "ui.meshery.io/category"
)

var (
	ConfigRootPath = path.Join(utils.GetHome(), ".meshery")
)

// New creates a new config instance
func New(provider string) (h config.Handler, err error) {
	opts := configprovider.Options{
		FilePath: ConfigRootPath,
		FileName: "consul",
		FileType: "yaml",
	}
	// Config provider
	switch provider {
	case configprovider.ViperKey:
		h, err = configprovider.NewViper(opts)
		if err != nil {
			return nil, err
		}
	case configprovider.InMemKey:
		h, err = configprovider.NewInMem(opts)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrEmptyConfig
	}

	// Setup Server config
	if err := h.SetObject(adapter.ServerKey, ServerDefaults); err != nil {
		return nil, err
	}

	// setup Mesh config
	if err := h.SetObject(adapter.MeshSpecKey, MeshSpecDefaults); err != nil {
		return nil, err
	}

	// setup Operation Config
	if err := h.SetObject(adapter.OperationsKey, Operations); err != nil {
		return nil, err
	}

	return h, nil
}

// NewKubeconfigBuilder returns a config handler based on the provider
//
// Valid providers are "viper" and "in-mem"
func NewKubeconfigBuilder(provider string) (config.Handler, error) {
	opts := configprovider.Options{
		FilePath: ConfigRootPath,
		FileType: "yaml",
		FileName: "kubeconfig",
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, ErrEmptyConfig
}

// RootPath returns the config root path for the adapter
func RootPath() string {
	return ConfigRootPath
}
