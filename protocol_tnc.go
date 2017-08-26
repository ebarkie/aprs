// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Refer to The KISS TNC: A simple Host-to-TNC communications
// protocol paper or Wikipedia's KISS (TNC) information.

package aprs

import (
	"bytes"
	"net"
)

const (
	fend  = 0xc0 // Frame end
	tfend = 0xdc // Transformed frame end
	fesc  = 0xdb // Frame escape
	tfesc = 0xdd // Transformed frame escape
)

func kissEscape(b []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	for i := range b {
		switch b[i] {
		case fend:
			buf.Write([]byte{fesc, tfend})
		case fesc:
			buf.Write([]byte{fesc, tfesc})
		default:
			buf.WriteByte(b[i])
		}
	}

	return buf.Bytes()
}

// SendTNC sends a Frame to the specified network TNC device
// using the KISS protocol for transmission over RF.
func (f Frame) SendTNC(dial string) (err error) {
	const (
		cmdData   = 0x00 // Frame contains data that should be sent out of the TNC
		cmdReturn = 0xff // Exit KISS mode
	)

	const port = 0 // XXX this can be made a variable if necessary

	conn, err := net.Dial("tcp", dial)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.Write([]byte{fend, cmdData | ((port & 0xf) << 4)})
	conn.Write(kissEscape(f.Bytes()))
	conn.Write([]byte{fend})

	return
}
