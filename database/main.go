package database

import (
	"github.com/munnik/gosk/message"
)

type DatabaseWriter interface {
	WriteRaw(raw *message.Raw)
	WriteMapped(mapped *message.Mapped)
}

type DatabaseReader interface {
	ReadRaw(where string, arguments ...interface{}) ([]message.Raw, error)
	ReadMapped(where string, arguments ...interface{}) ([]message.Mapped, error)
}
