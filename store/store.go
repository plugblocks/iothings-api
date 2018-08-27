package store

import (
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
)

type Store interface {
	CreateGroup(*models.User, *models.Group) error
	GetGroupById(*models.User, string) (*models.Group, error)
	UpdateGroup(*models.User, string, params.M) error
	GetAllGroups(*models.User) ([]models.Group, error)
	DeleteGroup(*models.User, string) error
	CountGroups() (int, error)

	CreateUser(*models.User) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(params.M) (*models.User, error)
	UpdateUser(*models.User, params.M) error
	GetUsers() ([]*models.User, error)
	CountUsers() (int, error)
	AssignOrganization(string, string) error
	GetUserOrganization(user *models.User) (*models.Organization, error)

	CreateFleet(*models.User, *models.Fleet) error
	AddDeviceToFleet(user *models.User, fleetId string, deviceId string) (*models.Fleet, error)
	GetFleetById(*models.User, string) (*models.Fleet, error)
	UpdateFleet(*models.User, string, params.M) error
	GetAllFleets(*models.User) ([]models.Fleet, error)
	DeleteFleet(*models.User, string) error
	CountFleets() (int, error)

	CreateDevice(*models.User, *models.Device) error
	GetDevices(*models.User) ([]*models.Device, error)
	UpdateDevice(*models.User, string, params.M) error
	DeleteDevice(*models.User, string) error
	GetDevice(*models.User, string) (*models.Device, error)
	GetDeviceFromSigfoxId(string) (*models.Device, error)
	CountDevices() (int, error)

	CreateSigfoxMessage(*sigfox.Message) error
	CreateSigfoxLocation(*sigfox.Location) error
	GetSigfoxLocations() ([]sigfox.Location, error)
	GetGeoJSON() (*models.GeoJSON, error)
	CountSigfoxMessages() (int, error)

	CreateGeolocation(*models.Geolocation) error
	//TODO: DANGER: Protect by auth device GeoJSON
	GetFleetGeoJSON( /* *models.User, */ string) (*models.GeoJSON, error)
	//TODO: DANGER: Protect by auth device GeoJSON
	GetDeviceGeoJSON( /* *models.User, */ string) (*models.GeoJSON, error)
	//TODO: DANGER: Protect by auth device GeoJSON
	GetFleetsGeoJSON( /**models.User*/ string, int, int, int) (*models.GeoJSON, error)
	GetUserFleetsGeoJSON(*models.User) (*models.GeoJSON, error)
	CountGeolocations() (int, error)

	CreateOrganization(*models.Organization) error
	GetOrganizationById(string) (*models.Organization, error)
	GetOrganizationUsers(string) ([]models.SanitizedUser, error)
	UpdateOrganization(string, params.M) error
	GetAllOrganizations() ([]models.Organization, error)
	DeleteOrganization(string) error
	CountOrganizations() (int, error)

	CreateObservation(*models.Observation) error
	GetDeviceObservations(string, string, int) ([]*models.Observation, error)
	GetFleetObservations(*models.User, string, string, int) ([]*models.Observation, error)
	CountObservations() (int, error)
}
