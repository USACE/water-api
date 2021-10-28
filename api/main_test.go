package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/cwms"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	config app.Config
	st     *app.PGStore
	err    error
	cs     cwms.Store
)

func init() {
	config = app.Config{
		ApplicationKey:        "appkey",
		AuthMocked:            true,
		DBUser:                "water_user",
		DBPass:                "water_pass",
		DBName:                "postgres",
		DBHost:                "localhost",
		DBSSLMode:             "disable",
		DBPoolMaxConns:        10,
		DBPoolMaxConnIdleTime: 10,
		DBPoolMinConns:        5,
	}

	// parse configuration from environment variables
	if err := envconfig.Process("water", &config); err != nil {
		log.Fatal(err.Error())
	}
	// create store (database pool) from configuration
	st, err = app.NewStore(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	cs = cwms.Store{Connection: st.Connection}
}

// Setup
func Setup() (*http.Request, *httptest.ResponseRecorder, echo.Context) {
	e := echo.New() // All Routes
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return req, rec, c
}

// TestListOffices
func TestListOffices(t *testing.T) {
	_, rec, c := Setup()
	c.SetPath("/offices")

	if assert.NoError(t, cs.ListOffices(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// TestListLevelKind
func TestListLevelKind(t *testing.T) {
	_, rec, c := Setup()
	c.SetPath("/levels/kind")

	if assert.NoError(t, cs.ListLevelKind(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// TestCreateLevelKind
func TestCreateLevelKind(t *testing.T) {
	_, rec, c := Setup()
	c.SetPath("/levels/kind/:name")
	c.SetParamNames("name")
	c.SetParamValues("Test Location Level")

	if assert.NoError(t, cs.CreateLevelKind(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// TestDeleteLevelKind
func TestDeleteLevelKind(t *testing.T) {
	_, rec, c := Setup()
	c.SetPath("/levels/kind/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("test-location-level")

	if assert.NoError(t, cs.DeleteLevelKind(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// CreateLocationLevels
func TestCreateLocationLevels(t *testing.T) {
	var (
		levelResponse = `[{"kind_id": "43e6ecff-32d0-4e03-ba79-f05a9ed5924d","levels": [{"time": "1900-01-01T06:00:00-06:00", "value": 663},{"time": "1900-12-31T06:00:00-06:00", "value": 663}]},{"kind_id": "b3e8fbb0-ae51-4f56-b2b0-b39658f72375","levels": [{"time": "1900-01-01T06:00:00-06:00", "value": 651},{"time": "1900-12-31T06:00:00-06:00", "value": 651}]}]`
	)
	var request = struct {
		path        string
		contenttype string
		body        string
	}{
		path:        "/levels/:location_id",
		contenttype: "application/json",
		body:        levelResponse,
	}
	// Setup the request parameters and body
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(levelResponse))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath(request.path)
	c.SetParamNames("location_id")
	c.SetParamValues("67d388ca-0c28-44b3-bc43-018e85e939be")

	if assert.NoError(t, cs.CreateLocationLevels(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// TestUpdateLocationLevels
func TestUpdateLocationLevels(t *testing.T) {
	var (
		levelResponse = `[{"kind_id": "43e6ecff-32d0-4e03-ba79-f05a9ed5924d","levels": [{"time": "1900-01-01T06:00:00-06:00", "value": 999},{"time": "1900-12-31T06:00:00-06:00", "value": 999}]},{"kind_id": "b3e8fbb0-ae51-4f56-b2b0-b39658f72375","levels": [{"time": "1900-01-01T06:00:00-06:00", "value": 999},{"time": "1900-12-31T06:00:00-06:00", "value": 999}]}]`
	)
	var request = struct {
		path        string
		contenttype string
		body        string
	}{
		path:        "/levels/:location_id",
		contenttype: "application/json",
		body:        levelResponse,
	}
	// Setup the request parameters and body
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(levelResponse))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath(request.path)
	c.SetParamNames("location_id")
	c.SetParamValues("67d388ca-0c28-44b3-bc43-018e85e939be")

	if assert.NoError(t, cs.UpdateLocationLevels(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// TestListLevelValues
func TestListLevelValues(t *testing.T) {
	_, rec, c := Setup()
	c.SetPath("/levels/:location_slug/:level_kind")
	c.SetParamNames("location_slug", "level_kind")
	c.SetParamValues("dale-hollow", "top-of-flood-control")

	if assert.NoError(t, cs.ListLevelValues(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
	}
}

// TestTimeseriesExtractWatershed
func TestTimeseriesExtractWatershed(t *testing.T) {
	e := echo.New() // All Routes

	q := make(url.Values)
	q.Set("after", "2021-10-15T00:00:00Z")
	q.Set("before", "2021-10-15T23:00:00Z")

	req := httptest.NewRequest(http.MethodGet, "/watersheds/:watershed_slug/extract?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("watershed_slug")
	c.SetParamValues("kanawha-river")

	if assert.NoError(t, cs.TimeseriesExtractWatershed(c)) {
		b := rec.Body.String()
		var out bytes.Buffer
		json.Indent(&out, []byte(b), "", "    ")
		fmt.Printf("%s", out.Bytes())
		// fmt.Printf("Body size: %T, %d\n", b, unsafe.Sizeof(b))
	}

}
