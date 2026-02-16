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

	before, after, ok := strings.Cut(addr, "-")
	if ok {
		a.Call = before
		a.SSID, err = strconv.Atoi(after)
		if err != nil {
			err = fmt.Errorf("address error: SSID is invalid: %s", err.Error())
			return
		}
	} else {
		a.Call = addr
	}

	if len(a.Call) > 6 {
		err = fmt.Errorf("address error: Callsign length %d > 6", len(a.Call))
		return
	}
	if a.SSID < 0 || a.SSID > 15 {
		err = fmt.Errorf("address error: SSID %d not in range 0-15", a.SSID)
		return
	}

	return
}

// FromString sets the Path from a string of comma separated
// addresses.
func (p *Path) FromString(path string) (err error) {
	if path == "" {
		return
	}
	for as := range strings.SplitSeq(path, ",") {
		a := Addr{}
		err = a.FromString(as)
		if err != nil {
			return
		}
		*p = append(*p, a)
	}

	return
}
