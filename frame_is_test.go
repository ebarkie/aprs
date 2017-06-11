// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressString(t *testing.T) {
	a := Address{Call: "N0CALL", SSID: 13}
	assert.Equal(t, "N0CALL-13", a.String(), "Address")
}

func TestFrameString(t *testing.T) {
	a := assert.New(t)

	f := Frame{}
	f.Src.FromString("N0CALL-13")
	f.Dst.FromString("APZ001")
	f.Path.FromString("WIDE1-1,WIDE2-1")
	f.Text = "Hello world"
	a.Equal("N0CALL-13>APZ001,WIDE1-1,WIDE2-1:Hello world",
		f.String())

	f = Frame{}
	f.FromBytes(tncWx1)
	a.Equal("KG4HIE>APK102,W4LBT-9,WIDE1,KD4PBS-3*,WIDE2:=3438.51N/07941.15W_120/001g004t073r   p   P000h  b     KU2k\r",
		f.String())

	f = Frame{}
	f.FromBytes(tncWx2)
	a.Equal("N4MTT-2>APX209,KD4PBS-3*,WIDE2-2:@270055z3548.41N/07846.35W_360/000g000t066r000P000p000h63b10183XU2k\r",
		f.String())
}
