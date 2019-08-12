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
	DeleteUser(*models.User, string) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(params.M) (*models.User, error)
	UpdateUser(*models.User, params.M) error
	ChangeLanguage(string, string) error
	GetUsers() ([]*models.User, error)
	CountUsers() (int, error)
	AssignOrganization(string, string) error
	GetUserOrganization(user *models.User) (*models.Organization, error)

	CreateCustomer(string, *models.Customer) error
	GetCustomerById(string, string) (*models.Customer, error)
	UpdateCustomer(string, string, params.M) error
	GetAllCustomers(string) ([]models.Customer, error)
	DeleteCustomer(string, string) error

	CreateFleet(*models.User, *models.Fleet) error
	AddDeviceToFleet(user *models.User, fleetId string, deviceId string) (*models.Fleet, error)
	GetFleetById(*models.User, string) (*models.Fleet, error)
	UpdateFleet(*models.User, string, params.M) error
	GetFleets(*models.User) ([]*models.Fleet, error)
	DeleteFleet(*models.User, string) error
	CountFleets() (int, error)
	GetDevicesFromFleet(*models.User, string) ([]*models.Device, error)

	CreateDevice( /**models.User*/ string, *models.Device) error
	GetDevices(*models.User) ([]*models.Device, error)
	GetAvailableDevices(string) ([]*models.Device, error)
	UpdateDevice(string, string, params.M) error
	DeleteDevice(*models.User, string) error
	GetDevice(string, string) (*models.Device, error)
	GetDeviceGeolocations(string, string, int, int, int) ([]*models.Geolocation, error)
	GetDevicePreciseGeolocations(string, int, int, int) ([]*models.Geolocation, error)
	GetDeviceFromSigfoxId(string) (*models.Device, error)
	CountDevices() (int, error)
	DeleteDeviceObservations(string) error
	DeleteDeviceGeolocations(string) error
	UpdateDeviceActivity(deviceId string, actDiff int) error
	GetDeviceMessages(string) ([]*sigfox.Message, error)

	GetAlerts(*models.User) ([]*models.Alert, error)
	CreateAlert(*models.User, *models.Alert) error
	GetAlert(*models.User, string) (*models.Alert, error)
	GetFleetAlerts(*models.User, string) ([]*models.Alert, error)
	GetDeviceAlerts(*models.User, string) ([]*models.Alert, error)
	UpdateAlert(*models.User, string, params.M) error
	DeleteAlert(*models.User, string) error
	CountAlerts() (int, error)

	CreateSigfoxMessage(*sigfox.Message) error
	CreateSigfoxLocation(*sigfox.Location) error
	GetSigfoxLocations() ([]sigfox.Location, error)
	GetGeoJSON() (*models.GeoJSON, error)
	CountSigfoxMessages() (int, error)

	CreateGeolocation(*models.Geolocation) error
	//TODO: DANGER: Protect by auth device GeoJSON
	GetFleetsGeoJSON( /**models.User*/ string, int, int, int) (*models.GeoJSON, error)
	//TODO: DANGER: Protect by auth device GeoJSON
	GetFleetGeoJSON( /* *models.User, */ string, string, int, int, int) (*models.GeoJSON, error)
	//TODO: DANGER: Protect by auth device GeoJSON
	GetDeviceGeoJSON( /* *models.User, */ string, string, int, int, int) (*models.GeoJSON, error)
	//TODO: DANGER: Protect by auth device GeoJSON
	GetDeviceGeolocation(*models.User, string, string) (*models.Geolocation, error)
	GetUserFleetsGeoJSON(*models.User) (*models.GeoJSON, error)
	CountGeolocations() (int, error)

	CreateOrganization(*models.Organization) error
	GetOrganizationById(string) (*models.Organization, error)
	GetOrganizationUsers(string) ([]models.SanitizedUser, error)
	UpdateOrganization(string, params.M) error
	GetAllOrganizations() ([]models.Organization, error)
	DeleteOrganization(string) error
	CountOrganizations() (int, error)
	GetOrganizationSubscription(string) (*models.Subscription, error)

	CreateSubscription(*models.Subscription) error
	GetSubscription(string) (*models.Subscription, error)
	UpdateSubscription(string, params.M) error
	GetSubscriptions() ([]models.Subscription, error)
	DeleteSubscription(string) error

	CreateObservation(*models.Observation) error
	GetDeviceObservations(string, string, string, int) ([]*models.Observation, error)
	GetFleetObservations(*models.User, string, string, string, int) ([]*models.Observation, error)
	CountObservations() (int, error)
	DeleteObservation(string) error

	CreateOrder(string, *models.Order) error
	GetOrderById(string, string) (*models.Order, error)
	UpdateOrder(string, string, params.M) error
	GetAllOrders(string) ([]models.Order, error)
	DeleteOrder(string, string) error
	TerminateOrder(string, string) error
	GetOrderGeolocations(string, string) ([]*models.Geolocation, error)

	CreateWarehouse(string, *models.Warehouse) error
	GetWarehouseById(string, string) (*models.Warehouse, error)
	UpdateWarehouse(string, string, params.M) error
	GetAllWarehouses(string) ([]models.Warehouse, error)
	DeleteWarehouse(string, string) error
}
