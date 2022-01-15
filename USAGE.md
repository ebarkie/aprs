# aprs

```go
import "github.com/ebarkie/aprs"
```

Package aprs works with APRS string and byte packets. It can upload those
packets via APRS-IS or transmit them via TNC KISS.

## Usage

```go
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
```
Errors.

```go
var SwName = "Go"
```
SwName is the default software name.

```go
var SwVers = "3"
```
SwVers is the default software version.

#### func  GenPass

```go
func GenPass(call string) (pass uint16)
```
GenPass generates a verification passcode for the given station.

#### func  RecvIS

```go
func RecvIS(ctx context.Context, dial string, user Addr, pass int, filters ...string) <-chan Frame
```
RecvIS receives APRS-IS frames over tcp from the specified server. Filter(s) are
optional and use the following syntax:

http://www.aprs-is.net/javAPRSFilter.aspx

#### type Addr

```go
type Addr struct {
	SSID     int
	Repeated bool

	Call string
}
```

Addr represents an APRS callsign, SSID, and associated metadata.

#### func (Addr) Bytes

```go
func (a Addr) Bytes() []byte
```
Bytes returns the address in AX.25 byte format.

#### func (*Addr) FromBytes

```go
func (a *Addr) FromBytes(addr []byte) error
```
FromBytes sets the address from an AX.25 byte slice.

#### func (*Addr) FromString

```go
func (a *Addr) FromString(addr string) (err error)
```
FromString sets the address from a string.

#### func (Addr) String

```go
func (a Addr) String() (addr string)
```
String returns the address as a TNC2 formatted string.

#### type Frame

```go
type Frame struct {
	Dst  Addr
	Src  Addr
	Path Path
	Text string
}
```

Frame represents a complete APRS frame.

#### func (Frame) Bytes

```go
func (f Frame) Bytes() []byte
```
Bytes returns the Frame in AX.25 byte format. This is suitable for sending to a
TNC.

#### func (*Frame) FromBytes

```go
func (f *Frame) FromBytes(frame []byte) error
```
FromBytes sets the Frame from an AX.25 byte slice.

#### func (*Frame) FromString

```go
func (f *Frame) FromString(frame string) (err error)
```
FromString sets the Frame from a TNC2 formatted string.

This strictly enforces the AX.25 specification and will return errors if
callsigns are greater than 6 characters or SSID's are not numeric values between
0 and 15.

#### func (Frame) SendHTTP

```go
func (f Frame) SendHTTP(dial string, pass int) (err error)
```
SendHTTP sends a Frame to the specified APRS-IS host over the HTTP protocol.
This scheme is the least efficient and requires a verified connection (real
callsign and passcode) but is reliable and provides acknowledgement of receipt.

#### func (Frame) SendIS

```go
func (f Frame) SendIS(dial string, pass int) error
```
SendIS sends a Frame to the specified APRS-IS dial string. The dial string
should be in the form scheme://host:port with scheme being http, tcp, or udp.
This is most commonly used for CWOP.

#### func (Frame) SendKISS

```go
func (f Frame) SendKISS(dial string) (err error)
```
SendKISS sends a Frame to the specified network TNC device using the KISS
protocol for transmission over RF.

#### func (Frame) SendTCP

```go
func (f Frame) SendTCP(dial string, pass int) (err error)
```
SendTCP sends a Frame to the specified APRS-IS host over the TCP protocol. This
scheme is the oldest, most compatible, and allows unverified connections.

#### func (Frame) SendUDP

```go
func (f Frame) SendUDP(dial string, pass int) (err error)
```
SendUDP sends a Frame to the specified APRS-IS host over the UDP protocol. This
scheme is the most efficient but requires a verified connection (real callsign
and passcode) and has no acknowledgement of receipt.

#### func (Frame) String

```go
func (f Frame) String() (frame string)
```
String returns the Frame as a TNC2 formatted string. This is suitable for
sending to APRS-IS servers.

#### type Freq

```go
type Freq struct {
	Mhz    float64 // Frequency in Mhz
	Tone   int     // Tone in hz
	CTCSS  int     // ctcss (mutually exclusive with tone and dcs)
	DCS    int     // dcs (mutually exclusive with tone and ctcss)
	Offset int     // +/- offset in mhz
	Range  int     // range in miles
	Narrow bool    // defaults to false for wideband
}
```

Freq contains an APRS FreqSpec frequency report
http://www.aprs.org/info/freqspec.txt

#### func (*Freq) Render

```go
func (f *Freq) Render() string
```
Render renders the frequency struct into a byte-slice

#### type Path

```go
type Path []Addr
```

Path represents the APRS digipath.

#### func (*Path) FromString

```go
func (p *Path) FromString(path string) (err error)
```
FromString sets the Path from a string of comma separated addresses.

#### type PositionReport

```go
type PositionReport struct {
	Timestamp      time.Time // NOTE: Position reports are NOT expected to contain a timestamp unless the report refers to "old" (not real-time) data.
	Lat            float64   // latitude
	Lon            float64   // longitude
	Altitude       int
	Symbol         string // 2 byte Map symbol; see Chapter 20 aprs101
	Extn           string // 7+ byte Data Extension field. See Chapter 7 pg27 aprs101
	Freq           *Freq  // freqspec compatible Frequency report
	Comment        string // free-form comment
	MessageCapable bool   // Stations without APRS messaging capability are typically stand-alone trackers or digipeaters.
}
```

PositionReport wraps all necessary metadata for a position report

#### func (*PositionReport) CSExtension

```go
func (p *PositionReport) CSExtension(course, speed, bearing, nrq int)
```
CSExtension sets a course/speed data-extension block in the report

#### func (*PositionReport) DFSExtension

```go
func (p *PositionReport) DFSExtension(str, gain, dir int, height byte)
```
DFSExtension sets a Direction-finding data-extension block

#### func (*PositionReport) DSExtension

```go
func (p *PositionReport) DSExtension(direction, speed int)
```
DSExtension Extension sets a wind direction/speed data-extension block

#### func (*PositionReport) PHGExtension

```go
func (p *PositionReport) PHGExtension(power, gain, dir int, height byte)
```
PHGExtension sets a power/height/gain data-extension block

#### func (*PositionReport) RNGExtension

```go
func (p *PositionReport) RNGExtension(miles int)
```
RNGExtension sets a pre-computed range data-extension block

#### func (*PositionReport) String

```go
func (p *PositionReport) String() string
```
String returns a rendered position report suitable for sending to a TNC

#### type Wx

```go
type Wx struct {
	Lat  float64
	Lon  float64
	Type string

	Timestamp time.Time

	Altimeter       float64
	Humidity        int
	RainLastHour    float64
	RainLast24Hours float64
	RainToday       float64
	SolarRad        int
	Temp            int
	WindDir         int
	WindGust        int
	WindSpeed       int
}
```

Wx represents a weather station observation.

#### func (Wx) String

```go
func (w Wx) String() (s string)
```
String returns an APRS packet for the provided measurements.

#### func (*Wx) Zero

```go
func (w *Wx) Zero()
```
Zero zeroes all measurements in the observation payload.
