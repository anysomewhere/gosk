package mapper

import (
	"fmt"
	"github.com/boatkit-io/n2k/pkg/adapter/canadapter"
	"github.com/boatkit-io/n2k/pkg/pgn"
	"github.com/boatkit-io/n2k/pkg/pkt"
	"github.com/brutella/can"
	"github.com/munnik/gosk/config"
	"github.com/munnik/gosk/logger"
	"github.com/munnik/gosk/message"
	"go.nanomsg.org/mangos/v3"
	"go.uber.org/zap"
)

type Nmea2000Mapper struct {
	config   config.MapperConfig
	protocol string
	multi    *canadapter.MultiBuilder
	current  *pkt.Packet
}

func NewNmea2000Mapper(c config.MapperConfig) (*Nmea2000Mapper, error) {
	return &Nmea2000Mapper{config: c, protocol: config.NMEA2000Type, multi: canadapter.NewMultiBuilder(nil)}, nil
}

func (m *Nmea2000Mapper) Map(subscriber mangos.Socket, publisher mangos.Socket) {
	process(subscriber, publisher, m)
}

func (m *Nmea2000Mapper) DoMap(raw *message.Raw) (*message.Mapped, error) {
	frm := can.Frame{}
	err := can.Unmarshal(raw.Value, &frm)
	if err != nil {
		return nil, err
	}

	pktInfo := canadapter.NewPacketInfo(canadapter.Frame(frm))
	m.current = pkt.NewPacket(pktInfo, frm.Data[:])
	packet := m.process()

	if packet == nil {
		return nil, nil
	}

	if len(packet.Decoders) > 0 {
		// call frame decoders, send valid return on.
		for _, decoder := range packet.Decoders {
			stream := pgn.NewPgnDataStream(packet.Data)
			ret, err := decoder(packet.Info, stream)
			if err != nil {
				packet.ParseErrors = append(packet.ParseErrors, err)
				continue
			} else {
				logger.GetLogger().Info("Decoded packet", zap.Any("packet", ret))
				break
			}
		}
	} else {
		// No valid decoder, so send on an UnknownPGN.
		packet.ParseErrors = append(packet.ParseErrors, fmt.Errorf("no matching decoder"))
	}

	return nil, nil
}

// process method is the worker function for Run
func (m *Nmea2000Mapper) process() *pkt.Packet {
	// https://endige.com/2050/nmea-2000-pgns-deciphered/

	if len(m.current.ParseErrors) > 0 {
		return m.current
	}

	if m.current.Info.PGN == 130824 {
		m.handlePGN130824()
	}
	if m.current.Fast {
		m.multi.Add(m.current)
	} else {
		m.current.Complete = true
	}
	if m.current.Complete {
		m.current.AddDecoders()
		return m.current
	}

	return nil
}

// handlePGN130824 method deals with the only PGN that has both a single and fast variant.
// PGN 130824 has 1 fast and 1 slow variant; all other PGNs are one or the other.
//
//	We validate this invariant on every import of canboat.json.
//
// If it changes we need to revisit this code.
// The slow variant starts 0x7D 0x81 (man code and industry). We'll look for this and if matched select
// it. The fast variant is length 9, fitting in 2 frames, so the first byte of either frame
// can't be 0x7D
func (m *Nmea2000Mapper) handlePGN130824() {
	var pInfo *pgn.PgnInfo
	m.current.Fast = true      // if slow match fails, the normal code will process this
	m.current.Complete = false //
	m.current.GetManCode()     // have to peak ahead in this special case
	for _, pInfo = range m.current.Candidates {
		if pInfo.Fast {
			break // we only want to check the slow variant here
		} else {
			if m.current.Manufacturer == pInfo.ManId {
				m.current.Fast = false
				m.current.Complete = true
				break
			}
		}
	}
}
