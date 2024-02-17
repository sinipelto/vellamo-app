package main

const CONFIG_DEFAULT = "config/config.json"

type ClientConfig struct {
	TimeoutSec *uint16            `json:"timeoutSec"`
	Headers    *map[string]string `json:"headers"`
}

type Config struct {
	RequestType   *string  `json:"requestType"`
	AreaName      *string  `json:"areaName"`
	SensorName    *string  `json:"sensorName"`
	GeoJsonSuffix *string  `json:"geoJsonSuffix"`
	AlertMin      *float64 `json:"alertMin"`
	AlertMax      *float64 `json:"alertMax"`
	MeasureUnit   *string  `json:"measureUnit"`
	Api           struct {
		Scheme     *string       `json:"scheme"`
		Server     *string       `json:"server"`
		Base       *string       `json:"base"`
		Area       *string       `json:"area"`
		Hourly     *string       `json:"hourly"`
		HttpClient *ClientConfig `json:"httpClient"`
	} `json:"api"`
	Ntfy struct {
		// Bin        *string       `json:"bin"`
		// Options    *string       `json:"options"`
		// PubCmd     *string       `json:"pubCmd"`
		Scheme     *string       `json:"scheme"`
		Server     *string       `json:"server"`
		Topic      *string       `json:"topic"`
		HttpClient *ClientConfig `json:"httpClient"`
	} `json:"ntfy"`
}
