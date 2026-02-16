// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// KG4HIE>APK102,W4LBT-9,WIDE1,KD4PBS-3*,WIDE2:=3438.51N/07941.15W_120/001g004t073r   p   P000h  b     KU2k<0x0d>
// U frame UI: p/f=0, No layer 3 protocol implemented., length = 105
//
//	dest    APK102  0 c/r=0 res=3 last=0
//	source  KG4HIE  0 c/r=1 res=3 last=0
//	digi 1  W4LBT   9   h=1 res=3 last=0
//	digi 2  WIDE1   0   h=1 res=3 last=0
//	digi 3  KD4PBS  3   h=1 res=3 last=0
//	digi 4  WIDE2   0   h=0 res=3 last=1
//
// Weather Report, WEATHER Station (blue), Kenwood D710
// N 34 38.5100, W 079 41.1500
// wind 1.2 mph, direction 120, gust 4, temperature 73, rain 0.00 since midnight, "KU2k"
var ax25Wx1 = []byte{
	0x82, 0xa0, 0x96, 0x62, 0x60, 0x64, 0x60, 0x96,
	0x8e, 0x68, 0x90, 0x92, 0x8a, 0xe0, 0xae, 0x68,
	0x98, 0x84, 0xa8, 0x40, 0xf2, 0xae, 0x92, 0x88,
	0x8a, 0x62, 0x40, 0xe0, 0x96, 0x88, 0x68, 0xa0,
	0x84, 0xa6, 0xe6, 0xae, 0x92, 0x88, 0x8a, 0x64,
	0x40, 0x61, 0x03, 0xf0, 0x3d, 0x33, 0x34, 0x33,
	0x38, 0x2e, 0x35, 0x31, 0x4e, 0x2f, 0x30, 0x37,
	0x39, 0x34, 0x31, 0x2e, 0x31, 0x35, 0x57, 0x5f,
	0x31, 0x32, 0x30, 0x2f, 0x30, 0x30, 0x31, 0x67,
	0x30, 0x30, 0x34, 0x74, 0x30, 0x37, 0x33, 0x72,
	0x20, 0x20, 0x20, 0x70, 0x20, 0x20, 0x20, 0x50,
	0x30, 0x30, 0x30, 0x68, 0x20, 0x20, 0x62, 0x20,
	0x20, 0x20, 0x20, 0x20, 0x4b, 0x55, 0x32, 0x6b,
	0x0d,
}

// N4MTT-2>APX209,KD4PBS-3*,WIDE2-2:@270055z3548.41N/07846.35W_360/000g000t066r000P000p000h63b10183XU2k<0x0d>
// U frame UI: p/f=0, No layer 3 protocol implemented., length = 98
//
//	dest    APX209  0 c/r=1 res=3 last=0
//	source  N4MTT   2 c/r=0 res=3 last=0
//	digi 1  KD4PBS  3   h=1 res=3 last=0
//	digi 2  WIDE2   2   h=0 res=3 last=1
//
// Weather Report, WEATHER Station (blue), Xastir
// N 35 48.4100, W 078 46.3500
// wind 0.0 mph, direction 360, gust 0, temperature 66, rain 0.00 in last hour, rain 0.00 since midnight, rain 0.00 in last 24 hours, humidity 63, barometer 30.07, "XU2k"
var ax25Wx2 = []byte{
	0x82, 0xa0, 0xb0, 0x64, 0x60, 0x72, 0xe0, 0x9c,
	0x68, 0x9a, 0xa8, 0xa8, 0x40, 0x64, 0x96, 0x88,
	0x68, 0xa0, 0x84, 0xa6, 0xe6, 0xae, 0x92, 0x88,
	0x8a, 0x64, 0x40, 0x65, 0x03, 0xf0, 0x40, 0x32,
	0x37, 0x30, 0x30, 0x35, 0x35, 0x7a, 0x33, 0x35,
	0x34, 0x38, 0x2e, 0x34, 0x31, 0x4e, 0x2f, 0x30,
	0x37, 0x38, 0x34, 0x36, 0x2e, 0x33, 0x35, 0x57,
	0x5f, 0x33, 0x36, 0x30, 0x2f, 0x30, 0x30, 0x30,
	0x67, 0x30, 0x30, 0x30, 0x74, 0x30, 0x36, 0x36,
	0x72, 0x30, 0x30, 0x30, 0x50, 0x30, 0x30, 0x30,
	0x70, 0x30, 0x30, 0x30, 0x68, 0x36, 0x33, 0x62,
	0x31, 0x30, 0x31, 0x38, 0x33, 0x58, 0x55, 0x32,
	0x6b, 0x0d,
}

func TestFrameFromBytesIncomplete(t *testing.T) {
	// Incomplete frames should not crash.
	for _, frame := range [][]byte{ax25Wx1, ax25Wx2} {
		for i := 0; i < len(frame); i++ {
			f := Frame{}
			f.FromBytes(frame[0:i])
		}
	}
}

func TestFrameFromBytes(t *testing.T) {
	a := assert.New(t)

	f := Frame{}
	err := f.FromBytes(ax25Wx2)
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
	f.Src = Addr{Call: "N4MTT", SSID: 2}
	f.Dst = Addr{Call: "APX209", Repeated: true}
	f.Path = []Addr{{Call: "KD4PBS", SSID: 3, Repeated: true}, {Call: "WIDE2", SSID: 2}}
	f.Text = "@270055z3548.41N/07846.35W_360/000g000t066r000P000p000h63b10183XU2k\r"

	a := assert.New(t)
	a.Equal(ax25Wx2, f.Bytes(), "TNC encoded frame")
}
