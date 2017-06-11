// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Refer to Automatic Position Reporting System (APRS) Protocol
// Reference - Protocol version 1.0.

package aprs

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	uiFrame    = 0x03
	protocolID = 0xf0
)

// Errors.
var (
	ErrFrameBadControl  = errors.New("Frame error: Control Field not UI-frame")
	ErrFrameBadProtocol = errors.New("Frame error: Protocol ID not no layer 3 protocol")
	ErrFrameIncomplete  = errors.New("Frame error: incomplete")
	ErrFrameNoLast      = errors.New("Frame error: incomplete or last path not set")
	ErrFrameShort       = errors.New("Frame error: too short (16-bytes minimum)")
)

// Bytes converts an Address into its TNC byte representation.
func (a Address) Bytes() []byte {
	// AX.25 addresses are always 7-bytes:
	//	6-bytes/characters 7-bit ASCI encoded for the callsign.
	//	1-byte for the SSID and other data.
	bs := []byte{0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x96}

	// Convert callsign to 7-bit ASCII.
	for i, c := range a.Call {
		bs[i] = byte(c) << 1
	}

	// The last byte has the following bit breakdown:
	//
	// Bit: 7        | 6        | 5        | 4 3 2 1 | 0       |
	//	-------- +----------+----------+---------+---------+
	//      Repeated | Reserved | Reserved | SSID    | Last    |
	//	---------+----------+----------+---------+---------+
	//      0 = No   | 1        | 1        | 7-bit   | 0 = No  |
	//	1 = Yes  |          |          |         | 1 = Yes |
	bs[6] = (byte(a.SSID) << 1) | 0x60
	if a.Repeated {
		bs[6] |= 0x80 // Set bit 7
	}
	if a.last {
		bs[6] |= 0x01 // Set bit 0
	}

	return bs
}

// FromBytes converts a TNC byte address into an Address.
func (a *Address) FromBytes(addr []byte) (err error) {
	if len(addr) != 7 {
		err = fmt.Errorf("Address error: size mismatch %d != 7-bytes", len(addr))
		return
	}

	// Convert call from 7-bit encoding back to 8-bit
	for i := 0; i < 6; i++ {
		a.Call += string(addr[i] >> 1)
	}
	a.Call = strings.Replace(a.Call, " ", "", -1)

	// The last byte has the following bit breakdown:
	//
	// Bit: 7        | 6        | 5        | 4 3 2 1 | 0       |
	//      -------- +----------+----------+---------+---------+
	//      Repeated | Reserved | Reserved | SSID    | Last    |
	//      ---------+----------+----------+---------+---------+
	//      0 = No   | 1        | 1        | 7-bit   | 0 = No  |
	//      1 = Yes  |          |          |         | 1 = Yes |
	a.SSID = int(addr[6] & 0x1e >> 1)
	if addr[6]&0x01 > 0 {
		a.last = true
	}
	if addr[6]&0x80 > 0 {
		a.Repeated = true
	}

	return
}

// Bytes converts a Frame into its TNC byte representation appropriate
// for sending via KISS.
func (f Frame) Bytes() []byte {
	// Frame format is:
	//
	// Destination Address | Source Address | Path (0-8) | Control Field | Protocol ID | Information Field
	//             7-bytes |        7-bytes | 0-56 bytes |        1-byte |      1-byte | 1-256 bytes

	buf := bytes.NewBuffer([]byte{})

	// If we have a path then set the last address as such,
	// otherwise the source address is the last one.
	if len(f.Path) > 0 {
		f.Path[len(f.Path)-1].last = true
	} else {
		f.Src.last = true
	}

	buf.Write(f.Dst.Bytes()) // Destination Address
	buf.Write(f.Src.Bytes()) // Source Address
	// Path (optional)
	for _, a := range f.Path {
		buf.Write(a.Bytes())
	}
	buf.WriteByte(uiFrame)    // Control Field (always UI-frame)
	buf.WriteByte(protocolID) // Protocol ID (always no layer 3 protocol)
	buf.WriteString(f.Text)   // Information Field

	return buf.Bytes()
}

// FromBytes converts a TNC byte Frame into a Frame.
func (f *Frame) FromBytes(frame []byte) (err error) {
	if len(frame) < 16 {
		err = ErrFrameShort
		return
	}
	f.Dst.FromBytes(frame[0:7])  // Destination address
	f.Src.FromBytes(frame[7:14]) // Source address
	i := 14

	// Path (optional)
	if !f.Src.last {
		for {
			if i+7 > len(frame) {
				err = ErrFrameNoLast
				return
			}

			a := Address{}
			err = a.FromBytes(frame[i : i+7])
			if err != nil {
				return
			}
			f.Path = append(f.Path, a)
			i += 7

			if a.last {
				break
			}
		}
	}

	// To be valid the frame must have at least 2 more bytes for
	// the Control Field and Protocol ID.
	if i+2 > len(frame) {
		err = ErrFrameIncomplete
		return
	}

	// Control Field (always UI-frame)
	if frame[i] != uiFrame {
		err = ErrFrameBadControl
		return
	}

	// Protocol ID (always no layer 3 protocol)
	if frame[i+1] != protocolID {
		err = ErrFrameBadProtocol
		return
	}

	f.Text = string(frame[i+2:]) // Information Field

	return
}
