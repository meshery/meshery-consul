package config

import "github.com/layer5io/meshkit/errors"

const (

	// ErrGetLatestReleasesCode represents the error which occurs during the process of getting
	// latest releases
	ErrGetLatestReleasesCode = "blah1"
	// ErrGetManifestNamesCode represents the error which occurs during the process of getting
	// manifest names
	ErrGetManifestNamesCode = "blah2"
)

// ErrGetLatestReleases is the error for fetching consul releases
func ErrGetLatestReleases(err error) error {
	return errors.New(ErrGetLatestReleasesCode, errors.Alert, []string{"Unable to fetch release info"}, []string{err.Error()}, []string{}, []string{})
}

// ErrGetManifestNames is the error for fetching consul manifest names
func ErrGetManifestNames(err error) error {
	return errors.New(ErrGetManifestNamesCode, errors.Alert, []string{"Unable to fetch manifest names from github"}, []string{err.Error()}, []string{}, []string{})
}
