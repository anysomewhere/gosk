package connector

import (
	"github.com/munnik/gosk/config"
	"github.com/munnik/gosk/logger"
	"go.nanomsg.org/mangos/v3"
	"go.uber.org/zap"

	"github.com/brutella/can"
)

type CanBusConnector struct {
	config *config.ConnectorConfig
}

func NewCanBusConnector(c *config.ConnectorConfig) (*CanBusConnector, error) {
	return &CanBusConnector{config: c}, nil
}

func (r *CanBusConnector) Publish(publisher mangos.Socket) {
	stream := make(chan []byte, 1)
	defer close(stream)
	go func() {
		for {
			if err := r.receive(stream); err != nil {
				logger.GetLogger().Warn(
					"Error while receiving data for the stream",
					zap.String("URL", r.config.URL.String()),
					zap.String("Error", err.Error()),
				)
			}
		}
	}()
	process(stream, r.config.Name, r.config.Protocol, publisher)
}

func (*CanBusConnector) Subscribe(subscriber mangos.Socket) {
	// do nothing
}

func (r *CanBusConnector) receive(stream chan<- []byte) error {
	bus, err := can.NewBusForInterfaceWithName(r.config.URL.Host)
	if err != nil {
		return err
	}
	defer bus.Disconnect()
	bus.SubscribeFunc(handleCanFrameStream(stream))
	err = bus.ConnectAndPublish()
	if err != nil {
		return err
	}
	return nil
}

func handleCanFrameStream(stream chan<- []byte) can.HandlerFunc {
	return func(frm can.Frame) {
		// fmt.Println(frm)
		bytes, err := can.Marshal(frm)
		if err != nil {
			return
		}
		stream <- bytes

	}
}
