package build

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-consul/internal/config"

	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

type Versions struct {
	AppVersion   string
	ChartVersion string
}

var DefaultGenerationMethod string
var LatestVersion string
var LatestAppVersion string
var WorkloadPath string
var MeshModelPath string
var CRDnames []string
var OverrideURL string
var AllVersions []Versions

var meshmodelmetadata = make(map[string]interface{})

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Orchestration & Management",
	Metadata: meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_CONSUL)],
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func GetDefaultURL(crd string, version string) string {
	if OverrideURL != "" {
		return OverrideURL
	}
	return strings.Join([]string{fmt.Sprintf("https://raw.githubusercontent.com/hashicorp/consul-k8s/%s/control-plane/config/crd/bases", version), crd}, "/")
}
func init() {
	wd, _ := os.Getwd()
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &meshmodelmetadata)
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	allVersions, _ := utils.GetLatestReleaseTagsSorted("hashicorp", "consul-k8s")
	if len(allVersions) == 0 {
		return
	}
	for i, v := range allVersions {
		if i == len(allVersions)-1 { //only get AppVersion of latest chart version
			//Executing the below function for all versions is redundant and takes time on startup, we only want to know the latest app version of latest version.
			av, err := kubernetes.HelmChartVersionToAppVersion("https://helm.releases.hashicorp.com", "consul", strings.TrimPrefix(v, "v"))
			if err != nil {
				log.Println("could not find app version for " + v + err.Error())
			}
			AllVersions = append(AllVersions, Versions{
				ChartVersion: v,
				AppVersion:   av,
			})
		} else {
			AllVersions = append(AllVersions, Versions{
				ChartVersion: v,
			})
		}
	}
	CRDnames, _ = config.GetFileNames("hashicorp", "consul-k8s", "control-plane/config/crd/bases/")
	LatestAppVersion = AllVersions[len(AllVersions)-1].AppVersion
	LatestVersion = AllVersions[len(AllVersions)-1].ChartVersion
	DefaultGenerationMethod = adapter.Manifests
}
