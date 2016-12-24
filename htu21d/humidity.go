/*
MIT License

Copyright (c) 2016 Nick Potts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

/*Package HTU21D wraps around a i2c device to control a HTU21D on address 0x40 (fixed).

There is only one ID that exists*/
package htu21d

import (
	"fmt"
	"github.com/sigurn/crc8"
	"golang.org/x/exp/io/i2c"
	"math"
	"time"
)

type htuCommand byte

const (
	htuMeasTempHold = 0xE3 //hold the master line
	htuMeasHumdHold = 0xE5 //hold the master line
	htuMeasTemp     = 0xF3 //dont hold the master line
	htuMeasHumd     = 0xF5 //dont hold the master line
	htuWriteReg     = 0xE6 //write to user registery
	htuReadReg      = 0xE7 //read from user registry
	htuReset        = 0xE7 //soft reset
)

/*HTU21D wraps around a i2c hygrometer*/
type HTU21D struct {
	mode0, mode1 bool
	heater       bool
	oTPReload    bool
	i2c          *i2c.Device
	coeffTemp    float64
	crcTable     *crc8.Table
}

func (h *HTU21D) checkCRC(packet []byte) bool {
	if len(packet) < 1 {
		return false
	}
	return packet[len(packet)-1] == crc8.Checksum(packet[0:len(packet)-1], h.crcTable)
}

/*NewHTU21D returns a new HTU device*/
func NewHTU21D(dev *i2c.Devfs, mode1, mode0, heater bool) (*HTU21D, error) {
	d, err := i2c.Open(dev, 0x40)
	if err != nil {
		return nil, err
	}

	h := &HTU21D{
		mode0:  mode0,
		mode1:  mode1,
		heater: heater,
		i2c:    d,
		crcTable: crc8.MakeTable(crc8.Params{
			Poly:   0x31,
			Init:   0x00,
			RefIn:  false,
			RefOut: false,
			XorOut: 0x00,
			Check:  0xFF,
			Name:   "CRC-HTU21D",
		}),
	}

	/*Init the HTU21D Device.

	Per the datasheet, we need to:
		- retrive the Reserved bits (3, 4, & 5)
		- write the desired settings, keeping the values
	*/
	b2b := func(i bool) byte {
		if i {
			return 1
		}
		return 0
	}
	buf := make([]byte, 1)
	if e := h.i2c.ReadReg(htuReadReg, buf); e != nil {
		return nil, e
	}

	var reg = (b2b(h.mode1) << 7) +
		(buf[0]&0x20)<<5 +
		(buf[0]&0x10)<<4 +
		(buf[0]&0x08)<<3 +
		b2b(h.heater)<<2 +
		b2b(h.oTPReload)<<1 +
		b2b(h.mode0)<<0

	if err := h.i2c.WriteReg(htuWriteReg, []byte{reg}); err != nil {
		return nil, err
	}
	return h, nil
}

//Temp  conversion
func (h *HTU21D) asTemp(sample uint64) float64 {
	return 175.72*float64(sample)/65536.0 - 46.85
}

//Humid conversion
func (h *HTU21D) asHumidity(sample uint64) float64 {
	return 125.0*float64(sample)/65536.0 - 6.0
}

const (
	a = 8.1332
	b = 1762.39
	c = 235.66
)

func (h *HTU21D) partialPressure(tambient float64) float64 {
	return math.Pow(10.0, a-b/(tambient+c))
}

/*dewpoint calculates a dewpoints based off tambient and a RH.*/
func (h *HTU21D) dewpoint(tambient, rh float64) float64 {
	return -(b/(math.Log10(rh*h.partialPressure(tambient)/100.00)-a) + c)
}

/*ReadHumidity reads the humidity, blocking while waiting for values */
func (h *HTU21D) read(cmd byte, backoff time.Duration) (uint64, error) {
	if err := h.i2c.Write([]byte{htuMeasHumd}); err != nil {
		return 0, err
	}
	time.Sleep(backoff) //wait for conversion
	buffer := make([]byte, 3)
	if err := h.i2c.Read(buffer); err != nil || len(buffer) != 3 {
		return 0, err
	}

	val := (uint64(buffer[0])<<8 | (uint64(buffer[1] & 0xFF)))
	if h.checkCRC(buffer) {
		return val, nil
	}
	return val, fmt.Errorf("Checksum missmatch")
}

/*RHMeasurement represents a RH Measurements*/
type RHMeasurement struct {
	Humidity, Temperature, Dewpoint float64
}

/*Measure returns */
func (h *HTU21D) Measure() (m RHMeasurement, err error) {
	rhm := RHMeasurement{math.NaN(), math.NaN(), math.NaN()}
	raw, err := h.read(htuMeasHumd, 75*time.Millisecond)
	if err != nil {
		return rhm, err
	}
	rhm.Humidity = h.asHumidity(raw)

	raw, err = h.read(htuMeasTemp, 75*time.Millisecond)
	if err != nil {
		return rhm, err
	}
	rhm.Temperature = h.asTemp(raw)
	rhm.Dewpoint = h.dewpoint(rhm.Temperature, rhm.Humidity)
	return rhm, nil
}
