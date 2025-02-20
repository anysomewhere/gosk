package writer

import (
	"encoding/json"
	"fmt"

	"github.com/munnik/gosk/logger"
	"github.com/munnik/gosk/message"
	"go.nanomsg.org/mangos/v3"
	"go.uber.org/zap"
)

type StdOutWriter struct {
}

func NewStdOutWriter() *StdOutWriter {
	return &StdOutWriter{}
}

func (w *StdOutWriter) WriteMapped(subscriber mangos.Socket) {
	for {
		received, err := subscriber.Recv()
		if err != nil {
			logger.GetLogger().Warn(
				"Could not receive a message from the publisher",
				zap.String("Error", err.Error()),
			)
			continue
		}
		var m message.Mapped
		if err := json.Unmarshal(received, &m); err != nil {
			logger.GetLogger().Warn(
				"Could not unmarshal a message from the publisher",
				zap.String("Error", err.Error()),
			)
			return
		}
		fmt.Println(m)
	}
}

func (w *StdOutWriter) WriteRaw(subscriber mangos.Socket) {
	for {
		received, err := subscriber.Recv()
		if err != nil {
			logger.GetLogger().Warn(
				"Could not receive a message from the publisher",
				zap.String("Error", err.Error()),
			)
			continue
		}
		raw := message.Raw{}
		if err := json.Unmarshal(received, &raw); err != nil {
			logger.GetLogger().Warn(
				"Could not unmarshal the received data",
				zap.ByteString("Received", received),
				zap.String("Error", err.Error()),
			)
			continue
		}
		fmt.Println(raw)
	}
}
