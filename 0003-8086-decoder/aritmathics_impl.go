package main

import "fmt"

var aritmaticsOpCodeMap map[byte]string = map[byte]string{
	0b000: "add",
	0b101: "sub",
	0b111: "cmp",
}

func (i *InstructionDecoder) IncRegToMem() string {
	fmt.Println("IRM", i.curIdx)
	w := i.getMaskedBits(0, 0b00000001)
	d := i.getMaskedBits(1, 0b00000001)
	opBinary := i.getMaskedBits(3, 0b00000111)

	i.Next()
	mod := i.getMaskedBits(6, 0b00000011)
	reg := i.getMaskedBits(3, 0b00000111)
	rm := i.getMaskedBits(0, 0b00000111)
	RToR := 0
	if mod == 0b11 {
		RToR = 1
	}

	regStr := RegisterMap[reg|w<<3|1<<4]
	rmStr := RegisterMap[rm|w<<3|byte(RToR)<<4]

	switch mod {
	case 0b00:
		if rmStr == "bp" {
			disp := []int16{}
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			rmStr = fmt.Sprintf("%d", uint16(disp[1]<<8|disp[0]))
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	case 0b01:
		i.Next()
		disp := int8(i.CurrentByte())
		if disp > 0 {
			rmStr = fmt.Sprintf("%s + %d", rmStr, disp)
		} else if disp < 0 {
			rmStr = fmt.Sprintf("%s - %d", rmStr, -disp)
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	case 0b10:
		disp := []int16{}
		i.Next()
		disp = append(disp, int16(i.CurrentByte()))
		i.Next()
		disp = append(disp, int16(i.CurrentByte()))

		dispVal := disp[1]<<8 | disp[0]
		if dispVal > 0 {
			rmStr = fmt.Sprintf("%s + %d", rmStr, dispVal)
		} else if dispVal < 0 {
			rmStr = fmt.Sprintf("%s - %d", rmStr, -dispVal)
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	}

	var op, src, dst string
	if d == 1 {
		src = rmStr
		dst = regStr
	} else {
		src = regStr
		dst = rmStr
	}
	op = aritmaticsOpCodeMap[opBinary]

	return fmt.Sprintf("%s %s, %s", op, dst, src)
}

func (i *InstructionDecoder) ImmediateToRM() string {
	// fmt.Print("ITRM", i.curIdx)
	w := i.getMaskedBits(0, 0b00000001)
	s := i.getMaskedBits(1, 0b00000001)

	i.Next()

	mod := i.getMaskedBits(6, 0b00000011)
	subOpCode := i.getMaskedBits(3, 0b00000111)
	rm := i.getMaskedBits(0, 0b00000111)
	RToR := 0
	if mod == 0b11 {
		RToR = 1
	}

	// regStr := RegisterMap[reg|w<<3|1<<4]
	rmStr := RegisterMap[rm|w<<3|byte(RToR)<<4]

	switch mod {
	case 0b00:
		if rmStr == "bp" {
			disp := []int16{}
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			rmStr = fmt.Sprintf("%d", uint16(disp[1]<<8|disp[0]))
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	case 0b01:
		i.Next()
		disp := int8(i.CurrentByte())
		if disp > 0 {
			rmStr = fmt.Sprintf("%s + %d", rmStr, disp)
		} else if disp < 0 {
			rmStr = fmt.Sprintf("%s - %d", rmStr, -disp)
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	case 0b10:
		disp := []int16{}
		i.Next()
		disp = append(disp, int16(i.CurrentByte()))
		i.Next()
		disp = append(disp, int16(i.CurrentByte()))

		dispVal := disp[1]<<8 | disp[0]
		if dispVal > 0 {
			rmStr = fmt.Sprintf("%s + %d", rmStr, dispVal)
		} else if dispVal < 0 {
			rmStr = fmt.Sprintf("%s - %d", rmStr, -dispVal)
		}
		rmStr = fmt.Sprintf("word [%s]", rmStr)
	}

	dataStr := ""
	switch w {
	case 0b1:
		if s == 1 {
			i.Next()
			data := int16(i.CurrentByte())
			dataStr = fmt.Sprintf("%d", data)
		} else {
			data := []int16{}
			i.Next()
			data = append(data, int16(i.CurrentByte()))
			i.Next()
			data = append(data, int16(i.CurrentByte()))

			dataVal := data[1]<<8 | data[0]
			dataStr = fmt.Sprintf("%d", dataVal)
		}

	case 0b0:
		i.Next()
		data := int8(i.CurrentByte())
		dataStr = fmt.Sprintf("%d", data)
		rmStr = fmt.Sprint("byte ", rmStr)
	}
	op := aritmaticsOpCodeMap[subOpCode]
	// fmt.Print(i.curIdx)
	return fmt.Sprintf("%s %s, %s", op, rmStr, dataStr)
}

func (i *InstructionDecoder) ImmediateToAcc() string {
	fmt.Println("IA", i.curIdx)
	w := i.getMaskedBits(0, 0b00000001)
	opBinary := i.getMaskedBits(3, 0b00000111)
	reg := "ax"
	dataStr := ""
	switch w {
	case 0b1:
		data := []int16{}
		i.Next()
		data = append(data, int16(i.CurrentByte()))
		i.Next()
		data = append(data, int16(i.CurrentByte()))

		dataVal := data[1]<<8 | data[0]
		dataStr = fmt.Sprintf("%d", dataVal)
	case 0b0:
		reg = "al"
		i.Next()
		data := int8(i.CurrentByte())
		dataStr = fmt.Sprintf("%d", data)
	}

	op := aritmaticsOpCodeMap[opBinary]

	return fmt.Sprintf("%s %s, %s", op, reg, dataStr)
}
