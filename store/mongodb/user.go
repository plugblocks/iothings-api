package mongodb

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
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

func (db *mongo) DeleteUser(user *models.User, userId string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Remove(bson.M{"_id": userId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_delete_failed", "Failed to delete the user", err)
	}

	return nil
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

func (db *mongo) ChangeLanguage(id string, language string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"language": language}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to change language", err)
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

func (db *mongo) AssignOrganization(user_id string, organization_id string) error {
	session := db.Session.Copy()
	defer session.Close()

	users := db.C(models.UsersCollection).With(session)
	err := users.Update(bson.M{"_id": user_id}, bson.M{"$set": bson.M{"organization_id": organization_id}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_change_failed", "Couldn't change the user organization", err)
	}
	return nil
}

func (db *mongo) GetUserOrganization(user *models.User) (*models.Organization, error) {
	session := db.Session.Copy()
	defer session.Close()
	organizations := db.C(models.OrganizationsCollection).With(session)

	organization := &models.Organization{}

	err := organizations.Find(bson.M{"_id": user.OrganizationId}).One(organization)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "organization_not_found", "Organization not found", err)
	}

	return organization, err
}

func (db *mongo) CountUsers() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	users := db.C(models.UsersCollection).With(session)

	nbr, err := users.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "users_not_found", "Users not found", err)
	}
	return nbr, nil
}
