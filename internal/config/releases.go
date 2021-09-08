package config

import (
	"encoding/json"
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
	Url string `json:"url"`
}
type lol struct {
	Tree []map[string]string `json:"tree"`
}

// GetFileNames takes the url of a github repo and the path to a directory. Then returns all the filenames from that directory
func GetFileNames(url string, path string) []string {
	res, err := http.Get(url + "/commits") //url="https://api.github.com/repos/hashicorp/consul-k8s"
	if err != nil {
		return []string{}
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	var r []rootdir
	json.Unmarshal(content, &r)
	shaurl := r[0].Commit.Tree.Url
	bpath := strings.Split(path, "/")
	ans := _getname(shaurl, bpath)
	return ans
}

func _getname(shaUrl string, bpath []string) []string {
	var ans []string
	res, err := http.Get(shaUrl)
	if err != nil {
		return []string{}
	}
	defer res.Body.Close()
	var t lol
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}
	}
	json.Unmarshal(content, &t)
	if len(bpath) != 0 {
		dirName := bpath[0]
		bpath = bpath[1:]

		for _, c := range t.Tree {
			if c["path"] == dirName {
				tempans := _getname(c["url"], bpath)
				return tempans
			}
		}
		return []string{}
	}

	//base case
	for _, c := range t.Tree {
		ans = append(ans, c["path"])
	}
	return ans
}
