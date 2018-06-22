package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateGroup(user *models.User, group *models.Group) error {
	session := db.Session.Copy()
	defer session.Close()
	groups := db.C(models.GroupsCollection).With(session)

	group.Id = bson.NewObjectId().Hex()
	group.UserId = user.Id

	err := groups.Insert(group)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "group_creation_failed", "Failed to insert the group in the database", err)
	}

	return nil
}
