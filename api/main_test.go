package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"unsafe"

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

// TestTimeseriesExtractWatershed
func TestTimeseriesExtractWatershed(t *testing.T) {
	e := echo.New() // All Routes

	q := make(url.Values)
	q.Set("after", "2021-10-15T00:00:00Z")
	q.Set("before", "2021-10-22T23:00:00Z")

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
		fmt.Printf("Body size: %T, %d\n", b, unsafe.Sizeof(b))
	}

}
