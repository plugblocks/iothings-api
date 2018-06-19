package mongodb

import (
	"net/http"

	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateUser(user *models.User) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user.Id = bson.NewObjectId().Hex()
	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	if count, _ := users.Find(bson.M{"email": user.Email}).Count(); count > 0 {
		return helpers.NewError(http.StatusConflict, "user_already_exists", "User already exists", err)
	}

	err = users.Insert(user)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the user in the database", err)
	}

	return nil
}

func (db *mongo) FindUserById(id string) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}
	err := users.FindId(id).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

func (db *mongo) FindUser(params params.M) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}

	err := users.Find(params).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

func (db *mongo) ActivateUser(activationKey string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Update(bson.M{"$and": []bson.M{{"_id": id}, {"activationKey": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to activate", err)
	}
	return nil
}

func (db *mongo) UpdateUser(user *models.User, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.UpdateId(user.Id, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user", err)
	}

	return nil
}

func (db *mongo) GetUsers() ([]*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()

	users := db.C(models.UsersCollection).With(session)

	list := []*models.User{}
	err := users.Find(params.M{}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "users_not_found", "Users not found", err)
	}

	return list, nil
}

func (db *mongo) UserIsAdmin(userId string) (bool, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}
	err := users.FindId(userId).One(user)
	if err != nil {
		return false, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}
	return user.Admin, nil
}


func (db *mongo) UserAttachFleet(userId string, fleetId string) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)
	user := &models.User{}

	err := users.FindId(userId).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	fleets := db.C(models.FleetsCollection).With(session)
	fleet := &models.Fleet{}

	err = fleets.FindId(fleetId).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "fleet_not_found", "Fleet not found", err)
	}

	user.FleetIds = append(user.FleetIds, fleet.Id)

	err = users.Update(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"fleet_ids": user.FleetIds}})
	if err != nil {
		return user, helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user", err)
	}

	return user, err
}
/*	UserAttachFleet(string, string) (*models.User, error)
	UserDetachFleet(string, string) (*models.User, error)
	UserGetFleet(string, string) (*models.User, error)
	UserGetFleets(string) ([]*models.User, error)
*/