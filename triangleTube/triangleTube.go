package triangleTube

import "fmt"

type TriangleTube struct {
	// Modbus
	modBusID int

	// Input Registers
	BoilerStatus           int8
	LockoutStatus1         uint16
	LockoutStatus2         uint16
	BoilerSupplyTemp       uint16
	BoilerReturnTemp       uint16
	BoilerFlueTemp         uint16
	OutdoorTemp            uint16
	FlameIonizationCurrent uint16
	BoilerFiringRate       uint16
	BoilerSetpoint         uint16

	// Holding registers
	ChDemand           uint16
	MaxFiringRate      uint16
	ChSetpoint         uint16
	Ch1MaxSetpoint     uint16
	DhwStorageSetpoint uint16
}

func NewBoiler(ID int) (b *TriangleTube) {
	b = new(TriangleTube)
	// Set initial status
	b.modBusID = ID

	return
}

func (b *TriangleTube) ProcessCommand(c []byte) (r int16, err error) {

	if len(c) < 3 {
		return 0, fmt.Errorf("malformed command")
	}

	command := c[0]
	switch command {
	case readInputRegisters:

	case readHoldingRegisters:
	default:
		return 0, fmt.Errorf("command not supported")
	}

	return

}

func (b *TriangleTube) readInputRegister(r int) (v int16, err error) {
	switch r {
	case boilerSupplyTemp:
		return int16(b.getBoilerSupplyTemp()), nil
	case boilerReturnTemp:
		return int16(b.getBoilerReturnTemp()), nil
	case outdoorTemp:
		return int16(b.getOutdoorTemp()), nil
	default:
		err = fmt.Errorf("command not supported")
		return 0, err
	}
}

func (b *TriangleTube) getBoilerSupplyTemp() uint16 {
	return b.BoilerSupplyTemp
}

func (b *TriangleTube) getBoilerReturnTemp() uint16 {
	return b.BoilerReturnTemp
}

func (b *TriangleTube) getOutdoorTemp() uint16 {
	return b.BoilerSupplyTemp
}
