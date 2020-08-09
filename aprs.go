// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package aprs works with APRS string and byte packets.  It can upload
// those packets via APRS-IS or transmit them via TNC KISS.
package aprs

import "errors"

// Errors.
var (
	ErrCallNotVerified = errors.New("callsign not verified")
	ErrFrameBadControl = errors.New("frame Control Field not UI-frame")
	ErrFrameBadProto   = errors.New("frame Protocol ID not no layer 3 protocol")
	ErrFrameIncomplete = errors.New("frame incomplete")
	ErrFrameInvalid    = errors.New("frame is invalid")
	ErrFrameNoLast     = errors.New("frame incomplete or last path not set")
	ErrFrameShort      = errors.New("frame too short (16-bytes minimum)")
	ErrProtoScheme     = errors.New("protocol scheme is unknown")
)

// SwName is the default software name.
var SwName = "Go"

// SwVers is the default software version.
var SwVers = "3"
