package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"time"
)

func PlanMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expiration := config.GetInt(c, "plan_expiration")
		mailCredit := config.GetInt(c, "plan_credit_mail")
		textCredit := config.GetInt(c, "plan_credit_text")
		wifiCredit := config.GetInt(c, "plan_credit_wifi")

		if expiration <= int(time.Now().Unix()) {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("plan_expired", "The plan is expired, please renew your plan", errors.New("The plan is expired")))
			return
		}

		if mailCredit <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The mail credit is empty, please buy a new pack", errors.New("Mail credit empty")))
			return
		}
		if textCredit <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The text credit is empty, please buy a new pack", errors.New("Text credit empty")))
			return
		}
		if wifiCredit <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The wifi resolving credit is empty, please buy a new pack", errors.New("Wifi credit empty")))
			return
		}

		c.Next()
	}
}
