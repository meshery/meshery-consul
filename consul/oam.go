package consul

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
)

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Consul, v1alpha1.Component, bool) (string, error)

func (h *Consul) HandleComponents(comps []v1alpha1.Component, isDel bool) (string, error) {
	var errs []error
	var msgs []string

	compFuncMap := map[string]CompHandler{
		"ConsulMesh": handleComponentConsulMesh,
	}
	for _, comp := range comps {
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleConsulCoreComponents(h, comp, isDel, "", "")
			if err != nil {
				errs = append(errs, err)
				continue
			}

			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(h, comp, isDel)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		msgs = append(msgs, msg)
	}
	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}
func (h *Consul) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Errorf(strings.Join(errMsgs, "\n"))
}

func mergeMsgs(strs []string) string {
	return strings.Join(strs, "\n")
}

func handleComponentConsulMesh(c *Consul, comp v1alpha1.Component, isDelete bool) (string, error) {
	// Get the consul version from the settings
	// we are sure that the version of consul would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Settings["version"].(string)

	msg, err := c.installConsul(isDelete, version, comp.Namespace)
	if err != nil {
		return fmt.Sprintf("%s: %s", comp.Name, msg), err
	}

	return fmt.Sprintf("%s: %s", comp.Name, msg), nil
}
func handleConsulCoreComponents(
	c *Consul,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string) (string, error) {
	if apiVersion == "" {
		apiVersion = getAPIVersionFromComponent(comp)
		if apiVersion == "" {
			return "", ErrConsulCoreComponentFail(fmt.Errorf("failed to get API Version for: %s", comp.Name))
		}
	}

	if kind == "" {
		kind = getKindFromComponent(comp)
		if kind == "" {
			return "", ErrConsulCoreComponentFail(fmt.Errorf("failed to get kind for: %s", comp.Name))
		}
	}
	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
		},
		"spec": comp.Spec.Settings,
	}

	// Convert to yaml
	yamlByt, err := yaml.Marshal(component)
	if err != nil {
		err = ErrParseConsulCoreComponent(err)
		c.Log.Error(err)
		return "", err
	}

	msg := fmt.Sprintf("created %s \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	if isDel {
		msg = fmt.Sprintf("deleted %s config \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	}

	return msg, c.MesheryKubeclient.ApplyManifest(yamlByt, mesherykube.ApplyOptions{
		Namespace: comp.Namespace,
		Update:    true,
		Delete:    isDel,
	})
}
func getAPIVersionFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sAPIVersion"]
}
func getKindFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sKind"]
}
