package config

import "github.com/layer5io/meshkit/errors"

const (

	// ErrGetLatestReleasesCode represents the error which occurs during the process of getting
	// latest releases
	ErrGetLatestReleasesCode = "1008"
	// ErrGetManifestNamesCode represents the error which occurs during the process of getting
	// manifest names
	ErrGetManifestNamesCode = "1009"
	// ErrEmptyConfigCode represents the error when the configuration is either empty
	// or is invalid
	ErrEmptyConfigCode = "1010"
)

var (
	// ErrEmptyConfig error is the error when config is invalid
	ErrEmptyConfig = errors.New(ErrEmptyConfigCode, errors.Alert, []string{"Config is empty"}, []string{}, []string{}, []string{})
)

// ErrGetLatestReleases is the error for fetching consul releases
func ErrGetLatestReleases(err error) error {
	return errors.New(ErrGetLatestReleasesCode, errors.Alert, []string{"Unable to fetch release info"}, []string{err.Error()}, []string{}, []string{})
}

// ErrGetManifestNames is the error for fetching consul manifest names
func ErrGetManifestNames(err error) error {
	return errors.New(ErrGetManifestNamesCode, errors.Alert, []string{"Unable to fetch manifest names from github"}, []string{err.Error()}, []string{}, []string{})
}
