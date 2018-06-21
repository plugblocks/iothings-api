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

	CreateDevice(*models.Device) error
	GetDevices(string) ([]*models.Device, error)
	UpdateDevice(string, params.M) error
	DeleteDevice(string) error
	GetDevice(string) (*models.Device, error)

	CreateGroup(*models.User, *models.Group) error
	GetGroupById(*models.User, string) (*models.Group, error)
	UpdateGroup(*models.User, string, params.M) error
	GetAllGroups(*models.User) ([]models.Group, error)
	DeleteGroup(*models.User, string) error

	CreateFleet(*models.User, *models.Fleet) error
	GetFleetById(*models.User, string) (*models.Fleet, error)
	UpdateFleet(*models.User, string, params.M) error
	GetAllFleets(*models.User) ([]models.Fleet, error)
	DeleteFleet(*models.User, string) error
}
