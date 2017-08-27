// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tnc2Bad = []string{
	"N0CALLL>APZ001:Hello world",
	"N0CALL-16>APZ001:Hello world",
	"N0CALL-XX>APZ001:Hello world",
	"N0CALL>APZ001,N0CALL-XX:Hello world",
}

var tnc2Good = []string{
	"N0CALL>APZ001:Hello world",
	"N0CALL-13>APZ001:Hello world",
	"N0CALL-13>APZ001,WIDE1-1,WIDE2-1:Hello world",
	"N0CALL-13>APZ001,TCPIP*,qAS,N0CALL:Hello world",
	"W4DEX-2>APU25N,TCPIP*,qAC,FIFTH:@270130z3515.35N/08022.82W_000/000g000t073r000p000P000h74b10177/WX",
	"K4CCC-9>APRS,NE4SC-12,WA4USN-3*,qAR,WA4USN-5:!!0000000402ED024E27CD0383--------00EE050700000000",
	"WR4AGC-5>APRX23,TCPIP*,qAC,SIXTH::WR4AGC-5 :PARM.Avg 10m,Avg 10m,RxPkts,IGateDropRx,TxPkts",
	"K4OGB-9>APN383,WIDE2-2,qAR,KM4FZA:!3521.61NS08017.87W#PHG7430/W3,NC3 Digi  Albemarle NC",
	"K4CCC-9>APRS,NE4SC-12,W4HRS-15,WIDE2*,qAR,WA4USN-5:CHESTERFIELD COUNTY AMATEUR RADIO SOCIETY",
	"K4JH-1>APRX28,TCPIP*,qAC,T2MCI:;443.100NC*013104h3529.51N/07850.14Wrhttp://www.carolina440.net 443.100 MHz PL 100.0 Hz",
	"KD4PBS-3>APN382,WIDE2-1,qAR,K4JH-1:!3540.59NN07832.12W#PHG7760/W2, NCn digi listening 147.39+88.5Hz",
	"W3AHL-2>APTPV1,KD4PBS-3*,WIDE2-1,qAR,K4RAX-10:@270130z3551.25N/07907.52W_000/000g000t072r000p001P001h81b08326.DsVP",
	"NC4LA-4>APN390,WIDE2-2,qAR,N4ILM-4:!3514.85NS07735.90W#PHG7650d/W2,NCn-N Kinston N.C.",
	"KG4AGD>APRS,TCPIP*,qAC,FOURTH:@270131z3521.18N/07833.48W_019/002g004t072r000p000P000h75b10206 VISR3760 400",
	"KC6URO-13>APTW14,KD4PBS-3*,WIDE2-1,qAR,K4JH-1:_04021046c110s000g000t068r000p000P000h..b.....tU2k",
	"KJ4GPT-1>APDW13,TCPIP*,qAC,T2PR:!3530.48NR08019.15W#Kahuna's RASPi DireWolf iGate Gold Hill,NC",
	"W4DEX-2>APU25N,TCPIP*,qAC,FIFTH:@270131z3515.35N/08022.82W_000/000g000t073r000p000P000h74b10177/WX",
}

func TestAddressString(t *testing.T) {
	a := Address{Call: "N0CALL", SSID: 13}
	assert.Equal(t, "N0CALL-13", a.String(), "Address")
}

func TestFrameFromString(t *testing.T) {
	a := assert.New(t)

	for _, s := range tnc2Bad {
		f := Frame{}
		err := f.FromString(s)
		a.NotNil(err, "Invalid frame: "+s)
	}

	for _, s := range tnc2Good {
		f := Frame{}
		err := f.FromString(s)
		a.Nil(err, "Valid frame: "+s)
	}
}

func TestFrameString(t *testing.T) {
	a := assert.New(t)

	f := Frame{}
	f.Src.FromString("N0CALL-13")
	f.Dst.FromString("APZ001")
	f.Path.FromString("WIDE1-1,WIDE2-1")
	f.Text = "Hello world"
	a.Equal("N0CALL-13>APZ001,WIDE1-1,WIDE2-1:Hello world",
		f.String())

	f = Frame{}
	f.FromBytes(ax25Wx1)
	a.Equal("KG4HIE>APK102,W4LBT-9,WIDE1,KD4PBS-3*,WIDE2:=3438.51N/07941.15W_120/001g004t073r   p   P000h  b     KU2k\r",
		f.String())

	f = Frame{}
	f.FromBytes(ax25Wx2)
	a.Equal("N4MTT-2>APX209,KD4PBS-3*,WIDE2-2:@270055z3548.41N/07846.35W_360/000g000t066r000P000p000h63b10183XU2k\r",
		f.String())
}
