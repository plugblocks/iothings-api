package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
)

func CreateSigfoxMessage(c context.Context, message *sigfox.Message) error {
	return FromContext(c).CreateSigfoxMessage(message)
}

func CreateSigfoxLocation(c context.Context, location *sigfox.Location) error {
	return FromContext(c).CreateSigfoxLocation(location)
}
