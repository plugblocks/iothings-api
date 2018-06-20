package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
)

func CreateSigfoxLocation(c context.Context, location *sigfox.Location) error {
	return FromContext(c).CreateSigfoxLocation(location)
}
