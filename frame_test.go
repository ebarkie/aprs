// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressFromString(t *testing.T) {
	a := assert.New(t)

	addr := Address{}
	err := addr.FromString("N0CALL-13")
	a.Nil(err, "Valid from string")
	a.Equal("N0CALL", addr.Call, "Call")
	a.Equal(13, addr.SSID, "SSID")

	addr = Address{}
	err = addr.FromString("N0NE")
	a.Nil(err, "Valid from string")
	a.Equal("N0NE", addr.Call, "Call")
	a.Equal(0, addr.SSID, "SSID")

	addr = Address{}
	err = addr.FromString("N0CALLS")
	a.NotNil(err, "Invalid from string")
}

func TestPathFromString(t *testing.T) {
	p := Path{}
	p.FromString("WIDE1-1,WIDE2-1")
	assert.Equal(t, Path{{Call: "WIDE1", SSID: 1}, {Call: "WIDE2", SSID: 1}}, p, "Path")
}
