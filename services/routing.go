package services

import (
	"context"
	"github.com/ryankurte/go-mapbox/lib"
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/ryankurte/go-mapbox/lib/directions"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"time"
)

// Routing enhancer
func CheckLocation(context context.Context, store store.Store, device *models.Device, location *models.Geolocation) {
	if device.OrderId != nil {
		order, err := store.GetOrderById(device.OrganizationId, *location.OrderId)
		if err != nil {
			return
		}

		mapBox, err := mapbox.NewMapbox(config.GetString(context, "MAPBOX_API_KEY"))

		var directionOpts directions.RequestOpts
		loc := []base.Location{
			{
				location.Longitude,
				location.Latitude,
			},
			{
				order.Destination.Longitude,
				order.Destination.Latitude,
			},
		}

		fetchedDirections, err := mapBox.Directions.GetDirections(loc, directions.RoutingDriving, &directionOpts)
		if err != nil {
			return
		}

		// FIRST MESSAGE OF ORDER
		if order.Status == models.Created.String() {
			order.Status = models.Transiting.String()
		}

		// TRANSITING
		if order.Status == models.Transiting.String() {

			var fastestRoute *directions.Route
			for _, route := range fetchedDirections.Routes {
				if fastestRoute == nil || route.Duration < fastestRoute.Duration {
					fastestRoute = &route
				}
			}

			s := GetTextSender(context)

			if !order.HasNotifiedDelay && float64(time.Now().Unix())+fastestRoute.Duration > float64(order.ExpectedArrivalTime) {
				data := models.TextData{PhoneNumber: order.ContactPhoneNumber, Subject: "Text Alert", Message: "La livraison " + order.Reference + " sera en retard."}
				err = s.SendText(data)
				if err != nil {
					return
				}
				order.HasNotifiedDelay = true
			}

			if fastestRoute.Distance <= 3000 {
				order.Status = models.Arrived.String()
				data := models.TextData{PhoneNumber: order.ContactPhoneNumber, Subject: "Text Alert", Message: "La livraison " + order.Reference + " est arrivée."}
				s.SendText(data)
			}
		}

		store.UpdateOrder(order.OrganizationId, order.Id, params.M{"$set": order})
	}
}
