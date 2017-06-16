// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Errors.
var (
	ErrNotVerified     = errors.New("Not verified but scheme requires it")
	ErrUnhandledScheme = errors.New("Unhandled scheme")
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

// SendIS sends a Frame to the specified APRS-IS dial string.  The
// dial string should be in the form scheme://host:port with
// scheme being http, tcp, or udp.  This is most commonly used for
// CWOP.
func (f Frame) SendIS(dial string, pass int) error {
	// Refer to Connecting to APRS-IS:
	// http://www.aprs-is.net/Connecting.aspx

	parts := strings.Split(strings.ToLower(dial), "://")
	if len(parts) != 2 {
		return net.InvalidAddrError(dial)
	}

	switch parts[0] {
	case "http":
		return f.SendHTTP(dial, pass)
	case "tcp":
		return f.SendTCP(parts[1], pass)
	case "udp":
		return f.SendUDP(parts[1], pass)
	}

	return ErrUnhandledScheme
}

// SendHTTP sends a Frame to the specified APRS-IS host over the
// HTTP protocol.  This scheme is the least efficient and requires
// a verified connection (real callsign and passcode) but is
// reliable and provides acknowledgement of receipt.
func (f Frame) SendHTTP(dial string, pass int) (err error) {
	if pass < 0 {
		err = ErrNotVerified
		return
	}

	data := fmt.Sprintf("user %s pass %d vers %s %s\r\n%s", f.Src, pass, SwName, SwVers, f)

	var req *http.Request
	req, err = http.NewRequest("POST", dial, bytes.NewBufferString(data))
	if err != nil {
		return
	}
	req.Header.Set("Accept-Type", "text/plain")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP request returned non-OK status code %d", resp.StatusCode)
	}

	return
}

// SendUDP sends a Frame to the specified APRS-IS host over the
// UDP protocol.  This scheme is the most efficient but requires
// a verified connection (real callsign and passcode) and has no
// acknowledgement of receipt.
func (f Frame) SendUDP(dial string, pass int) (err error) {
	if pass < 0 {
		err = ErrNotVerified
		return
	}

	var conn net.Conn
	conn, err = net.Dial("udp", dial)
	if err != nil {
		return
	}
	defer conn.Close()

	// Send data packet
	fmt.Fprintf(conn, "user %s pass %d vers %s %s\r\n%s", f.Src, pass, SwName, SwVers, f)

	return
}

// SendTCP sends a Frame to the specified APRS-IS host over the
// TCP protocol.  This scheme is the oldest, most compatible, and
// allows unverified connections.
func (f Frame) SendTCP(dial string, pass int) (err error) {
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
