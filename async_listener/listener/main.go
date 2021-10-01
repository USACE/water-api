package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/USACE/go-simple-asyncer/asyncer"
	"github.com/kelseyhightower/envconfig"
	"github.com/lib/pq"
)

// Message holds Fn (function name) and details (string to send to queue / used by specific worker)
type Message struct {
	Fn      string `json:"fn"`
	Details string `json:"details"`
}

// NotificationHandler is a function that takes a notification and returns an error
type NotificationHandler func(string) error

// Config holds application configuration variables
type Config struct {
	DBUser    string `envconfig:"WATER_DB_USER"`
	DBPass    string `envconfig:"WATER_DB_PASS"`
	DBName    string `envconfig:"WATER_DB_NAME"`
	DBHost    string `envconfig:"WATER_DB_HOST"`
	DBSSLMode string `envconfig:"WATER_DB_SSLMODE"`
	// AsyncEnginePackager         string `envconfig:"ASYNC_ENGINE_PACKAGER"`
	// AsyncEnginePackagerTarget   string `envconfig:"ASYNC_ENGINE_PACKAGER_TARGET"`
	// AsyncEngineStatistics       string `envconfig:"ASYNC_ENGINE_STATISTICS"`
	// AsyncEngineStatisticsTarget string `envconfig:"ASYNC_ENGINE_STATISTICS_TARGET"`
	AsyncEngineGeoprocess       string `envconfig:"ASYNC_ENGINE_GEOPROCESS"`
	AsyncEngineGeoprocessTarget string `envconfig:"ASYNC_ENGINE_GEOPROCESS_TARGET"`
	MaxReconn                   string `envconfig:"MAX_RECONN"`
	MinReconn                   string `envconfig:"MIN_RECONN"`
}

// connStr returns a database connection string
func (c Config) connStr() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=%s binary_parameters=yes",
		c.DBUser, c.DBPass, c.DBName, c.DBHost, c.DBSSLMode,
	)
}

func (c Config) minReconn() time.Duration {
	d, err := time.ParseDuration(c.MinReconn)
	if err != nil {
		panic(err.Error())
	}
	return d
}

func (c Config) maxReconn() time.Duration {
	d, err := time.ParseDuration(c.MaxReconn)
	if err != nil {
		panic(err.Error())
	}
	return d
}

// NewAsyncNotificationHandler handles dependency injection of asyncer.Asyncer
func NewAsyncNotificationHandler(a asyncer.Asyncer) NotificationHandler {
	return func(d string) error {
		if err := a.CallAsync([]byte(d)); err != nil {
			fmt.Println("Error calling async")
			fmt.Println(err.Error())
			return err
		}
		return nil
	}
}

func waitForNotification(l *pq.Listener, handlers map[string]NotificationHandler) {
	select {
	case n := <-l.Notify:
		fmt.Println("notification on channel: " + n.Channel)
		var m Message
		if err := json.Unmarshal([]byte(n.Extra), &m); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		if handler, ok := handlers[m.Fn]; ok {
			go handler(m.Details)
		} else {
			fmt.Printf("Unimplemented handler for Function (fn) %s\n", m.Fn)
		}
	case <-time.After(90 * time.Second):
		go l.Ping()
		fmt.Println("received no work for 90 seconds; checking for new work")
	}
}

func reportProblem(eq pq.ListenerEventType, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {

	var cfg Config
	if err := envconfig.Process("water", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	// Database Listener
	listener := pq.NewListener(cfg.connStr(), cfg.minReconn(), cfg.maxReconn(), reportProblem)
	// Start Listening for Productfiles
	if err := listener.Listen("water_new"); err != nil {
		panic(err)
	}

	// // downloadAsyncer defines async engine used to package DSS files for download
	// downloadAsyncer, err := asyncer.NewAsyncer(asyncer.Config{Engine: cfg.AsyncEnginePackager, Target: cfg.AsyncEnginePackagerTarget})
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// d := NewAsyncNotificationHandler(downloadAsyncer)

	// acquirablefileAsyncer defines async engine for processing new acquirable files
	geoprocessAsyncer, err := asyncer.NewAsyncer(asyncer.Config{Engine: cfg.AsyncEngineGeoprocess, Target: cfg.AsyncEngineGeoprocessTarget})
	if err != nil {
		log.Fatal(err.Error())
	}
	g := NewAsyncNotificationHandler(geoprocessAsyncer)

	// Map of handlers
	handlers := map[string]NotificationHandler{
		"geoprocess-shapefile-upload": g,
	}

	fmt.Println("entering main loop")
	for {
		waitForNotification(listener, handlers)
	}
}
