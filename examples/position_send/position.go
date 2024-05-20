package main

import (
	"log"

	"github.com/acobaugh/aprs"
)

// N0CALL-13>APZ001,WIDE1-1,WIDE2-1:!3542.00N\07842.00Wj146.520MHz Flat tire; Send beer!
func main() {
	p := aprs.PositionReport{ // create a position report
		Lat:    35.7,
		Lon:    -78.7,
		Symbol: `\j`, // symbol for a jeep
		Freq: &aprs.Freq{ // add a frequency report (optional)
			Mhz: 146.52, // I am listening on US 2m simplex
		},
		Comment: "Flat tire; Send beer!", // self explainatory
	}

	f := aprs.Frame{}                                    // create an ax25 frame
	f.Dst.FromString("APZ001")                           // Experimental v0.0.1
	f.Src.FromString("N0CALL-13")                        // your callsign here
	f.Path.FromString("WIDE1-1,WIDE2-1")                 // http://blog.aprs.fi/2020/02/how-aprs-paths-work.html
	f.Text = p.String()                                  // pack the position report into the frame
	if err := f.SendKISS("localhost:8001"); err != nil { // send it
		log.Printf("Network TNC transmit error: %s", err) // or not
	}
}
