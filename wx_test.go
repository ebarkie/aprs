// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWxText(t *testing.T) {
	a := assert.New(t)

	w := Wx{
		Lat:  35.7,
		Lon:  -78.7,
		Type: "Stn",
	}
	w.Zero()
	w.Timestamp = time.Date(2016, time.November, 5, 20, 35, 0, 0, time.UTC)

	p := "@052035z3542.00N/07842.00W_"
	s := SwName + SwVers + "-Stn"

	// Zero state
	a.Equal(p+".../...g...t...r...p...P...h..b....."+s, w.String())

	// Custom software
	SwName = "GoTst"
	SwVers = "9"
	s = "GoTst9-Stn"
	a.Equal(p+".../...g...t...r...p...P...h..b....."+s, w.String())

	// Parameters
	w.WindDir = 0
	a.Equal(p+"000/...g...t...r...p...P...h..b....."+s, w.String())
	w.WindDir = 180
	a.Equal(p+"180/...g...t...r...p...P...h..b....."+s, w.String())
	w.WindSpeed = 0
	a.Equal(p+"180/000g...t...r...p...P...h..b....."+s, w.String())
	w.WindSpeed = 8
	a.Equal(p+"180/008g...t...r...p...P...h..b....."+s, w.String())
	w.WindGust = 0
	a.Equal(p+"180/008g...t...r...p...P...h..b....."+s, w.String())
	w.WindGust = 16
	a.Equal(p+"180/008g016t...r...p...P...h..b....."+s, w.String())

	w.Temp = -20
	a.Equal(p+"180/008g016t-20r...p...P...h..b....."+s, w.String())
	w.Temp = 0
	a.Equal(p+"180/008g016t000r...p...P...h..b....."+s, w.String())
	w.Temp = 72
	a.Equal(p+"180/008g016t072r...p...P...h..b....."+s, w.String())

	w.RainRate = 0.0
	a.Equal(p+"180/008g016t072r000p...P...h..b....."+s, w.String())
	w.RainRate = 0.54
	a.Equal(p+"180/008g016t072r054p...P...h..b....."+s, w.String())
	w.RainLast24Hours = 0
	a.Equal(p+"180/008g016t072r054p000P...h..b....."+s, w.String())
	w.RainLast24Hours = 0.23
	a.Equal(p+"180/008g016t072r054p023P...h..b....."+s, w.String())
	w.RainToday = 0.0
	a.Equal(p+"180/008g016t072r054p023P000h..b....."+s, w.String())
	w.RainToday = 0.21
	a.Equal(p+"180/008g016t072r054p023P021h..b....."+s, w.String())

	w.Humidity = 61
	a.Equal(p+"180/008g016t072r054p023P021h61b....."+s, w.String())
	w.Humidity = 100
	a.Equal(p+"180/008g016t072r054p023P021h00b....."+s, w.String())

	w.Altimeter = 29.87
	a.Equal(p+"180/008g016t072r054p023P021h00b10115"+s, w.String())

	a.Equal(p+"180/008g016t072r054p023P021h00b10115"+s, w.String())
	w.SolarRad = 864
	a.Equal(p+"180/008g016t072r054p023P021h00b10115L864"+s, w.String())
	w.SolarRad = 1864
	a.Equal(p+"180/008g016t072r054p023P021h00b10115l864"+s, w.String())
}
