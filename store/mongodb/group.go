package mongodb

import (
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/helpers"
)

func (db *mongo) CreateGroup(group *models.Group) error {
	session := db.Session.Copy()
	defer session.Close()
	groups := db.C(models.GroupsCollection).With(session)

	group.Id = bson.NewObjectId().Hex()
	err := groups.Insert(group)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the group in the database", err)
	}

	return nil
}