package connector

import (
	"github.com/munnik/gosk/config"
	"github.com/munnik/gosk/logger"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.nanomsg.org/mangos/v3"
	"go.uber.org/zap"
	"net/url"
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

func (*CanBusConnector) Subscribe(mangos.Socket) {
	// do nothing
}

func networkFromURL(u *url.URL) string {
	switch u.Scheme {
	case "sock":
		return "can"
	default:
		return u.Scheme
	}
}

func (r *CanBusConnector) receive(stream chan<- []byte) error {
	conn, err := socketcan.Dial(networkFromURL(r.config.URL), r.config.URL.Host)
	if err != nil {
		return err
	}
	recv := socketcan.NewReceiver(conn, socketcan.ReceiverFrameInterceptor(handleCanFrameStream(stream)))
	defer recv.Close()

	for recv.Receive() {
		recv.Frame()
	}
	return nil
}

func handleCanFrameStream(stream chan<- []byte) socketcan.FrameInterceptor {
	return func(frm can.Frame) {
		stream <- []byte(frm.String())
	}
}
