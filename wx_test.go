// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"fmt"
	"time"
)

var testWx = Wx{
	Lat:  35.7,
	Lon:  -78.7,
	Type: "Stn",
}

func init() {
	testWx.Zero()
	testWx.Timestamp = time.Date(2016, time.November, 5, 20, 35, 0, 0, time.UTC)
}

func ExampleWx_String() {
	w := testWx
	fmt.Println(w)

	SwName = "GoTst"
	SwVers = "9"
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h..b.....Go3-Stn
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h..b.....GoTst9-Stn
}

func ExampleWx_String_pressure() {
	w := testWx

	w.Pressure = 29.87
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h..b10115GoTst9-Stn
}

func ExampleWx_String_humidity() {
	w := testWx

	w.Humidity = 61
	fmt.Println(w)

	w.Humidity = 100
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h61b.....GoTst9-Stn
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h00b.....GoTst9-Stn
}

func ExampleWx_String_luminosity() {
	w := testWx

	w.SolarRad = 864
	fmt.Println(w)

	w.SolarRad = 1864
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h..b.....L864GoTst9-Stn
	// @052035z3542.00N/07842.00W_.../...g...t...r...p...P...h..b.....l864GoTst9-Stn
}

func ExampleWx_String_rain() {
	w := testWx

	w.RainLastHour = 0.0
	w.RainLast24Hours = 0.0
	w.RainToday = 0.0
	fmt.Println(w)

	w.RainLastHour = 0.54
	w.RainLast24Hours = 0.23
	w.RainToday = 0.21
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_.../...g...t...r000p000P000h..b.....GoTst9-Stn
	// @052035z3542.00N/07842.00W_.../...g...t...r054p023P021h..b.....GoTst9-Stn
}

func ExampleWx_String_temp() {
	w := testWx

	w.Temp = -20
	fmt.Println(w)

	w.Temp = 0
	fmt.Println(w)

	w.Temp = 72
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_.../...g...t-20r...p...P...h..b.....GoTst9-Stn
	// @052035z3542.00N/07842.00W_.../...g...t000r...p...P...h..b.....GoTst9-Stn
	// @052035z3542.00N/07842.00W_.../...g...t072r...p...P...h..b.....GoTst9-Stn
}

func ExampleWx_String_wind() {
	w := testWx

	w.WindDir = 0
	w.WindSpeed = 0
	w.WindGust = 0
	fmt.Println(w)

	w.WindDir = 180
	w.WindSpeed = 8
	w.WindGust = 16
	fmt.Println(w)

	// Output:
	// @052035z3542.00N/07842.00W_000/000g000t...r...p...P...h..b.....GoTst9-Stn
	// @052035z3542.00N/07842.00W_180/008g016t...r...p...P...h..b.....GoTst9-Stn
}
