package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Set by ldflags during compilation
var DEBUGS string
var DEBUG = DEBUGS == "true"
var CONFIGS string

// Configure logger, timestamps, lineno etc
var logger *log.Logger = log.New(os.Stdout, "[vellamo] ", log.Ldate|log.Ltime|log.Lshortfile)

func read_config(fname string) *Config {
	exc, err := os.Executable()
	if err != nil {
		logger.Panicln("ERROR: Failed to get executable path. ERR:", err)
	}

	cwd := filepath.Dir(exc)
	cpth := filepath.Join(cwd, fname)

	f, err := os.Open(cpth)
	if err != nil {
		logger.Panicf("ERROR: Failed to open config file. Double check config file '%s' exists in path: %s\n", fname, cwd)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		logger.Panicln("ERROR: Failed to read config file:", err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		logger.Panicln("ERROR: Failed to parse config from config file. Double check syntax and values. Error:", err)
	}

	return &config
}

func BuildRequest(method string, uri *url.URL, config *ClientConfig, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, uri.String(), body)
	if config.Headers != nil {
		for header, value := range *config.Headers {
			req.Header.Add(header, value)
		}
	}
	return req, err
}

func SendRequest(req *http.Request, config *ClientConfig) (resp *http.Response, err error) {
	// Use default client as base
	client := http.DefaultClient

	if config != nil {
		// Custom timeout if set
		if config.TimeoutSec != nil {
			client.Timeout = time.Duration(*config.TimeoutSec) * time.Second
		}

		if config.Headers != nil {
			for header, value := range *config.Headers {
				req.Header.Add(header, value)
			}
		}
	}

	// Allow all redirs
	// TODO from config?
	client.CheckRedirect = nil

	return client.Do(req)
}

func NtfyPub(config *Config, msg string) {
	// NtfyPub(*config.Ntfy.Bin, config.Ntfy.Options, config.Ntfy.PubCmd, &targetUrlStr, &s)

	urlObj := &url.URL{
		Scheme: *config.Ntfy.Scheme,
		Host:   *config.Ntfy.Server,
		Path:   *config.Ntfy.Topic,
	}

	if DEBUG {
		logger.Println("NTFY URLOBJ:", urlObj.String())
	}

	ntfyUrl, err := url.Parse(urlObj.String())

	if err != nil {
		logger.Panicln("ERROR: Could not parse Ntfy URL:", err)
	}

	if DEBUG {
		logger.Println("NTFY_URL:", ntfyUrl.String())
	}

	bodyBuf := bytes.NewBufferString(msg)
	req, err := BuildRequest(http.MethodPost, ntfyUrl, config.Ntfy.HttpClient, bodyBuf)
	// Check request ok
	if err != nil {
		logger.Panicln("ERROR: Failed to construct Ntfy request:", err)
	}

	// Send request to api
	resp, err := SendRequest(req, config.Ntfy.HttpClient)
	if err != nil {
		logger.Panicln("ERROR: Failed to send Ntfy request:", err)
	}
	// Check response code
	// We follow indefinitely, so will be 2xx at final if OK
	if resp.StatusCode != 200 {
		logger.Panicln("ERROR: Ntfy response was not OK. Status was:", resp.Status)
	}
	defer resp.Body.Close()
}

func main() {
	// Separator for logfile
	logger.Println("---------- PROGRAM START ----------")

	logger.Println("DEBUG MODE:", DEBUG)

	// Fetch config from json file
	var config *Config
	if CONFIGS != "" {
		config = read_config(CONFIGS)
	} else {
		config = read_config(CONFIG_DEFAULT)
	}

	// TODO switch case request type => uri
	// Eg area, hour, month .geojson...

	// Build api uri
	// /api/v1
	path := *config.Api.Base

	switch strings.ToLower(*config.RequestType) {
	case "area":
		// /api/v1 /area /hervanta-3 .geojson
		path, _ = url.JoinPath(path, *config.Api.Area, *config.SensorName)
		path += *config.GeoJsonSuffix
	default:
		logger.Panicln("ERROR: Api request type was not recognized. Was:", *config.RequestType)
	}

	if DEBUG {
		logger.Println("PATH:", path)
	}

	// https:// server.com /api/v1/area
	urlObj := &url.URL{
		Scheme: *config.Api.Scheme,
		Host:   *config.Api.Server,
		Path:   path,
	}

	if DEBUG {
		logger.Println("URLOBJ:", urlObj.String())
	}

	// Final validation
	sensorUrl, err := url.Parse(urlObj.String())
	if err != nil {
		logger.Panicln("ERROR: Could not parse Sensor API URL:", err)
	}

	if DEBUG {
		logger.Println("SENSOR_URL:", sensorUrl.String())
	}

	// Build api request
	apiReq, err := BuildRequest(http.MethodGet, sensorUrl, config.Api.HttpClient, nil)

	// Check request ok
	if err != nil {
		logger.Panicln("ERROR: Failed to construct GET request:", err)
	}

	// Send request to api
	apiResp, err := SendRequest(apiReq, config.Api.HttpClient)
	if err != nil {
		logger.Panicln("ERROR: Failed to GET sensor data:", err)
	}
	// Check response code
	if apiResp.StatusCode != 200 {
		logger.Panicln("ERROR: Failed to GET sensor data. Response was not OK. Status was:", apiResp.Status)
	}
	defer apiResp.Body.Close()

	body, err := io.ReadAll(apiResp.Body)

	if err != nil {
		logger.Panicln("ERROR: Failed to read response body:", err)
	}

	var bodyObj AreaGeoJsonData
	err = json.Unmarshal(body, &bodyObj)
	if err != nil {
		logger.Panicln("ERROR: Failed to unmarshal response body JSON to object:", err)
	}

	// TODO: measuredunits as array
	// loop through
	// check min max for each
	// Optionally: map each unit to its own ntfy topic
	var val float64
	switch strings.ToLower(*config.MeasureUnit) {
	case "t":
		val = bodyObj.Properties.T
	case "cl":
		val = bodyObj.Properties.Cl
	case "hardness":
		val = bodyObj.Properties.Hardness
	case "ph":
		val = bodyObj.Properties.Ph
	default:
		logger.Panicln("ERROR: Target Measurement unit was not recognized. Was:", *config.MeasureUnit)
	}

	logger.Printf("INFO: Measured Unit: '%s' Current value: '%v'\n", *config.MeasureUnit, val)

	// Build ntfy target url = server + topic
	urlObj = &url.URL{
		Scheme: *config.Ntfy.Scheme,
		Host:   *config.Ntfy.Server,
		Path:   *config.Ntfy.Topic,
	}

	targetUrl, err := url.Parse(urlObj.String())

	if err != nil {
		logger.Panicln("ERROR: Could not parse ntfy URL from config:", err)
	}

	if DEBUG {
		logger.Println("TARGET URL:", targetUrl.String())
	}

	alerted := false

	if config.AlertMin != nil {
		if val <= *config.AlertMin {
			s := fmt.Sprintf("INFO: Unit '%s' Value below or MIN. Threshold: %v Value: %v", *config.MeasureUnit, *config.AlertMin, val)
			logger.Println(s)
			NtfyPub(config, s)
			alerted = true
		}
	}

	if config.AlertMax != nil {
		if val >= *config.AlertMax {
			s := fmt.Sprintf("INFO: Unit '%s' Value above or MAX. Threshold: %v Value: %v", *config.MeasureUnit, *config.AlertMax, val)
			logger.Println(s)
			NtfyPub(config, s)
			alerted = true
		}
	}

	if !alerted {
		logger.Println("INFO: No alerts were triggered.")
	}

	logger.Println("---------- END OF PROGRAM ----------")
}
