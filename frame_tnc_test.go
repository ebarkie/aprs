// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrameFromBytesIncomplete(t *testing.T) {
	// Incomplete frames should not crash.
	for _, frame := range [][]byte{tncWx1, tncWx2} {
		for i := 0; i < len(frame); i++ {
			f := Frame{}
			f.FromBytes(frame[0:i])
		}
	}
}

func TestFrameFromBytes(t *testing.T) {
	a := assert.New(t)

	f := Frame{}
	err := f.FromBytes(tncWx2)
	if err != nil {
		a.FailNow(err.Error())
	}

	//t.Logf("%+v", f)

	a.Equal("APX209", f.Dst.Call, "Destination call")
	a.Equal(0, f.Dst.SSID, "Destination SSID")
	a.Equal(true, f.Dst.Repeated, "Destination repeated")
	a.Equal(false, f.Dst.last, "Destination last")

	a.Equal("N4MTT", f.Src.Call, "Source call")
	a.Equal(2, f.Src.SSID, "Source SSID")
	a.Equal(false, f.Src.Repeated, "Source repeated")
	a.Equal(false, f.Src.last, "Source last")

	a.Equal(2, len(f.Path), "Path length")
	a.Equal("KD4PBS", f.Path[0].Call, "Path 0 call")
	a.Equal(3, f.Path[0].SSID, "Path 0 SSID")
	a.Equal(true, f.Path[0].Repeated, "Path 0 repeated")
	a.Equal(false, f.Path[0].last, "Path 0 last")
	a.Equal("WIDE2", f.Path[1].Call, "Path 1 call")
	a.Equal(2, f.Path[1].SSID, "Path 1 SSID")
	a.Equal(false, f.Path[1].Repeated, "Path 1 repeated")
	a.Equal(true, f.Path[1].last, "Path 1 last")

	a.Equal("@270055z3548.41N/07846.35W_360/000g000t066r000P000p000h63b10183XU2k\r", f.Text, "Text")
}

func TestFrameToBytes(t *testing.T) {
	f := Frame{}
	f.Src = Address{Call: "N4MTT", SSID: 2}
	f.Dst = Address{Call: "APX209", Repeated: true}
	f.Path = []Address{{Call: "KD4PBS", SSID: 3, Repeated: true}, {Call: "WIDE2", SSID: 2}}
	f.Text = "@270055z3548.41N/07846.35W_360/000g000t066r000P000p000h63b10183XU2k\r"

	a := assert.New(t)
	a.Equal(tncWx2, f.Bytes(), "TNC encoded frame")
}
