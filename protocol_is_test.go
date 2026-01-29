// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenPass(t *testing.T) {
	p := GenPass("N0CALL")
	assert.Equal(t, 13023, int(p), "Passcode mismatch")
}

func TestGenPassOdd(t *testing.T) {
	p := GenPass("N0CAL")
	assert.Equal(t, 12947, int(p), "Passcode mismatch")
}
