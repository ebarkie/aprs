package aprs

import (
	"testing"
	"time"
)

func TestRenderDataType(t *testing.T) {
	tests := []struct {
		name string
		pr   *PositionReport
		want byte
	}{
		{
			name: "no message, no time",
			pr:   &PositionReport{},
			want: byte(33),
		},
		{
			name: "yes message, no time",
			pr: &PositionReport{
				MessageCapable: true,
			},
			want: byte(61),
		},
		{
			name: "yes message, yes time",
			pr: &PositionReport{
				Timestamp:      time.Now(),
				MessageCapable: true,
			},
			want: byte(64),
		},
		{
			name: "no message, yes time",
			pr: &PositionReport{
				Timestamp: time.Now(),
			},
			want: byte(47),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.pr.renderDataType(); got != tc.want {
				t.Fatalf("Wanted: %v. Got %v", tc.want, got)
			}
		})
	}
}

func TestRenderTimestamp(t *testing.T) {
	tests := []struct {
		name string
		pr   *PositionReport
		want string
	}{{
		name: `2006-03-18T08:25:05Z`,
		pr:   prTimeHelper(t, `2006-03-18T08:25:05Z`),
		want: "180825z",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.pr.renderTimestamp(); string(got) != tc.want {
				t.Fatalf("Wanted: %s. Got %s", tc.want, string(got))
			}
		})
	}
}

func TestRenderCoords(t *testing.T) {
	tests := []struct {
		name string
		pr   *PositionReport
		want string
	}{
		{
			name: "Wyoming",
			pr: &PositionReport{
				Lat:    44.1083775,
				Lon:    -107.9386725,
				Symbol: []byte("/j"),
			},
			want: `4406.50N/10756.32Wj`,
		},
		{
			name: "Awamangu",
			pr: &PositionReport{
				Lat:    -46.071795,
				Lon:    169.6652273,
				Symbol: []byte(`\#`),
			},
			want: `4604.31S\16939.91E#`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.pr.renderCoords(); string(got) != tc.want {
				t.Fatalf("Wanted: %s. Got %s", tc.want, string(got))
			}
		})
	}
}

func TestRenderAltitude(t *testing.T) {
	tests := []struct {
		name string
		pr   *PositionReport
		want string
	}{{
		name: `2006-03-18T08:25:05Z`,
		pr: &PositionReport{
			Altitude: 30,
		},
		want: "/A=000030",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.pr.renderAltitude(); string(got) != tc.want {
				t.Fatalf("Wanted: %s. Got %s", tc.want, string(got))
			}
		})
	}
}

func TestExtensions(t *testing.T) {
	tests := []struct {
		name string
		pr   *PositionReport
		op   func(*PositionReport)
		want string
	}{
		{
			name: "CS-Extension-Without-DF",
			pr:   &PositionReport{},
			op: func(p *PositionReport) {
				p.CSExtension(360, 50, 0, 0)
			},
			want: "360/050",
		},
		{
			name: "CS-Extension-With-DF",
			pr:   &PositionReport{},
			op: func(p *PositionReport) {
				p.CSExtension(360, 50, 180, 9)
			},
			want: "360/050",
		},
		{
			name: "DS-Extension",
			pr:   &PositionReport{},
			op: func(p *PositionReport) {
				p.DSExtension(360, 50)
			},
			want: "360/050",
		},
		{
			name: "PHG-Extension",
			pr:   &PositionReport{},
			op: func(p *PositionReport) {
				p.PHGExtension(5, 1, 8, byte(59))
			},
			want: "PHG5;18",
		},
		{
			name: "RNG-Extension",
			pr:   &PositionReport{},
			op: func(p *PositionReport) {
				p.RNGExtension(25)
			},
			want: "RNG0025",
		},
		{
			name: "DFS-Extension",
			pr:   &PositionReport{},
			op: func(p *PositionReport) {
				p.DFSExtension(5, 1, 8, byte(59))
			},
			want: "DFS5;18",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.op(tc.pr)
			if got := string(tc.pr.Extn); got != tc.want {
				t.Fatalf("Wanted: %s. Got %s", tc.want, string(got))
			}
		})
	}
}

// prTimeHelper helps deal with the err returned by time.Parse()
func prTimeHelper(t *testing.T, tString string) *PositionReport {
	t.Helper()
	tParsed, err := time.Parse(time.RFC3339, tString)
	if err != nil {
		t.Fatalf("Error parsing time %s: %v", tString, err)
	}
	return &PositionReport{Timestamp: tParsed}
}
