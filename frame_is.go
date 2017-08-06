// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Refer to Automatic Position Reporting System (APRS) Protocol
// Reference - Protocol version 1.0.

package aprs

import "fmt"

// String converts an Address into its text representation.
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

// String converts a Frame into its text representation appropriate
// for printing or sending via APRS-IS.
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
