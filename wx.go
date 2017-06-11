// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"fmt"
	"math"
	"time"
)

// Wx holds the weather station information.
type Wx struct {
	Lat  float64
	Lon  float64
	Type string
	o    obs
}

// obs represents an observation to be encoded.
type obs struct {
	altimeter       *float64
	humidity        *int
	rainRate        *float64
	rainLast24Hours *float64
	rainToday       *float64
	solarRadiation  *int
	temperature     *int
	timestamp       time.Time
	windDirection   *int
	windGust        *int
	windSpeed       *int
}

// Altimeter is the barometric pressure as altimeter
// (station plus elevation correction) in inches.
func (w *Wx) Altimeter(f float64) {
	// Specification calls for millibars - convert in to mb.
	v := f * 33.8637526
	w.o.altimeter = &v
}

// Humidity is the relative humidity.
func (w *Wx) Humidity(i int) {
	// Specification states that 100 should be represented as 0.
	v := i
	if v == 100 {
		v = 0
	}
	w.o.humidity = &v
}

// RainRate is the rain accumulation for the last hour
// in inches.
func (w *Wx) RainRate(f float64) {
	w.o.rainRate = &f
}

// RainLast24Hours is the rain accumulation for the last 24 hours
// in inches.
func (w *Wx) RainLast24Hours(f float64) {
	w.o.rainLast24Hours = &f
}

// RainToday is the rain accumulation so far today in inches.
func (w *Wx) RainToday(f float64) {
	w.o.rainToday = &f
}

// SolarRadiation is the solar radiation in W/m^2.
func (w *Wx) SolarRadiation(i int) {
	w.o.solarRadiation = &i
}

// Temperature is the temperature in F.
func (w *Wx) Temperature(f float64) {
	v := int(f)
	w.o.temperature = &v
}

// Timestamp is the time of the observation.
func (w *Wx) Timestamp(t time.Time) {
	w.o.timestamp = t
}

// WindDirection is the wind direction in degrees.
func (w *Wx) WindDirection(i int) {
	w.o.windDirection = &i
}

// WindGust is the peak wind speed for the previous 5 minutes
// in mph.
func (w *Wx) WindGust(i int) {
	w.o.windGust = &i
}

// WindSpeed is the wind speed in mph.
func (w *Wx) WindSpeed(i int) {
	w.o.windSpeed = &i
}

// Clear clears all measurements in the observation payload.
func (w *Wx) Clear() {
	w.o = obs{}
}

// llDecToDMS takes a float latitude or longitude and converts it to
// degrees, minutes (second as decimal), and a hemisphere string.
func llDecToDMS(ll float64, hems [2]string) (float64, float64, string) {
	deg, frac := math.Modf(math.Abs(ll))
	min := frac * 60.0
	hem := hems[0]
	if ll < 0 {
		hem = hems[1]
	}

	return deg, min, hem
}

// String returns an APRS packet for the provided measurements.
func (w Wx) String() (t string) {
	// Refer to APRS Protocol Reference 1.0
	// Chapter 12: Weather Reports
	//
	// Parameters:
	//   c = wind direction (in degrees) [4 bytes]
	//   s = sustained one-minute wind speed (in mph) [4 bytes]
	//   g = gust (peak wind speed in mph in the last 5 minutes) [4 bytes]
	//   t = temperature (in degrees Fahrenheit). Temperatures below zero
	//       are expressed as -01 to -99 [4 bytes]
	//   r = rainfall (in hundredths of an inch) in the last hour [4 bytes]
	//   p = rainfall (in hundredths of an inch) in the last 24 hours [4 bytes]
	//   P = rainfall (in hundredths of an inch) since midnight [4 bytes]
	//   h = humidity (in %. 00 = 100%) [3 bytes]
	//   b = barometric altimeter pressure (in tenths of millibars/tenths of
	//       hPascal) [5 bytes]
	//
	// Optional parameters:
	//   L = luminosity (in watts per square meter) 999 and below [4 bytes]
	//   l = luminosity (in watts per square meter) 1000 and above [4 bytes]
	//   (L is inserted in place of one of the rain values)
	//   s = snowfall (in inches) in the last 24 hours
	//   # = raw rain counter

	// Timestamp as UTC two digit day, hour, and minute.
	if w.o.timestamp.IsZero() {
		w.o.timestamp = time.Now()
	}
	t += fmt.Sprintf("@%sz", w.o.timestamp.In(time.UTC).Format("021504"))

	// Location
	deg, min, hem := llDecToDMS(w.Lat, [2]string{"N", "S"})
	t += fmt.Sprintf("%02.0f%02.2f%s", deg, min, hem)
	deg, min, hem = llDecToDMS(w.Lon, [2]string{"E", "W"})
	t += fmt.Sprintf("/%03.0f%02.2f%s", deg, min, hem)

	// Parameters
	if w.o.windDirection == nil {
		t += "_..."
	} else {
		t += fmt.Sprintf("_%03d", *w.o.windDirection)
	}

	if w.o.windSpeed == nil {
		t += "/..."
	} else {
		t += fmt.Sprintf("/%03d", *w.o.windSpeed)
	}

	if w.o.windGust == nil {
		t += "g..."
	} else {
		t += fmt.Sprintf("g%03d", *w.o.windGust)
	}

	if w.o.temperature == nil {
		t += "t..."
	} else {
		t += fmt.Sprintf("t%03d", *w.o.temperature)
	}

	if w.o.rainRate == nil {
		t += "r..."
	} else {
		t += fmt.Sprintf("r%03.0f", *w.o.rainRate*100.0)
	}

	if w.o.rainLast24Hours == nil {
		t += "p..."
	} else {
		t += fmt.Sprintf("p%03.0f", *w.o.rainLast24Hours*100.0)
	}

	if w.o.rainToday == nil {
		t += "P..."
	} else {
		t += fmt.Sprintf("P%03.0f", *w.o.rainToday*100.0)
	}

	if w.o.humidity == nil {
		t += "h.."
	} else {
		t += fmt.Sprintf("h%02d", *w.o.humidity)
	}

	if w.o.altimeter == nil {
		t += "b....."
	} else {
		t += fmt.Sprintf("b%05.0f", *w.o.altimeter*10.0)
	}

	if w.o.solarRadiation != nil {
		if *w.o.solarRadiation >= 1000 {
			t += fmt.Sprintf("l%03d", *w.o.solarRadiation-1000)
		} else {
			t += fmt.Sprintf("L%03d", *w.o.solarRadiation)
		}
	}

	// Software and weather unit type
	t += fmt.Sprintf("%s_%s", SwName, SwVers)
	if w.Type != "" {
		t += fmt.Sprintf("-%s", w.Type)
	}

	return
}
