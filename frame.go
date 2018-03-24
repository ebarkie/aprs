// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"fmt"
	"strconv"
	"strings"
)

// Frame represents a complete APRS frame.
type Frame struct {
	Dst  Addr
	Src  Addr
	Path Path
	Text string
}

// Addr represents an APRS callsign, SSID, and associated
// metadata.
type Addr struct {
	SSID     int
	Repeated bool
	last     bool
	Call     string
}

// Path represents the APRS digipath.
type Path []Addr

// FromString sets the address from a string.
func (a *Addr) FromString(addr string) (err error) {
	if strings.HasSuffix(addr, "*") {
		a.Repeated = true
		addr = addr[:len(addr)-1]
	}

	dash := strings.Index(addr, "-")
	if dash > -1 {
		a.Call = addr[:dash]
		a.SSID, err = strconv.Atoi(addr[dash+1:])
		if err != nil {
			err = fmt.Errorf("Address error: SSID is invalid: %s", err.Error())
			return
		}
	} else {
		a.Call = addr
	}

	if len(a.Call) > 6 {
		err = fmt.Errorf("Address error: Callsign length %d > 6", len(a.Call))
		return
	}
	if a.SSID < 0 || a.SSID > 15 {
		err = fmt.Errorf("Address error: %d not > 0 & < 15", a.SSID)
		return
	}

	return
}

// FromString sets the Path from a string of comma separated
// addresses.
func (p *Path) FromString(path string) (err error) {
	for _, as := range strings.Split(path, ",") {
		a := Addr{}
		err = a.FromString(as)
		if err != nil {
			return
		}
		*p = append(*p, a)
	}

	return
}
