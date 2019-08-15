package models

type HereRequest struct {
	Wlan []Wlan `json:"wlan"`
}
type Wlan struct {
	Mac string `json:"mac"`
}

type HereError struct {
	Error struct {
		Code        int    `json:"code"`
		Message     string `json:"message"`
		Description string `json:"description"`
	} `json:"error"`
}

type HereLocation struct {
	Location struct {
		Lat      float64 `json:"lat"`
		Lng      float64 `json:"lng"`
		Accuracy int     `json:"accuracy"`
	} `json:"location"`
}
