// Copyright (c) 2016 Eric Barkie. All rights reserved.  // Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"fmt"
)

// Freq contains an APRS FreqSpec frequency report
// http://www.aprs.org/info/freqspec.txt
type Freq struct {
	Mhz    float64 // Frequency in Mhz
	Tone   int     // Tone in hz
	CTCSS  int     // ctcss (mutually exclusive with tone and dcs)
	DCS    int     // dcs (mutually exclusive with tone and ctcss)
	Offset int     // +/- offset in mhz
	Range  int     // range in miles
	Narrow bool    // defaults to false for wideband
}

// Render renders the frequency struct into a byte-slice
func (f *Freq) Render() string {
	// set the frequency
	out := fmt.Sprintf("%07.03fMHz ", f.Mhz)

	// check for a tone
	if f.Tone > 0 {
		t := "T"
		if f.Narrow {
			t = "t"
		}
		out = fmt.Sprintf("%s%s%s ", out, t, z3p(f.Tone))
	}

	// check for CTCSS
	if f.CTCSS > 0 {
		c := "C"
		if f.Narrow {
			c = "c"
		}
		out = fmt.Sprintf("%s%s%s ", out, c, z3p(f.CTCSS))
	}

	// check for  DCS
	if f.DCS > 0 {
		d := "D"
		if f.Narrow {
			d = "d"
		}
		out = fmt.Sprintf("%s%s%s ", out, d, z3p(f.DCS))
	}

	// check for offset
	if f.Offset != 0 {
		out = fmt.Sprintf("%s%s ", out, z3p(f.Offset))
	}

	// check range
	if f.Range != 0 {
		out = fmt.Sprintf("%sR%sm ", out, z3p(f.Range))
	}
	return out
}
