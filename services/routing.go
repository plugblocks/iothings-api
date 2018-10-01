package services

import (
	"context"
	"fmt"
	"github.com/ryankurte/go-mapbox/lib"
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/ryankurte/go-mapbox/lib/directions"
	"github.com/ryankurte/go-mapbox/lib/map_matching"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"time"
	"strconv"
)

// Routing enhancer
func CheckLocation(context context.Context, store store.Store, device *models.Device, location *models.Geolocation) {
	if device.OrderId != nil {
		order, err := store.GetOrderById(device.OrganizationId, *location.OrderId)
		if err != nil {
			return
		}

		mapBox, err := mapbox.NewMapbox(config.GetString(context, "MAPBOX_API_KEY"))

		directionOpts := directions.RequestOpts{}
		loc := []base.Location{
			{
				location.Latitude,
				location.Longitude,
			},
			{
				order.Destination.Latitude,
				order.Destination.Longitude,
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
			if !order.HasNotifiedDelay && float64(time.Now().Unix())+fastestRoute.Duration > float64(order.ExpectedArrivalTime) {
				data := models.TextData{PhoneNumber: order.ContactPhoneNumber, Subject: "Text Alert", Message: "La livraison " + order.Reference + " sera en retard. ETA: " + secondsToHours(int(fastestRoute.Duration))}
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

func GetMatchingRouteFromGeolocations(context context.Context, locations []*models.Geolocation, order *models.Order) (*mapmatching.MatchingResponse, error) {
	mapBox, err := mapbox.NewMapbox(config.GetString(context, "MAPBOX_API_KEY"))
	if err != nil {
		return nil, err
	}

	var locs []base.Location
	var timestampsData []int64

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

		timestampsData = append(timestampsData, geoloc.Timestamp)
	}

	timestampsString := ""
	for _, timestamp := range timestampsData {
		timestampsString += strconv.FormatInt(timestamp,10) + ";"
	}

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
