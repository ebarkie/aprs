package aprs

import (
	"testing"
)

func TestRenderFreq(t *testing.T) {
	tests := []struct {
		name string
		pr   *PositionReport
		want string
	}{
		{
			name: `simple_146.72`,
			pr: &PositionReport{
				Freq: &Freq{
					Mhz: 146.72,
				},
			},
			want: "146.720MHz ",
		},
		{
			name: `146.72-range-25m`,
			pr: &PositionReport{
				Freq: &Freq{
					Mhz:   146.72,
					Range: 25,
				},
			},
			want: "146.720MHz R025m ",
		},
		{
			name: `440.050-repeater`,
			pr: &PositionReport{
				Freq: &Freq{
					Mhz:    440.050,
					CTCSS:  100,
					Offset: -500,
					Range:  25,
				},
			},
			want: "440.050MHz C100 -500 R025m ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.pr.renderFreq(); string(got) != tc.want {
				t.Fatalf("Wanted: `%s`. Got: `%s`", tc.want, got)
			}
		})
	}
}
