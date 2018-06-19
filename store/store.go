package store

import (
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

type Store interface {
	CreateUser(*models.User) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(params.M) (*models.User, error)
	UpdateUser(*models.User, params.M) error
	GetUsers() ([]*models.User, error)

	UserIsAdmin(string) (bool, error)
	UserAttachFleet(string, string) (*models.User, error)
	UserDetachFleet(string, string) (*models.User, error)
	UserGetFleet(string, string) (*models.User, error)
	UserGetFleets(string) ([]*models.User, error)

	CreateDevice(*models.Device) error
	GetDevices(string) ([]*models.Device, error)
	UpdateDevice(string, params.M) error
	DeleteDevice(string) error
	GetDevice(string) (*models.Device, error)

	CreateGroup(*models.Group) error
}
