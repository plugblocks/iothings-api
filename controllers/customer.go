package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type CustomerController struct{}

func NewCustomerController() CustomerController {
	return CustomerController{}
}

func (uc CustomerController) GetCustomer(c *gin.Context) {
	customer, err := store.FindCustomerById(c, c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("customer_not_found", "The customer does not exist", err))
		return
	}

	c.JSON(http.StatusOK, customer.Sanitize())
}

func (uc CustomerController) CreateCustomer(c *gin.Context) {
	customer := &models.Customer{}

	if err := c.BindJSON(customer); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	appName := config.GetString(c, "mail_sender_name")
	subject := "Welcome to " + appName + "! Please confirm your account"
	templateLink := "./templates/html/mail_activate_account.html"

	s := services.GetEmailSender(c)
	data := models.EmailData{ReceiverMail: customer.Email, ReceiverName: customer.Firstname + " " + customer.Lastname /*User: customer,*/, Subject: subject, ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
	err := s.SendEmailFromTemplate(c, &data, templateLink)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_credit_spent", "Your mail credit is spent", err))
		return
	}

	if err := store.CreateCustomer(c, customer); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, customer.Sanitize())
}

func (uc CustomerController) DeleteCustomer(c *gin.Context) {
	err := store.DeleteCustomer(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uc CustomerController) ActivateCustomer(c *gin.Context) {
	if err := store.ActivateCustomer(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	//c.JSON(http.StatusOK, nil)

	/*customer, err := store.FindCustomerById(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("customer_not_found", "The customer does not exist", err))
		return
	}

	vars := gin.H{
		"Customer": customer,
		"AppName": config.GetString(c, "mail_sender_name"),
		"AppUrl": config.GetString(c, "front_url"),
	}

	c.HTML(http.StatusOK, "./templates/html/page_account_activated.html", vars)*/

	c.Redirect(http.StatusMovedPermanently, "https://"+config.GetString(c, "front_url"))
}

func (uc CustomerController) GetCustomers(c *gin.Context) {
	customers, err := store.GetCustomers(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (uc CustomerController) AssignOrganization(c *gin.Context) {
	customerId := c.Param("id")
	organizationId := c.Param("organization_id")

	err := store.AssignOrganization(c, customerId, organizationId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uc CustomerController) GetCustomerOrganization(c *gin.Context) {
	customer, err := store.FindCustomerById(c, c.Param("id"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	organization, err := store.GetCustomerOrganization(c, customer)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, organization)
}
