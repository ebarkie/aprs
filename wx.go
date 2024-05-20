// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"fmt"
	"time"
)

// Wx represents a weather station observation.
type Wx struct {
	Lat  float64
	Lon  float64
	Type string

	Timestamp time.Time

	Humidity        int     // 0-100 %
	Pressure        float64 // mBar
	RainLastHour    float64
	RainLast24Hours float64
	RainToday       float64
	SolarRad        int
	Temp            int
	WindDir         int
	WindGust        int
	WindSpeed       int
}

// Zero zeroes all measurements in the observation payload.
func (w *Wx) Zero() {
	w.Timestamp = time.Time{}
	w.Humidity = -1
	w.Pressure = 0
	w.RainLastHour = -1.0
	w.RainLast24Hours = -1.0
	w.RainToday = -1.0
	w.SolarRad = -1
	w.Temp = -100
	w.WindDir = -1
	w.WindGust = -1
	w.WindSpeed = -1
}

// String returns an APRS packet for the provided measurements.
func (w Wx) String() (s string) {
	// Refer to APRS Protocol Reference 1.0
	// Chapter 12: Weather Reports
	//
	// Parameters:
	//   _ = wind direction (in degrees) [3 chars]
	//   / = sustained one-minute wind speed (in mph) [3 chars]
	//   g = gust (peak wind speed in mph in the last 5 minutes) [3 chars]
	//   t = temperature (in degrees Fahrenheit). Temperatures below zero
	//       are expressed as -01 to -99 [3 chars]
	//   r = rainfall (in hundredths of an inch) in the last hour [3 chars]
	//   p = rainfall (in hundredths of an inch) in the last 24 hours [3 chars]
	//   P = rainfall (in hundredths of an inch) since midnight [3 chars]
	//   h = humidity (in %. 00 = 100%) [2 chars]
	//   b = barometric altimeter pressure (in tenths of millibars/tenths of
	//       hPascal) [4 chars]
	//
	//   L = luminosity (in watts per square meter) 999 and below [3 chars]
	//   l = luminosity (in watts per square meter) 1000 and above [3 chars]
	//   (L is inserted in place of one of the rain values)
	//   s = snowfall (in inches) in the last 24 hours
	//   # = raw rain counter

	// Timestamp as UTC two digit day, hour, and minute.
	if w.Timestamp.IsZero() {
		w.Timestamp = time.Now()
	}

	// Base prefix
	latDeg, latMin, latHem := decToDMS(w.Lat, [2]string{"N", "S"})
	lonDeg, lonMin, lonHem := decToDMS(w.Lon, [2]string{"E", "W"})
	s = fmt.Sprintf("@%sz%02.0f%05.2f%s/%03.0f%05.2f%s",
		w.Timestamp.In(time.UTC).Format("021504"),
		latDeg, latMin, latHem,
		lonDeg, lonMin, lonHem)

	// Parameters
	if w.WindDir < 0 {
		s += "_..."
	} else {
		s += fmt.Sprintf("_%03d", w.WindDir)
	}

	if w.WindSpeed < 0 {
		s += "/..."
	} else {
		s += fmt.Sprintf("/%03d", w.WindSpeed)
	}

	if w.WindGust < 0 {
		s += "g..."
	} else {
		s += fmt.Sprintf("g%03d", w.WindGust)
	}

	if w.Temp < -99 {
		s += "t..."
	} else {
		s += fmt.Sprintf("t%03d", w.Temp)
	}

	if w.RainLastHour < 0.0 {
		s += "r..."
	} else {
		s += fmt.Sprintf("r%03.0f", w.RainLastHour*100.0)
	}

	if w.RainLast24Hours < 0.0 {
		s += "p..."
	} else {
		s += fmt.Sprintf("p%03.0f", w.RainLast24Hours*100.0)
	}

	if w.RainToday < 0.0 {
		s += "P..."
	} else {
		s += fmt.Sprintf("P%03.0f", w.RainToday*100.0)
	}

	if w.Humidity < 0 {
		s += "h.."
	} else {
		s += fmt.Sprintf("h%02d", w.Humidity%100)
	}

	if w.Pressure <= 0.0 {
		s += "b....."
	} else {
		s += fmt.Sprintf("b%05.0f", w.Pressure*10.0)
	}

	if w.SolarRad >= 1000 {
		s += fmt.Sprintf("l%03d", w.SolarRad-1000)
	} else if w.SolarRad >= 0 {
		s += fmt.Sprintf("L%03d", w.SolarRad)
	}

	// Software
	if w.Type != "" {
		s += w.Type
	} else {
		s += "GolangAPRS"
	}

	return
}
