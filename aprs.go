// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package aprs works with APRS string and byte packets.  It can upload
// those packets via APRS-IS or transmit them via TNC KISS.
package aprs

import "errors"

// Errors.
var (
	ErrFrameBadControl  = errors.New("Frame error: Control Field not UI-frame")
	ErrFrameBadProtocol = errors.New("Frame error: Protocol ID not no layer 3 protocol")
	ErrFrameIncomplete  = errors.New("Frame error: incomplete")
	ErrFrameNoLast      = errors.New("Frame error: incomplete or last path not set")
	ErrFrameShort       = errors.New("Frame error: too short (16-bytes minimum)")
	ErrNotVerified      = errors.New("Not verified but scheme requires it")
	ErrUnhandledScheme  = errors.New("Unhandled scheme")
)

// SwName is the default software name.
var SwName = "Go"

// SwVers is the default software version.
var SwVers = "2"
