package config

import "github.com/layer5io/meshkit/errors"

const (

	// ErrGetLatestReleasesCode represents the error which occurs during the process of getting
	// latest releases
	ErrGetLatestReleasesCode = "blah1"
)

// ErrGetLatestReleases is the error for fetching nsm-mesh releases
func ErrGetLatestReleases(err error) error {
	return errors.New(ErrGetLatestReleasesCode, errors.Alert, []string{"Unable to fetch release info"}, []string{err.Error()}, []string{}, []string{})
}
