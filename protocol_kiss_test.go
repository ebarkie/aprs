// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKISSEscape(t *testing.T) {
	for _, test := range []struct {
		from []byte
		to   []byte
	}{
		{[]byte("test"), []byte{0x74, 0x65, 0x73, 0x74}},
		{[]byte{0x74, 0x65, 0x73, fesc, 0x74}, []byte{0x74, 0x65, 0x73, fesc, tfesc, 0x74}},
		{[]byte{fend}, []byte{fesc, tfend}},
		{[]byte{fesc, fesc}, []byte{fesc, tfesc, fesc, tfesc}},
	} {
		assert.Equal(t, test.to, kissEscape(test.from), "Escaped string")
	}
}
