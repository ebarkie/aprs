// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// GenPass generates a verification passcode for the given station.
func GenPass(call string) (pass uint16) {
	// Refer to aprsc:
	// https://github.com/hessu/aprsc

	// Upper case callsign and strip SSID if it was included
	c := strings.ToUpper(call)
	dash := strings.Index(c, "-")
	if dash > -1 {
		c = c[:dash]
	}

	pass = 0x73e2 // The key/seed.
	for i := 0; i < len(c); i += 2 {
		pass ^= uint16(c[i]) << 8
		pass ^= uint16(c[i+1])
	}

	// Mask off the high bit so number is always positive
	pass &= 0x7fff
	return
}

func readLine(conn net.Conn) (string, error) {
	const socketTimeout = 5 * time.Second

	conn.SetReadDeadline(time.Now().Add(socketTimeout))
	s, err := bufio.NewReader(conn).ReadString('\n')
	return strings.TrimSpace(s), err
}

// SendIS sends a Frame to the specified APRS-IS host.  It is
// most commonly used for CWOP.
func (f Frame) SendIS(dial string, pass int) (err error) {
	// Refer to Connecting to APRS-IS:
	// http://www.aprs-is.net/Connecting.aspx

	var conn net.Conn
	conn, err = net.Dial("tcp", dial)
	if err != nil {
		return
	}
	defer conn.Close()

	// Read welcome banner
	_, err = readLine(conn)
	if err != nil {
		return
	}

	// Login
	fmt.Fprintf(conn, "user %s pass %d vers %s %s\r\n", f.Src, pass, SwName, SwVers)

	// # logresp CWxxxx unverified, server CWOP-7
	// # logresp CWxxxx unverified, server THIRD
	_, err = readLine(conn)
	if err != nil {
		return
	}

	// Send frame
	fmt.Fprintf(conn, "%s\r\n", f)

	return
}
