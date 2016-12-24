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

package mpl3115

import (
	"fmt"
	"golang.org/x/exp/io/i2c"
	"math"
)

/*Baraometer wraps around a MPL3115 strictly as a barometer in raw mode*/
type Baraometer struct {
	i2c *i2c.Device
}

/*NewMPL3115 returns a new MPL3115 configures a barometer device*/
func NewMPL3115(dev *i2c.Devfs) (*Baraometer, error) {
	d, err := i2c.Open(dev, 0x60)
	if err != nil {
		return nil, err
	}

	m := &Baraometer{
		i2c: d,
	}

	/*
		Per AN4519 (http://cache.freescale.com/files/sensors/doc/app_note/AN4519.pdf
		and the datasheet, the MPL3115 should be configured in RAW & Barometer mode
	*/
	cmds := []regCmd{
		{regFifoSetup, 0x00}, //disable fifo mode and watermark
		{regCtrlReg1, 0x79},  //0b01111001
		{regCtrlReg2, 0x00},  //0b00000000
		{regCtrlReg3, 0x00},  //0b00000000
		{regCtrlReg4, 0x00},  //0b00000000
		{regCtrlReg5, 0x00},  //0b00000000
	}
	for _, rg := range cmds {
		if err := m.i2c.WriteReg(rg.reg, []byte{rg.cmd}); err != nil {
			return nil, err
		}
	}

	return m, nil
}

/*Measure measures */
func (b *Baraometer) Measure() (Measurement, error) {
	var err error
	m := Measurement{math.NaN(), math.NaN()}
	if m.PressAlt, err = b.pressure(); err != nil {
		return m, err
	}

	if m.Temperature, err = b.temperature(); err != nil {
		return m, err
	}

	return m, nil
}

func (b *Baraometer) pressure() (float64, error) {
	buf := make([]byte, 6)
	if e := b.i2c.ReadReg(regSenStatus, buf); e != nil {
		return 0, e
	}
	fmt.Println("Buffer:", buf)
	return 1, nil

}

func (b *Baraometer) temperature() (float64, error) {
	return 0.0, nil

}

// //Temp  conversion
// func (m *MPL3115) asTemp(sample uint64) float64 {
// 	return 175.72*float64(sample)/65536.0 - 46.85
// }

// //Pressure conversion
// func (m *MPL3115) asHumidity(sample uint64) float64 {
// 	return 125.0*float64(sample)/65536.0 - 6.0
// }

// const (
// 	a = 8.1332
// 	b = 1762.39
// 	c = 235.66
// )

// func (h *HTU21D) partialPressure(tambient float64) float64 {
// 	return math.Pow(10.0, a-b/(tambient+c))
// }

// /*dewpoint calculates a dewpoints based off tambient and a RH.*/
// func (h *HTU21D) dewpoint(tambient, rh float64) float64 {
// 	return -(b/(math.Log10(rh*h.partialPressure(tambient)/100.00)-a) + c)
// }

// /*ReadHumidity reads the humidity, blocking while waiting for values */
// func (h *HTU21D) read(cmd byte, backoff time.Duration) (uint64, error) {
// 	if err := h.i2c.Write([]byte{htuMeasHumd}); err != nil {
// 		return 0, err
// 	}
// 	time.Sleep(backoff) //wait for conversion
// 	buffer := make([]byte, 3)
// 	if err := h.i2c.Read(buffer); err != nil || len(buffer) != 3 {
// 		return 0, err
// 	}

// 	val := (uint64(buffer[0])<<8 | (uint64(buffer[1] & 0xFF)))
// 	if h.checkCRC(buffer) {
// 		return val, nil
// 	}
// 	return val, fmt.Errorf("Checksum missmatch")
// }

// /*RHMeasurement represents a RH Measurements*/
// type RHMeasurement struct {
// 	Humidity, Temperature, Dewpoint float64
// }

// /*Measure returns */
// func (h *HTU21D) Measure() (m RHMeasurement, err error) {
// 	rhm := RHMeasurement{math.NaN(), math.NaN(), math.NaN()}
// 	raw, err := h.read(htuMeasHumd, 75*time.Millisecond)
// 	if err != nil {
// 		return rhm, err
// 	}
// 	rhm.Humidity = h.asHumidity(raw)

// 	raw, err = h.read(htuMeasTemp, 75*time.Millisecond)
// 	if err != nil {
// 		return rhm, err
// 	}
// 	rhm.Temperature = h.asTemp(raw)
// 	rhm.Dewpoint = h.dewpoint(rhm.Temperature, rhm.Humidity)
// 	return rhm, nil
// }
