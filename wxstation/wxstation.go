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

/*Package wxstation wraps around several i2c and other system devices to
measure and return measured data


There is only one ID that exists*/
package wxstation

import (
	"github.com/npotts/sensors/htu21d"
	"golang.org/x/exp/io/i2c"
	"math"
)

/*Measurement is a weather station measurement*/
type Measurement struct {
	Pressure  float64
	PressT    float64
	Humidity  float64
	HumidityT float64
	Dewpoint  float64
}

func (m *Measurement) clear() {
	m.Pressure = math.NaN()
	m.PressT = math.NaN()
	m.Humidity = math.NaN()
	m.HumidityT = math.NaN()
	m.Dewpoint = math.NaN()
}

/*Station is a weather station*/
type Station struct {
	i2c *i2c.Devfs
	rh  *htu21d.HTU21D
}

/*NewStation returns a intialited Station*/
func NewStation(device string) (*Station, error) {
	dev := &i2c.Devfs{Dev: device}
	rh, err := htu21d.NewHTU21D(dev, false, false, true)
	// rh, err := NewHTU21D(dev, false, false, true)
	if err != nil {
		return nil, err
	}

	return &Station{
		i2c: dev,
		rh:  rh,
	}, nil
}

/*Measure returns some measurments from the hardware*/
func (s *Station) Measure() (m Measurement, e error) {
	m.clear()
	if rh, err := s.rh.Measure(); err == nil {
		m.Humidity, m.HumidityT, m.Dewpoint = rh.Humidity, rh.Temperature, rh.Dewpoint
	}
	return
}
