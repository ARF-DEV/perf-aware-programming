package internal

const (
	MOV_REGISTER_FROM_TO_MEMORY OpMode = iota
	MOV_IMMEDIATE_TO_REGISTER
	MOV_IMMEDIATE_TO_REGISTER_MEMORY
	MOV_ACCULUMATOR_FROM_TO_MEMORY
)

type MovInstruction struct {
	d, w, mod, reg, rm uint8
	lo, hi             int8
	data               int16
	op                 OpMode
}

func (i *MovInstruction) String() string {
	return ""
}

func (i *MovInstruction) isInstruction()

// TODO(arief) TBD: need to add mode map to know what function to use to convert data structure to string format
