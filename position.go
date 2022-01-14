package aprs

import (
	"fmt"
	"time"
)

// PositionReport wraps all necessary metadata for a position report
type PositionReport struct {
	Timestamp      time.Time // NOTE: Position reports are NOT expected to contain a timestamp unless the report refers to "old" (not real-time) data.
	Lat            float64   // latitude
	Lon            float64   // longitude
	Altitude       int
	Symbol         []string // 2 byte Map symbol; see Chapter 20 aprs101
	Extn           string   // 7+ byte Data Extension field. See Chapter 7 pg27 aprs101
	Freq           *Freq    // freqspec compatible Frequency report
	Comment        string   // free-form comment
	MessageCapable bool     // Stations without APRS messaging capability are typically stand-alone trackers or digipeaters.
}

// String returns a rendered position report suitable for sending to a TNC
func (p *PositionReport) String() string {
	// Refer to APRS protocol reference 1.0
	// Chapter 8: position and df report data formats
	var out string

	// render the data type
	out = string(p.renderDataType())

	// render the timestamp
	if !p.Timestamp.IsZero() {
		out = out + p.renderTimestamp()
	}

	// render the lat/long coords
	out = out + p.renderCoords()

	// render the data extension block (must be at least 7 bytes)
	if len(p.Extn) >= 7 {
		out = out + p.Extn
	}

	// render the Freq if it exists
	if p.Freq != nil {
		out = out + p.renderFreq()
	}

	// render altitude if it exists
	if p.Altitude > 0 {
		out = out + p.renderAltitude()
	}

	// add any other comments
	out = out + p.Comment

	// truncate to the size of the ui frame and return
	return out[:255]
}

// renderDataType returns the report Data-type (based on timestamp and messaging setting)
func (p *PositionReport) renderDataType() byte {
	if p.MessageCapable {
		if p.Timestamp.IsZero() {
			return byte(61) //("=") messaging, no timestamp
		}
		return byte(64) //("@") messaging timestamp
	}
	if p.Timestamp.IsZero() {
		return byte(33) //("!") no messaging, no timestamp
	}
	return byte(47) //("/") no messaging, timestamp
}

// renderTimestamp returns the rendered timestamp from the position report
func (p *PositionReport) renderTimestamp() string {
	return fmt.Sprintf("%sz", p.Timestamp.Format("021504"))
}

// renderCoords returns the rendered latitude and longitude from the position report
func (p *PositionReport) renderCoords() string {
	latDeg, latMin, latHem := decToDMS(p.Lat, [2]string{"N", "S"})
	lonDeg, lonMin, lonHem := decToDMS(p.Lon, [2]string{"E", "W"})

	return fmt.Sprintf("%02.0f%05.2f%s%s%03.0f%05.2f%s%s",
		latDeg, latMin, latHem,
		p.Symbol[0],
		lonDeg, lonMin, lonHem,
		p.Symbol[1])
}

// renderFreq returns the rendered freqspec compatible Frequency
func (p *PositionReport) renderFreq() string {
	var out string
	// add a delimiter if a data-extension exists
	if len(p.Extn) >= 7 {
		out = out + `/`
	}
	return (out + p.Freq.Render())
}

// renderAltitude returns the rendered freqspec compatible Frequency
func (p *PositionReport) renderAltitude() string {
	return fmt.Sprintf("/A=%06d", p.Altitude)
}

// Data Extensions:
// A fixed-length 7-byte field may follow APRS position data. This field is an
// APRS Data Extension. The extension may be one of the following:
// • CSE/SPD Course and Speed (this may be followed by a further 8 bytes
// containing DF bearing and Number/Range/Quality parameters)
// • DIR/SPD Wind Direction and Wind Speed
// • PHGphgd Station Power and Effective Antenna Height/Gain/
// Directivity
// • RNGrrrr Pre-Calculated Radio Range
// • DFSshgd DF Signal Strength and Effective Antenna Height/Gain
// • Tyy/Cxx Area Object Descriptor

// CSExtension sets a course/speed data-extension block in the report
func (p *PositionReport) CSExtension(course, speed, bearing, nrq int) {
	if bearing+nrq > 0 {
		// if Direction-finding data is included, hard-code the map-symbol to DF
		// aprs 101, pg 34; df reports
		p.Symbol = []string{`/`, `\`}
		p.Extn = fmt.Sprintf("%s/%s/%s/%s",
			z3p(course),
			z3p(speed),
			z3p(bearing),
			z3p(nrq))
	}
	p.Extn = fmt.Sprintf("%s/%s",
		z3p(course),
		z3p(speed))
}

// DSExtension Extension sets a wind direction/speed data-extension block
func (p *PositionReport) DSExtension(direction, speed int) {
	p.Extn = fmt.Sprintf("%s/%s",
		z3p(direction),
		z3p(speed))
}

// PHGExtension sets a power/height/gain data-extension block
func (p *PositionReport) PHGExtension(power, gain, dir int, height byte) {
	//-------------+------+------+-----+-------+------+-------+------+-------+------+------+-------+
	//|  phgdCode:  |  0   |  1   |  2  |   3   |  4   |   5   |  6   |   7   |  8   |  9   | Units |
	//+-------------+------+------+-----+-------+------+-------+------+-------+------+------+-------+
	//| Power       | 0    | 1    | 4   | 9     | 16   | 25    | 36   | 49    | 64   |   81 | watts |
	//| Height      | 10   | 20   | 40  | 80    | 160  | 320   | 640  | 1280  | 2560 | 5120 | feet  |
	//| Gain        | 0    | 1    | 2   | 3     | 4    | 5     | 6    | 7     | 8    |    9 | dB    |
	//| Directivity | omni | 45NE | 90E | 135SE | 180S | 225SW | 270W | 315NW | 360N |      | de    |
	//+-------------+------+------+-----+-------+------+-------+------+-------+------+------+-------+
	p.Extn = fmt.Sprintf("PHG%d%s%d%d",
		max(9, power),
		string(height), // `The height code may in fact be any ASCII character 0–9 and above`
		max(9, gain),
		max(8, dir))
}

// RNGExtension sets a pre-computed range data-extension block
func (p *PositionReport) RNGExtension(miles int) {
	p.Extn = fmt.Sprintf("RNG%04d", miles)
}

// DFSExtension sets a Direction-finding data-extension block
func (p *PositionReport) DFSExtension(str, gain, dir int, height byte) {
	//+-------------+------+------+-----+-------+------+-------+------+-------+------+------+----------+
	//| shgd-Code:  |  0   |  1   |  2  |   3   |  4   |   5   |  6   |   7   |  8   |  9   |  Units   |
	//+-------------+------+------+-----+-------+------+-------+------+-------+------+------+----------+
	//| Strength    | 0    | 1    | 2   | 3     | 4    | 5     | 6    | 7     | 8    | 9    | S-points |
	//| Height      | 10   | 20   | 40  | 80    | 160  | 320   | 640  | 1280  | 2560 | 5120 | feet     |
	//| Gain        | 0    | 1    | 2   | 3     | 4    | 5     | 6    | 7     | 8    | 9    | dB       |
	//| Directivity | omni | 45NE | 90E | 135SE | 180S | 225SW | 270W | 315NW | 360N |      | deg      |
	//+-------------+------+------+-----+-------+------+-------+------+-------+------+------+----------+
	p.Symbol = []string{`/`, `\`}
	p.Extn = fmt.Sprintf("DFS%d%s%d%d",
		max(9, str),
		string(height),
		max(9, gain),
		max(8, dir))
}
