package oam

import "github.com/layer5io/meshkit/errors"

const (
	ErrLoadingPathsetCode = "1011"
)

func ErrLoadingPathset(err error) error {
	return errors.New(ErrLoadingPathsetCode, errors.Alert, []string{"Could not create a pathset for static component generation"}, []string{err.Error()}, []string{}, []string{})
}
