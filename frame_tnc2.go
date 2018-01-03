// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Refer to Automatic Position Reporting System (APRS) Protocol
// Reference - Protocol version 1.0.

package aprs

import (
	"fmt"
	"regexp"
	"strings"
)

// String returns the Address as a TNC2 formatted string.
func (a Address) String() (addr string) {
	addr = a.Call
	if a.SSID > 0 {
		addr += fmt.Sprintf("-%d", a.SSID)
	}
	if a.Repeated {
		addr += "*"
	}

	return
}

// FromString sets the Frame from a TNC2 formatted string.
//
// This strictly enforces the AX.25 specification and will
// return errors if callsigns are greater than 6 characters or
// SSID's are not numeric values between 0 and 15.
func (f *Frame) FromString(frame string) (err error) {
	// SRC>DST[,PATH]:TEXT
	const reCall = "[[:alnum:]]{1,6}(?:-(?:[0-9]|1[0-5]))?"
	re := regexp.MustCompile(fmt.Sprintf("^(%[1]s)>(%[1]s)((?:(?:,)(?:%[1]s\\*?))*):(.*)", reCall))

	matches := re.FindStringSubmatch(frame)
	if matches == nil {
		err = ErrFrameInvalid
		return
	}

	// Nothing should ever error unless there is a mistake
	// in the regular expression.
	err = f.Src.FromString(matches[1])
	if err != nil {
		return
	}
	err = f.Dst.FromString(matches[2])
	if err != nil {
		return
	}
	err = f.Path.FromString(strings.TrimLeft(matches[3], ","))
	if err != nil {
		return
	}
	f.Text = matches[4]

	return
}

// String returns the Frame as a TNC2 formatted string.  This is
// suitable for sending to APRS-IS servers.
func (f Frame) String() (frame string) {
	// We have to manipulate the Addresses a little because only
	// the last repeated address should have an asterisk.

	// Destination and Source Addresses
	frame = fmt.Sprintf("%s>%s",
		Address{Call: f.Src.Call, SSID: f.Src.SSID},
		Address{Call: f.Dst.Call, SSID: f.Dst.SSID})
	// Path (optional)
	for i := 0; i < len(f.Path); i++ {
		a := f.Path[i]
		// Is there another address in the path and is it repeated?
		if i+1 < len(f.Path) && f.Path[i+1].Repeated {
			a.Repeated = false
		}
		frame += fmt.Sprintf(",%s", a)
	}
	frame += fmt.Sprintf(":%s", f.Text) // Information Field

	return
}
