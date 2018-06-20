package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
)

func CreateSigfoxMessage(c context.Context, message *sigfox.Message) error {
	return FromContext(c).CreateSigfoxMessage(message)
}
