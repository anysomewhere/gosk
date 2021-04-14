package nmea

import (
	"fmt"

	goNMEA "github.com/adrianmo/go-nmea"
	"github.com/martinlindhe/unit"
)

type VHW struct {
	goNMEA.BaseSentence
	TrueHeading            Float64
	MagneticHeading        Float64
	SpeedThroughWaterKnots Float64
	SpeedThroughWaterKPH   Float64
}

func init() {
	goNMEA.MustRegisterParser("VHW", func(s goNMEA.BaseSentence) (goNMEA.Sentence, error) {
		p := goNMEA.NewParser(s)
		result := VHW{
			BaseSentence: s,
		}
		if p.Fields[0] != "" {
			result.TrueHeading = NewFloat64(WithValue(p.Float64(0, "true heading")))
		} else {
			result.TrueHeading = NewFloat64()
		}
		if p.Fields[2] != "" {
			result.MagneticHeading = NewFloat64(WithValue(p.Float64(2, "magnetic heading")))
		} else {
			result.MagneticHeading = NewFloat64()
		}
		if p.Fields[4] != "" {
			result.SpeedThroughWaterKnots = NewFloat64(WithValue(p.Float64(4, "speed through water in knots")))
		} else {
			result.SpeedThroughWaterKnots = NewFloat64()
		}
		if p.Fields[6] != "" {
			result.SpeedThroughWaterKPH = NewFloat64(WithValue(p.Float64(6, "speed through water in kilometers per hour")))
		} else {
			result.SpeedThroughWaterKPH = NewFloat64()
		}
		return result, p.Err()
	})
}

// GetMagneticHeading retrieves the magnetic heading from the sentence
func (s VHW) GetMagneticHeading() (float64, error) {
	if !s.MagneticHeading.isNil {
		return (unit.Angle(s.MagneticHeading.value) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetTrueHeading retrieves the true heading from the sentence
func (s VHW) GetTrueHeading() (float64, error) {
	if !s.TrueHeading.isNil {
		return (unit.Angle(s.TrueHeading.value) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetSpeedThroughWater retrieves the speed through water from the sentence
func (s VHW) GetSpeedThroughWater() (float64, error) {
	if !s.SpeedThroughWaterKPH.isNil {
		return (unit.Speed(s.SpeedThroughWaterKPH.value) * unit.KilometersPerHour).MetersPerSecond(), nil
	}
	if !s.SpeedThroughWaterKnots.isNil {
		return (unit.Speed(s.SpeedThroughWaterKnots.value) * unit.Knot).MetersPerSecond(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
