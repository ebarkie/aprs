package main

import (
	"context"
	"log"
	"time"

	"github.com/acobaugh/aprs"
)

func main() {
	log.Println("Receiving")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: "N0CALL"}, -1, "t/w")
	for f := range fc {
		log.Println(f)
	}
	log.Println("Close")
}
