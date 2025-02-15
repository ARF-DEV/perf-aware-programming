package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InstructionFunc func() string
type InstructionInfo struct {
	Length uint8
	Ins    InstructionFunc
}

// key = {w}{register} -> [1bit][3bit]
// TBD: key = {RtoR}{w}{Reg/RM} -> [1bit][1bit][3bit]
var RegisterMap map[byte]string = map[byte]string{
	// MOD = 11
	0b10000: "al",
	0b10001: "cl",
	0b10010: "dl",
	0b10011: "bl",
	0b10100: "ah",
	0b10101: "ch",
	0b10110: "dh",
	0b10111: "bh",

	0b11000: "ax",
	0b11001: "cx",
	0b11010: "dx",
	0b11011: "bx",
	0b11100: "sp",
	0b11101: "bp",
	0b11110: "si",
	0b11111: "di",

	// the others i.e :
	// MOD = 00
	// MOD = 10
	// MOD = 01

	0b00000: "bx + si",
	0b00001: "bx + di",
	0b00010: "bp + si",
	0b00011: "bp + di",
	0b00100: "si",
	0b00101: "di",
	0b00110: "bp",
	0b00111: "bx",

	0b01000: "bx + si",
	0b01001: "bx + di",
	0b01010: "bp + si",
	0b01011: "bp + di",
	0b01100: "si",
	0b01101: "di",
	0b01110: "bp",
	0b01111: "bx",
}

type InstructionDecoder struct {
	buf     []byte
	curIdx  int64
	nextIdx int64
	builder strings.Builder

	instructionFuncs map[byte]InstructionFunc
}

func NewDecoder(filename string) *InstructionDecoder {
	ins := InstructionDecoder{}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	ins.buf = make([]byte, info.Size())
	_, err = r.Read(ins.buf)
	if err != nil {
		panic(err)
	}

	ins.initMap()
	ins.builder = strings.Builder{}
	ins.builder.WriteString("bits 16\n\n")
	return &ins
}

func (i *InstructionDecoder) initMap() {
	i.instructionFuncs = map[byte]InstructionFunc{
		0b100010:  i.MovInstruction,
		0b1011:    i.MovIRegInstruction,
		0b1100011: i.MovIRMInstruction,
		0b1010000: i.MovMToAcc,
		0b1010001: i.MovAccToM,
	}
}

func (i *InstructionDecoder) MovIRegInstruction() string {
	// fmt.Println("MOV IREG")
	reg := i.getMaskedBits(0, 0b00000111)
	w := i.getMaskedBits(3, 0b00000001)

	regStr := RegisterMap[reg|w<<3|1<<4]
	var data int16
	switch w {
	case 0b1:
		var b []uint16
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		i.Next()
		b = append(b, uint16(i.CurrentByte()))

		data = int16((b[1] << 8) | b[0])
	case 0b0:
		i.Next()
		// convert properly to int8 first then convert it to 16 bit
		data = int16(int8(i.CurrentByte()))
	}
	return fmt.Sprintf("mov %s, %d", regStr, data)
}

func (i *InstructionDecoder) MovIRMInstruction() string {
	// fmt.Println("MOV RRM")
	w := i.getMaskedBits(0, 0b00000001)

	i.Next()
	mod := i.getMaskedBits(6, 0b00000011)
	RtR := 0
	if mod&0b11 == 0b11 {
		RtR = 1
	}

	rm := i.getMaskedBits(0, 0b00000111)
	rmStr := RegisterMap[rm|(w<<3)|byte(RtR)<<4]

	switch mod {
	case 0b00:
		rmStr = fmt.Sprint("[", rmStr, "]")
	case 0b01:
		i.Next()
		var dispLO int8
		dispLO = int8(i.CurrentByte())
		if dispLO > 0 {
			rmStr = fmt.Sprintf("%s + %d", rmStr, dispLO)
		} else if dispLO < 0 {
			rmStr = fmt.Sprintf("%s - %d", rmStr, -dispLO)
		}
		rmStr = fmt.Sprint("[", rmStr, "]")
	case 0b10:
		var disp []uint16
		i.Next()
		disp = append(disp, uint16(i.CurrentByte()))
		i.Next()
		disp = append(disp, uint16(i.CurrentByte()))

		dispVal := int16(disp[1]<<8 | disp[0])
		if dispVal > 0 {
			rmStr = fmt.Sprintf("%s + %d", rmStr, dispVal)
		} else if dispVal < 0 {
			rmStr = fmt.Sprintf("%s - %d", rmStr, -dispVal)
		}

		rmStr = fmt.Sprint("[", rmStr, "]")
	}

	switch w {
	case 0b1:
		var b []uint16
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		return fmt.Sprintf("mov %s, word %d", rmStr, int16((b[1]<<8)|b[0]))
	case 0b0:
		i.Next()
		return fmt.Sprintf("mov %s, byte %d", rmStr, int8(i.CurrentByte()))
	}
	return ""
}
func (i *InstructionDecoder) MovMToAcc() string {
	w := i.getMaskedBits(0, 0b00000001)
	regStr := "ax"
	res := ""
	switch w {
	case 0b1:
		var b []uint16
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		res = fmt.Sprintf("mov %s, [%d]", regStr, int16((b[1]<<8)|b[0]))
	case 0b0:
		i.Next()
		res = fmt.Sprintf("mov %s, [%d]", regStr, int8(i.CurrentByte()))
	}
	return res
}
func (i *InstructionDecoder) MovAccToM() string {
	w := i.getMaskedBits(0, 0b00000001)
	regStr := "ax"
	res := ""
	switch w {
	case 0b1:
		var b []uint16
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		res = fmt.Sprintf("mov [%d], %s", int16((b[1]<<8)|b[0]), regStr)
	case 0b0:
		i.Next()
		res = fmt.Sprintf("mov [%d], %s", int8(i.CurrentByte()), regStr)
	}
	return res
}
func (i *InstructionDecoder) Decode() string {
	// i.printBytes()
	for i.Next() {
		b := i.CurrentByte()
		for j := 0; j < 8; j++ {
			ins, ok := i.instructionFuncs[b>>j]
			if !ok {
				continue
			}
			res := ins()
			i.builder.WriteString(res + "\n")
			break
		}
	}
	return i.builder.String()
}
func (i *InstructionDecoder) Next() bool {
	i.curIdx = i.nextIdx
	i.nextIdx++
	return int(i.curIdx) < len(i.buf)
}
func (i *InstructionDecoder) CurrentByte() byte {
	return i.buf[i.curIdx]
}
func (i *InstructionDecoder) NextByte() byte {
	return i.buf[i.nextIdx]
}

func (i *InstructionDecoder) MovInstruction() string {
	// fmt.Printf("%08b %08b\n", i.CurrentByte(), i.NextByte())
	// fmt.Println("MOV")
	w := i.getWide()
	d := i.isDestination()

	i.Next()

	mod := i.getMod()
	RtR := 0
	if mod&0b11 == 0b11 {
		RtR = 1
	}

	Reg := i.getReg()
	RegStr := RegisterMap[Reg|(w<<3)|byte(1)<<4]

	RM := i.getRM()
	RMStr := RegisterMap[RM|(w<<3)|(byte(RtR)<<4)]

	switch mod {
	case 0b00:
		if RMStr == "bp" {
			disp := []int16{}
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			RMStr = fmt.Sprintf("%d", uint16(disp[1]<<8|disp[0]))
		}
		RMStr = fmt.Sprint("[", RMStr, "]")
	case 0b01:
		i.Next()
		var dispLO int8
		dispLO = int8(i.CurrentByte())
		if dispLO > 0 {
			RMStr = fmt.Sprintf("%s + %d", RMStr, dispLO)
		} else if dispLO < 0 {
			RMStr = fmt.Sprintf("%s - %d", RMStr, -dispLO)
		}
		RMStr = fmt.Sprint("[", RMStr, "]")
	case 0b10:
		var disp []uint16
		i.Next()
		disp = append(disp, uint16(i.CurrentByte()))
		i.Next()
		disp = append(disp, uint16(i.CurrentByte()))

		dispVal := int16(disp[1]<<8 | disp[0])

		if dispVal > 0 {
			RMStr = fmt.Sprintf("%s + %d", RMStr, dispVal)
		} else if dispVal < 0 {
			RMStr = fmt.Sprintf("%s - %d", RMStr, -dispVal)
		}
		RMStr = fmt.Sprint("[", RMStr, "]")
	default:
	}

	var dst, src string
	if d {
		dst = RegStr
		src = RMStr
	} else {
		src = RegStr
		dst = RMStr
	}

	// fmt.Printf("%08b(%b)(%b)(%b), %08b(%b)(%b)(%b) - %b\n\n", Reg|(w<<3)|byte(RtR)<<4, RtR, w, Reg, RM|(w<<3)|byte(RtR)<<4, RtR, w, RM, mod)
	return fmt.Sprintf("mov %s, %s", dst, src)
}

func (i *InstructionDecoder) isDestination() bool {
	return (i.CurrentByte()>>1)&1 == 1
}

func (i *InstructionDecoder) isWide() bool {
	return i.getWide() == 1
}
func (i *InstructionDecoder) getWide() byte {
	return (i.CurrentByte() & 1)
}

func (i *InstructionDecoder) getReg() byte {
	return (i.CurrentByte() >> 3) & 0b00000111
}
func (i *InstructionDecoder) getRM() byte {
	return i.CurrentByte() & 0b00000111
}

func (i *InstructionDecoder) getMod() byte {
	return (i.CurrentByte() >> 6) & 0b00000011
}

func (i *InstructionDecoder) getMaskedBits(shift int64, mask byte) byte {
	return (i.CurrentByte() >> shift) & mask
}

func (i *InstructionDecoder) printBytes() {
	fmt.Printf("%08b\n", i.buf)
}
