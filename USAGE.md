# aprs

```go
import "github.com/ebarkie/aprs"
```

Package aprs works with APRS string and byte packets. It can upload those
packets via APRS-IS or transmit them via TNC KISS.

## Usage

```go
var (
	ErrFrameBadControl  = errors.New("Frame error: Control Field not UI-frame")
	ErrFrameBadProtocol = errors.New("Frame error: Protocol ID not no layer 3 protocol")
	ErrFrameIncomplete  = errors.New("Frame error: incomplete")
	ErrFrameNoLast      = errors.New("Frame error: incomplete or last path not set")
	ErrFrameShort       = errors.New("Frame error: too short (16-bytes minimum)")
)
```
Errors.

```go
var (
	ErrNotVerified     = errors.New("Not verified but scheme requires it")
	ErrUnhandledScheme = errors.New("Unhandled scheme")
)
```
Errors.

```go
var SwName = "Go"
```
SwName is the default software name.

```go
var SwVers = "2"
```
SwVers is the default software version.

#### func  GenPass

```go
func GenPass(call string) (pass uint16)
```
GenPass generates a verification passcode for the given station.

#### type Address

```go
type Address struct {
	SSID     int
	Repeated bool

	Call string
}
```

Address represents an APRS callsign, SSID, and associated metadata.

#### func (Address) Bytes

```go
func (a Address) Bytes() []byte
```
Bytes converts an Address into its TNC byte representation.

#### func (*Address) FromBytes

```go
func (a *Address) FromBytes(addr []byte) error
```
FromBytes converts a TNC byte address into an Address.

#### func (*Address) FromString

```go
func (a *Address) FromString(addr string) (err error)
```
FromString converts a text address into an Address.

#### func (Address) String

```go
func (a Address) String() (addr string)
```
String converts an Address into its text representation.

#### type Frame

```go
type Frame struct {
	Dst  Address
	Src  Address
	Path Path
	Text string
}
```

Frame represents a complete APRS frame.

#### func (Frame) Bytes

```go
func (f Frame) Bytes() []byte
```
Bytes converts a Frame into its TNC byte representation appropriate for sending
via KISS.

#### func (*Frame) FromBytes

```go
func (f *Frame) FromBytes(frame []byte) error
```
FromBytes converts a TNC byte Frame into a Frame.

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

#### func (Frame) SendTCP

```go
func (f Frame) SendTCP(dial string, pass int) (err error)
```
SendTCP sends a Frame to the specified APRS-IS host over the TCP protocol. This
scheme is the oldest, most compatible, and allows unverified connections.

#### func (Frame) SendTNC

```go
func (f Frame) SendTNC(dial string) (err error)
```
SendTNC sends a Frame to the specified network TNC device using the KISS
protocol for transmission over RF.

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
String converts a Frame into its text representation appropriate for printing or
sending via APRS-IS.

#### type Path

```go
type Path []Address
```

Path represents the APRS digipath.

#### func (*Path) FromString

```go
func (p *Path) FromString(path string) (err error)
```
FromString converts a list of comma separated addreses into a Path.

#### type Wx

```go
type Wx struct {
	Lat  float64
	Lon  float64
	Type string
}
```

Wx holds the weather station information.

#### func (*Wx) Altimeter

```go
func (w *Wx) Altimeter(f float64)
```
Altimeter is the barometric pressure as altimeter (station plus elevation
correction) in inches.

#### func (*Wx) Clear

```go
func (w *Wx) Clear()
```
Clear clears all measurements in the observation payload.

#### func (*Wx) Humidity

```go
func (w *Wx) Humidity(i int)
```
Humidity is the relative humidity.

#### func (*Wx) RainLast24Hours

```go
func (w *Wx) RainLast24Hours(f float64)
```
RainLast24Hours is the rain accumulation for the last 24 hours in inches.

#### func (*Wx) RainRate

```go
func (w *Wx) RainRate(f float64)
```
RainRate is the rain accumulation for the last hour in inches.

#### func (*Wx) RainToday

```go
func (w *Wx) RainToday(f float64)
```
RainToday is the rain accumulation so far today in inches.

#### func (*Wx) SolarRadiation

```go
func (w *Wx) SolarRadiation(i int)
```
SolarRadiation is the solar radiation in W/m^2.

#### func (Wx) String

```go
func (w Wx) String() (t string)
```
String returns an APRS packet for the provided measurements.

#### func (*Wx) Temperature

```go
func (w *Wx) Temperature(f float64)
```
Temperature is the temperature in F.

#### func (*Wx) Timestamp

```go
func (w *Wx) Timestamp(t time.Time)
```
Timestamp is the time of the observation.

#### func (*Wx) WindDirection

```go
func (w *Wx) WindDirection(i int)
```
WindDirection is the wind direction in degrees.

#### func (*Wx) WindGust

```go
func (w *Wx) WindGust(i int)
```
WindGust is the peak wind speed for the previous 5 minutes in mph.

#### func (*Wx) WindSpeed

```go
func (w *Wx) WindSpeed(i int)
```
WindSpeed is the wind speed in mph.
