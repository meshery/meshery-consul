package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
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

// GetLatestReleases fetches the latest releases from the traefik mesh repository
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

	body, err := ioutil.ReadAll(resp.Body)
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

type rootdir struct {
	Commit commit `json:"commit"`
}
type commit struct {
	Tree roottree `json:"tree"`
}
type roottree struct {
	URL string `json:"url"`
}
type treeslice struct {
	Tree []map[string]interface{} `json:"tree"`
}

// GetFileNames takes the url of a github repo and the path to a directory. Then returns all the filenames from that directory
func GetFileNames(url string, path string) ([]string, error) {
	res, err := http.Get(url + "/commits") //url="https://api.github.com/repos/hashicorp/consul-k8s"
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r []rootdir
	err = json.Unmarshal(content, &r)
	if err != nil {
		return nil, err
	}
	shaurl := r[0].Commit.Tree.URL
	bpath := strings.Split(path, "/")
	ans, err := _getname(shaurl, bpath)
	if err == nil {
		return ans, nil
	}
	return ans, ErrGetManifestNames(err)
}

func _getname(shaURL string, bpath []string) ([]string, error) {
	var ans []string
	res, err := http.Get(shaURL)
	if err != nil {
		return []string{}, err
	}
	defer res.Body.Close()
	var t treeslice
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}
	err = json.Unmarshal(content, &t)
	if err != nil {
		return nil, err
	}
	if len(bpath) != 0 {
		dirName := bpath[0]
		bpath = bpath[1:]

		for _, c := range t.Tree {
			var (
				url string
				ok  bool
			)
			if url, ok = c["url"].(string); !ok {
				return nil, errors.New("Invalid URL field")
			}
			if c["path"] == dirName {
				tempans, err := _getname(url, bpath)
				return tempans, err
			}
		}
		return []string{}, err
	}

	//base case
	for _, c := range t.Tree {
		var (
			path string
			ok   bool
		)
		if path, ok = c["path"].(string); !ok {
			return nil, errors.New("Invalid path field")
		}
		ans = append(ans, path)
	}
	return ans, nil
}
