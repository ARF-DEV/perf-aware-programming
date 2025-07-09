package internal

import "fmt"

func (i *InstructionDecoder) MovIRegInstruction() string {
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
		if rmStr == "bp" {
			disp := []int16{}
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			i.Next()
			disp = append(disp, int16(i.CurrentByte()))
			rmStr = fmt.Sprintf("%d", uint16(disp[1]<<8|disp[0]))
		}
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

func (i *InstructionDecoder) MovInstruction() string {
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

	return fmt.Sprintf("mov %s, %s", dst, src)
}
