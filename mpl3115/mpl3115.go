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

/*Package mpl3115 wraps around a i2c device to control devices in the MPL3115A.

Default address, which can be changed by the factory, is .

There is only one ID that exists*/
package mpl3115

import ()

const (
	regSenStatus = 0x00
	regPresMSB   = 0x01
	regPresCSB   = 0x02
	regPresLSB   = 0x03
	regTempMSB   = 0x04
	regTempLSB   = 0x05
	regFifoSetup = 0x0F
	regCtrlReg1  = 0x26
	regCtrlReg2  = 0x27
	regCtrlReg3  = 0x28
	regCtrlReg4  = 0x29
	regCtrlReg5  = 0x2A
	regToffset   = 0x2B
	regPoffset   = 0x2B
)

type regCmd struct {
	reg, cmd byte
}

/*Measurement is a Pressure/Altitude and Temp measurement
off of the component in either altimeter or barometer mode*/
type Measurement struct {
	PressAlt, Temperature float64
}
