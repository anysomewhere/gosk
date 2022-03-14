package writer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/munnik/gosk/config"
	"github.com/munnik/gosk/database"
	"github.com/munnik/gosk/logger"
	"github.com/munnik/gosk/message"
	"go.nanomsg.org/mangos/v3"
	"go.uber.org/zap"
)

const (
	SignalKEndpointsPath = "/signalk"
	SignalKHTTPPath      = "/signalk/v3/api/"
	SignalKWSPath        = "/signalk/v3/stream"
)

type SignalKWriter struct {
	config           *config.SignalKConfig
	database         *database.PostgresqlDatabase
	cache            *database.BigCache
	wg               *sync.WaitGroup
	websocketClients []websocketClient
}

func NewSignalKWriter(c *config.SignalKConfig) *SignalKWriter {
	return &SignalKWriter{
		config:           c,
		database:         database.NewPostgresqlDatabase(c.PostgresqlConfig),
		cache:            database.NewBigCache(c.BigCacheConfig),
		wg:               &sync.WaitGroup{},
		websocketClients: make([]websocketClient, 0),
	}
}

func (w *SignalKWriter) WriteMapped(subscriber mangos.Socket) {
	// fill the cache with data from the database
	w.wg.Add(1)
	go w.readFromDatabase()
	go w.receive(subscriber)

	router := chi.NewRouter()
	router.Use(middleware.Compress(5))

	router.Get(SignalKHTTPPath+"*", w.serveFullDataModel)
	router.Get(SignalKEndpointsPath, w.serveEndpoints)
	router.Get(SignalKWSPath, w.serveWebsocket)

	// listen to port
	err := http.ListenAndServe(fmt.Sprintf("%s", w.config.URL.Host), router)
	if err != nil {
		logger.GetLogger().Fatal(
			"Could not listen and serve",
			zap.String("Host", w.config.URL.Host),
			zap.String("Error", err.Error()),
		)
	}
}

func (w *SignalKWriter) receive(subscriber mangos.Socket) {
	for {
		received, err := subscriber.Recv()
		if err != nil {
			logger.GetLogger().Warn(
				"Could not receive a message from the publisher",
				zap.String("Error", err.Error()),
			)
			continue
		}
		var mapped message.Mapped
		if err := json.Unmarshal(received, &mapped); err != nil {
			logger.GetLogger().Warn(
				"Could not unmarshal the received data",
				zap.ByteString("Received", received),
				zap.String("Error", err.Error()),
			)
			continue
		}
		go w.updateFullDataModel(mapped)
		go w.updateWebsocket(mapped)
	}
}
