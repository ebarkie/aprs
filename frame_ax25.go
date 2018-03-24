// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Refer to Automatic Position Reporting System (APRS) Protocol
// Reference - Protocol version 1.0.

package aprs

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	uiFrame    = 0x03
	protocolID = 0xf0
)

// Bytes returns the address in AX.25 byte format.
func (a Addr) Bytes() []byte {
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
	// Bit  7        | 6        | 5        | 4 3 2 1 | 0
	//     ----------+----------+----------+---------+---------
	//      Repeated | Reserved | Reserved | SSID    | Last
	//     ----------+----------+----------+---------+---------
	//      0 = No   | 1        | 1        | 7-bit   | 0 = No
	//      1 = Yes  |          |          |         | 1 = Yes
	bs[6] = (byte(a.SSID) << 1) | 0x60
	if a.Repeated {
		bs[6] |= 0x80 // Set bit 7
	}
	if a.last {
		bs[6] |= 0x01 // Set bit 0
	}

	return bs
}

// FromBytes sets the address from an AX.25 byte slice.
func (a *Addr) FromBytes(addr []byte) error {
	if len(addr) != 7 {
		return fmt.Errorf("Address error: size mismatch %d != 7-bytes", len(addr))
	}

	// Convert call from 7-bit encoding back to 8-bit
	for i := 0; i < 6; i++ {
		a.Call += string(addr[i] >> 1)
	}
	a.Call = strings.Replace(a.Call, " ", "", -1)

	// The last byte has the following bit breakdown:
	//
	// Bit  7        | 6        | 5        | 4 3 2 1 | 0
	//     ----------+----------+----------+---------+---------
	//      Repeated | Reserved | Reserved | SSID    | Last
	//     ----------+----------+----------+---------+---------
	//      0 = No   | 1        | 1        | 7-bit   | 0 = No
	//      1 = Yes  |          |          |         | 1 = Yes
	a.SSID = int(addr[6] & 0x1e >> 1)
	if addr[6]&0x01 > 0 {
		a.last = true
	}
	if addr[6]&0x80 > 0 {
		a.Repeated = true
	}

	return nil
}

// Bytes returns the Frame in AX.25 byte format.  This is suitable for
// sending to a TNC.
func (f Frame) Bytes() []byte {
	// Frame format is:
	//
	// Destination address | Source address | Path (0-8) | Control field | Protocol ID | Information field
	//             7-bytes |        7-bytes | 0-56 bytes |        1-byte |      1-byte | 1-256 bytes

	buf := bytes.NewBuffer([]byte{})

	// If we have a path then set the last address as such,
	// otherwise the source address is the last one.
	if len(f.Path) > 0 {
		f.Path[len(f.Path)-1].last = true
	} else {
		f.Src.last = true
	}

	buf.Write(f.Dst.Bytes()) // Destination address
	buf.Write(f.Src.Bytes()) // Source address
	// Path (optional)
	for _, a := range f.Path {
		buf.Write(a.Bytes())
	}
	buf.WriteByte(uiFrame)    // Control field (always UI-frame)
	buf.WriteByte(protocolID) // Protocol ID (always no layer 3 protocol)
	buf.WriteString(f.Text)   // Information field

	return buf.Bytes()
}

// FromBytes sets the Frame from an AX.25 byte slice.
func (f *Frame) FromBytes(frame []byte) error {
	if len(frame) < 16 {
		return ErrFrameShort
	}
	f.Dst.FromBytes(frame[0:7])  // Destination address
	f.Src.FromBytes(frame[7:14]) // Source address
	i := 14

	// Path (optional)
	if !f.Src.last {
		for {
			if i+7 > len(frame) {
				return ErrFrameNoLast
			}

			a := Addr{}
			err := a.FromBytes(frame[i : i+7])
			if err != nil {
				return err
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
		return ErrFrameIncomplete
	}

	// Control Field (always UI-frame)
	if frame[i] != uiFrame {
		return ErrFrameBadControl
	}

	// Protocol ID (always no layer 3 protocol)
	if frame[i+1] != protocolID {
		return ErrFrameBadProto
	}

	f.Text = string(frame[i+2:]) // Information field

	return nil
}
