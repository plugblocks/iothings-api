package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ryankurte/go-mapbox/lib"
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/ryankurte/go-mapbox/lib/directions"
	"github.com/ryankurte/go-mapbox/lib/map_matching"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"strconv"
	"strings"
	"time"
)

// Routing enhancer
func CheckLocation(context *gin.Context, store store.Store, device *models.Device, location *models.Geolocation) {
	if device.OrderId != nil {
		order, err := store.GetOrderById(device.OrganizationId, *location.OrderId)
		if err != nil {
			return
		}

		mapBox, err := mapbox.NewMapbox(config.GetString(context, "MAPBOX_API_KEY"))

		directionOpts := directions.RequestOpts{}
		loc := []base.Location{
			{
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
			},
			{
				Latitude:  order.Destination.Latitude,
				Longitude: order.Destination.Longitude,
			},
		}

		fetchedDirections, err := mapBox.Directions.GetDirections(loc, directions.RoutingDrivingTraffic, &directionOpts)
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
			fmt.Println(fastestRoute.Duration)
			order.LiveETA = time.Now().Add(time.Duration(fastestRoute.Duration) * time.Second).Unix()

			s := GetTextSender(context)
			subscription, err := store.GetOrganizationSubscription(device.OrganizationId)
			if err != nil {
				fmt.Println("Get organization subscription error")
				return
			}

			if !order.HasNotifiedDelay && float64(time.Now().Unix())+fastestRoute.Duration > float64(order.ExpectedArrivalTime) {
				textData := models.TextData{PhoneNumber: order.ContactPhoneNumber, Subject: "Text Alert", Message: "La livraison " + order.Reference + " sera en retard. ETA: " + secondsToHours(int(fastestRoute.Duration))}
				err = s.SendText(context, subscription, textData)
				if err != nil {
					return
				}
				order.HasNotifiedDelay = true

				mailData := models.EmailData{ReceiverMail: "adrien@plugblocks.com", ReceiverName: "Adrien Chapelet", Body: "La livraison " + order.Reference + " sera en retard. ETA: " + secondsToHours(int(fastestRoute.Duration)), Subject: "La livraison " + order.Reference + " sera en retard.", AppName: config.GetString(context, "mail_sender_name")}
				s2 := GetEmailSender(context)
				s2.SendEmail(&mailData)
			}

			if fastestRoute.Distance <= 3000 {
				order.Status = models.Arrived.String()
				textData := models.TextData{PhoneNumber: order.ContactPhoneNumber, Subject: "Text Alert", Message: "La livraison " + order.Reference + " est arrivée."}
				s.SendText(context, subscription, textData)

				unixTimeUTC := time.Unix(int64(order.ExpectedArrivalTime), 0) //gives unix time stamp in utc
				mailData := models.EmailData{ReceiverMail: "adrien@plugblocks.com", ReceiverName: "Adrien Chapelet", Body: "La livraison " + order.Reference + " est arrivée. Arrivée estimée:" + unixTimeUTC.Format(time.RFC3339), Subject: "La livraison " + order.Reference + " est arrivée.", AppName: config.GetString(context, "mail_sender_name")}
				s2 := GetEmailSender(context)
				s2.SendEmail(&mailData)
			}
		}

		store.UpdateOrder(order.OrganizationId, order.Id, params.M{"$set": order})
	}
}

func GetMatchingRouteFromGeolocations(context context.Context, locations []*models.Geolocation, order *models.Order) (*mapmatching.MatchingResponse, error) {
	mapBox, err := mapbox.NewMapbox(config.GetString(context, "MAPBOX_API_KEY"))
	if err != nil {
		return nil, err
	}

	var locs []base.Location
	var timestampsData []string

	for _, geoloc := range locations {
		if len(locs) == 100 { // max size is 100
			locs = locs[1:]
		}

		if len(timestampsData) == 100 {
			timestampsData = timestampsData[1:]
		}

		locs = append(locs, base.Location{
			Latitude:  geoloc.Latitude,
			Longitude: geoloc.Longitude,
		})

		timestampsData = append(timestampsData, strconv.FormatInt(geoloc.Timestamp, 10))
	}

	timestampsString := strings.Join(timestampsData, ";")

	ops := &mapmatching.RequestOpts{
		Geometries: mapmatching.GeometryGeojson,
		Overview:   mapmatching.OverviewFull,
		Timestamps: timestampsString,
	}
	matching, err := mapBox.MapMatching.GetMatching(locs, mapmatching.RoutingDriving, ops)
	if err != nil {
		return nil, err
	}

	return matching, nil
}

func secondsToHours(inSeconds int) string {
	hours := inSeconds / 3600
	minutes := inSeconds % 3600 / 60
	str := fmt.Sprintf("%dH%dm", int(hours), int(minutes))
	return str
}
