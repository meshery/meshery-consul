package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshkit/utils/walker"
)

// Release is used to save the release informations
type Release struct {
	ID      int             `json:"id,omitempty"`
	TagName string          `json:"tag_name,omitempty"`
	Name    adapter.Version `json:"name,omitempty"`
	Draft   bool            `json:"draft,omitempty"`
	Assets  []*Asset        `json:"assets,omitempty"`
}

// Asset describes the github release asset object
type Asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

// GetLatestReleases fetches the latest releases from the Consul mesh repository
func GetLatestReleases(releases uint) ([]*Release, error) {
	releaseAPIURL := "https://api.github.com/repos/hashicorp/consul-k8s/releases?per_page" + fmt.Sprint(releases)
	// We need a variable url here hence using nosec
	// #nosec
	resp, err := http.Get(releaseAPIURL)
	if err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	if resp.StatusCode != http.StatusOK {
		return []*Release{}, ErrGetLatestReleases(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	var releaseList []*Release

	if err = json.Unmarshal(body, &releaseList); err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	if err = resp.Body.Close(); err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	return releaseList, nil
}

// GetFileNames takes the url of a github repo and the path to a directory. Then returns all the filenames from that directory
func GetFileNames(owner string, repo string, path string) ([]string, error) {
	g := walker.NewGit()
	var filenames []string
	err := g.Owner(owner).Repo(repo).Branch("master").Root(path).RegisterFileInterceptor(func(f walker.File) error {
		filenames = append(filenames, f.Name)
		return nil
	}).Walk()
	return filenames, err
}
