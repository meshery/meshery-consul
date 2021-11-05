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
	"github.com/layer5io/meshkit/errors"
)

var (
	ErrApplyOperationCode           = "1000"
	ErrParseOAMComponentCode        = "1001"
	ErrParseOAMConfigCode           = "1002"
	ErrProcessOAMCode               = "1003"
	ErrApplyHelmChartCode           = "1004"
	ErrMeshConfigCode               = "1005"
	ErrConsulCoreComponentFailCode  = "1006"
	ErrParseConsulCoreComponentCode = "1007"
	ErrParseOAMComponent            = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occurred while prasing application component in the OAM request made"}, []string{"Invalid OAM component passed in OAM request"}, []string{"Check if your request has vaild OAM components"})
	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occurred while prasing component config in the OAM request made"}, []string{"Invalid OAM config passed in OAM request"}, []string{"Check if your request has vaild OAM config"})
)

func ErrApplyOperation(err error) error {
	return errors.New(ErrApplyOperationCode, errors.Alert, []string{"Error applying operation", err.Error()}, []string{}, []string{}, []string{})
}

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyHelmChart is the error which occurs in the process of applying helm chart
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"Error applying helm chart"}, []string{err.Error()}, []string{"Chart could be invalid"}, []string{"Use `helm verify` and `helm lint` to verify chart path and validity"})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error(), "Error getting MeshSpecKey config from in-memory configuration"}, []string{}, []string{})
}

// ErrConsulCoreComponentFail is the error when core Consul component processing fails
func ErrConsulCoreComponentFail(err error) error {
	return errors.New(ErrConsulCoreComponentFailCode, errors.Alert, []string{"error in Consul core component"}, []string{err.Error()}, []string{"API version or Kind passed is empty"}, []string{"Make sure API version and Kind are not empty"})
}

// ErrParseConsulCoreComponent is the error when Consul core component manifest parsing fails
func ErrParseConsulCoreComponent(err error) error {
	return errors.New(ErrParseConsulCoreComponentCode, errors.Alert, []string{"Consul core component manifest parsing failing"}, []string{err.Error()}, []string{"Could not marshall generated component to YAML"}, []string{})
}
