// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func beginTest(w *Wx) {
	w.Timestamp(time.Date(2016, time.November, 5, 20, 35, 0, 0, time.UTC))
}

func endTest(w *Wx) {
	w.Clear()
}

func TestWxText(t *testing.T) {
	const pre = "@052035z3542.00N/07842.00W_"
	const suf = "GoTest_9.9-Stn"

	a := assert.New(t)

	SwName = "GoTest"
	SwVers = "9.9"
	w := Wx{
		Lat:  35.7,
		Lon:  -78.7,
		Type: "Stn",
	}

	beginTest(&w)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h..b....."+suf)
	w.WindDirection(180)
	a.Equal(w.String(), pre+"180/...g...t...r...p...P...h..b....."+suf)
	w.WindSpeed(8)
	a.Equal(w.String(), pre+"180/008g...t...r...p...P...h..b....."+suf)
	w.WindGust(16)
	a.Equal(w.String(), pre+"180/008g016t...r...p...P...h..b....."+suf)
	endTest(&w)

	beginTest(&w)
	w.Temperature(72.0)
	a.Equal(w.String(), pre+".../...g...t072r...p...P...h..b....."+suf)
	endTest(&w)

	beginTest(&w)
	w.RainRate(0.54)
	a.Equal(w.String(), pre+".../...g...t...r054p...P...h..b....."+suf)
	w.RainLast24Hours(0.23)
	a.Equal(w.String(), pre+".../...g...t...r054p023P...h..b....."+suf)
	w.RainToday(0.21)
	a.Equal(w.String(), pre+".../...g...t...r054p023P021h..b....."+suf)
	endTest(&w)

	beginTest(&w)
	w.Humidity(61)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h61b....."+suf)
	w.Humidity(100)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h00b....."+suf)
	endTest(&w)

	beginTest(&w)
	w.Altimeter(29.87)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h..b10115"+suf)
	endTest(&w)

	beginTest(&w)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h..b....."+suf)
	w.SolarRadiation(864)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h..b.....L864"+suf)
	w.SolarRadiation(1864)
	a.Equal(w.String(), pre+".../...g...t...r...p...P...h..b.....l864"+suf)
	endTest(&w)
}
