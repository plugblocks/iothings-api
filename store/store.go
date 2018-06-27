package store

import (
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
)

type Store interface {
	CreateUser(*models.User) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(params.M) (*models.User, error)
	UpdateUser(*models.User, params.M) error
	GetUsers() ([]*models.User, error)
	AssignOrganization(string, string) error
	GetUserOrganization(user *models.User) (*models.Organization, error)

	CreateDevice(*models.User, *models.Device) error
	GetDevices(*models.User, string) ([]*models.Device, error)
	UpdateDevice(*models.User, string, params.M) error
	DeleteDevice(*models.User, string) error
	GetDevice(*models.User, string) (*models.Device, error)
	GetDeviceFromSigfoxId(string) (*models.Device, error)

	CreateSigfoxMessage(*sigfox.Message) error
	CreateSigfoxLocation(location *sigfox.Location) error

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

	CreateOrganization(*models.Organization) error
	GetOrganizationById(string) (*models.Organization, error)
	GetOrganizationUsers(string) ([]models.SanitizedUser, error)
	UpdateOrganization(string, params.M) error
	GetAllOrganizations() ([]models.Organization, error)
	DeleteOrganization(string) error

	CreateObservation(*models.Observation) error
	GetDeviceObservations(string, string, string) ([]models.Observation, error)
	GetDeviceLatestObservation(string, string, string) (*models.Observation, error)
	GetFleetObservations(*models.User, string, string) ([]models.Observation, error)
	GetFleetLatestObservation(*models.User, string, string) ([]models.Observation, error)
}
