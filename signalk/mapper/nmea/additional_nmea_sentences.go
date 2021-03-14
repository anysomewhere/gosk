package nmea

import (
	log "github.com/sirupsen/logrus"

	"github.com/adrianmo/go-nmea"
)

const (
	// TypeMDA for MDA messages
	TypeMDA = "MDA"
	// TypeMWV for MWV messages
	TypeMWV = "MWV"
	// TypeVWR for VWR messages
	TypeVWR = "VWR"
)

// MDA Meteorological Composite
type MDA struct {
	nmea.BaseSentence
	BarometricPressureInInchesOfMercury float64
	BarometricPressureInBar             float64
	AirTemperature                      float64
	WaterTemperature                    float64
	RelativeHumidity                    float64
	AbsoluteHumidity                    float64
	DewPoint                            float64
	WindDirectionTrue                   float64
	WindDirectionMagnetic               float64
	WindSpeedInKnots                    float64
	WindSpeedInMetersPerSecond          float64
}

// MWV Wind Speed and Angle
type MWV struct {
	nmea.BaseSentence
	Angle         float64
	Reference     string
	WindSpeed     float64
	WindSpeedUnit string
	Status        string
}

// VWR Relative Wind Speed and Angle
type VWR struct {
	nmea.BaseSentence
	Direction                    float64
	LeftRightOfBow               string
	WindSpeedInKnots             float64
	WindSpeedInMetersPerSecond   float64
	WindSpeedInKilometersPerHour float64
}

func init() {
	if err := nmea.RegisterParser("MDA", func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		p.AssertType(TypeMDA)
		return MDA{
			BaseSentence:                        s,
			BarometricPressureInInchesOfMercury: p.Float64(0, "BarometricPressureInInchesOfMercury"),
			BarometricPressureInBar:             p.Float64(2, "BarometricPressureInBar"),
			AirTemperature:                      p.Float64(4, "AirTemperature"),
			WaterTemperature:                    p.Float64(6, "WaterTemperature"),
			RelativeHumidity:                    p.Float64(8, "RelativeHumidity"),
			AbsoluteHumidity:                    p.Float64(9, "AbsoluteHumidity"),
			DewPoint:                            p.Float64(10, "DewPoint"),
			WindDirectionTrue:                   p.Float64(12, "WindDirectionTrue"),
			WindDirectionMagnetic:               p.Float64(14, "WindDirectionMagnetic"),
			WindSpeedInKnots:                    p.Float64(16, "WindSpeedInKnots"),
			WindSpeedInMetersPerSecond:          p.Float64(18, "WindSpeedInMetersPerSecond"),
		}, p.Err()
	}); err != nil {
		log.Panic(err)
	}
	if err := nmea.RegisterParser("MWV", func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		p.AssertType(TypeMWV)
		return MWV{
			BaseSentence:  s,
			Angle:         p.Float64(0, "Angle"),
			Reference:     p.String(1, "Reference"),
			WindSpeed:     p.Float64(2, "WindSpeed"),
			WindSpeedUnit: p.String(3, "WindSpeedUnit"),
			Status:        p.String(4, "Status"),
		}, p.Err()
	}); err != nil {
		log.Panic(err)
	}
	if err := nmea.RegisterParser("VWR", func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		p.AssertType(TypeVWR)
		return VWR{
			BaseSentence:                 s,
			Direction:                    p.Float64(0, "Direction"),
			LeftRightOfBow:               p.String(1, "LeftRightOfBow"),
			WindSpeedInKnots:             p.Float64(2, "SpeedInKnots"),
			WindSpeedInMetersPerSecond:   p.Float64(4, "SpeedInMetersPerSecond"),
			WindSpeedInKilometersPerHour: p.Float64(6, "SpeedInKilometersPerHour"),
		}, p.Err()
	}); err != nil {
		log.Panic(err)
	}

}