// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenPass(t *testing.T) {
	a := assert.New(t)

	p := GenPass("N0CALL")
	a.Equal(13023, int(p), "Passcode mismatch")
}
