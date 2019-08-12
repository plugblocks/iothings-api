package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
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

func (db *mongo) GetAllGroups(user *models.User) ([]models.Group, error) {
	session := db.Session.Copy()
	defer session.Close()
	groupCollection := db.C(models.GroupsCollection).With(session)

	groups := []models.Group{}
	err := groupCollection.Find(bson.M{"user_id": user.Id}).All(&groups)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_groups_failed", "Failed to get the groups: "+err.Error(), err)
	}

	return groups, nil
}

func (db *mongo) GetGroupById(user *models.User, id string) (*models.Group, error) {
	session := db.Session.Copy()
	defer session.Close()
	groupCollection := db.C(models.GroupsCollection).With(session)

	group := &models.Group{}
	err := groupCollection.Find(bson.M{"_id": id, "user_id": user.Id}).One(group)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "group_not_found", "Could not find the group", err)
	}

	return group, nil
}

func (db *mongo) UpdateGroup(user *models.User, id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	groups := db.C(models.GroupsCollection).With(session)

	err := groups.Update(bson.M{"_id": id, "user_id": user.Id}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "group_update_failed", "Failed to update the groups: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteGroup(user *models.User, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	groups := db.C(models.GroupsCollection).With(session)

	err := groups.Remove(bson.M{"_id": id, "user_id": user.Id})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "group_delete_failed", "Failed to delete the group", err)
	}

	return nil
}

func (db *mongo) CountGroups() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	groups := db.C(models.GroupsCollection).With(session)

	nbr, err := groups.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "hroups_not_found", "Groups not found", err)
	}
	return nbr, nil
}
