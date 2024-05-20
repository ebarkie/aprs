package main

import (
	"log"

	"github.com/acobaugh/aprs"
)

func main() {
	w := aprs.Wx{
		Lat:  35.7,
		Lon:  -78.7,
		Type: "DvsVP2+",
	}
	w.Pressure = 29.70
	w.Humidity = 90
	w.RainLastHour = 0.0
	w.RainLast24Hours = 0.10
	w.Temp = 85
	w.WindDir = 180
	w.WindSpeed = 5

	f := aprs.Frame{
		Dst:  aprs.Addr{Call: "APRS"},
		Src:  aprs.Addr{Call: "aWnnnn"},
		Path: aprs.Path{aprs.Addr{Call: "TCPIP", Repeated: true}},
		Text: w.String(),
	}
	err := f.SendIS("tcp://cwop.aprs.net:14580", -1)
	if err != nil {
		log.Printf("Upload error: %s", err)
	}

	f = aprs.Frame{}
	f.Dst.FromString("APZ001") // Experimental v0.0.1
	f.Src.FromString("N0CALL-13")
	f.Path.FromString("WIDE1-1,WIDE2-1")
	f.Text = w.String()
	err = f.SendKISS("direwolf:8001")
	if err != nil {
		log.Printf("Network TNC transmit error: %s", err)
	}
}
