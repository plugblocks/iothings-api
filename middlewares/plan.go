package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/store"
	"time"
	"strconv"
)

func PlanMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := store.Current(c)

		orga, err := store.GetOrganizationById(c, user.OrganizationId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("organization_find_failed", "Failed to find the organization related to user", err))
			return
		}

		exp, _ := strconv.Atoi(orga.PlanExpiration)
		mail, _ := strconv.Atoi(orga.PlanCreditMails)
		text, _ := strconv.Atoi(orga.PlanCreditTexts)
		wifi, _ := strconv.Atoi(orga.PlanCreditWifi)

		if exp <= int(time.Now().Unix()) {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("plan_expired", "The plan is expired, please renew your plan", errors.New("The user is not administrator")))
			return
		}

		if mail <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The mail credit is empty, please buy a new pack", errors.New("Orga: "+orga.Name+" mail credit empty")))
			return
		}
		if text <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The texts credit is empty, please buy a new pack", errors.New("Orga: "+orga.Name+" texts credit empty")))
			return
		}
		if wifi <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The wifi resolving credit is empty, please buy a new pack", errors.New("Orga: "+orga.Name+" wifi credit empty")))
			return
		}

		c.Next()
	}
}
