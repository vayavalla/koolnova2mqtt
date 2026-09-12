package modbus

import (
	"errors"
)

type Mock struct {
	State map[byte][]uint16
}

func NewMock() *Mock {
	ms := &Mock{
		State: map[byte][]uint16{
			49: {3, 68, 41, 41, 3, 68, 41, 41, 3, 68, 41, 45, 3, 68, 41, 45, 3, 68, 41, 42, 3, 52, 41, 40, 3, 52, 41, 44, 3, 68, 41, 41, 3, 68, 41, 40, 3, 68, 41, 41, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 3, 4, 2, 49, 3, 7, 1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			50: {3, 68, 41, 41, 3, 68, 41, 41, 3, 68, 41, 45, 3, 68, 41, 45, 3, 68, 41, 42, 3, 52, 41, 40, 3, 52, 41, 44, 3, 68, 41, 41, 3, 68, 41, 40, 3, 68, 41, 41, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 68, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 3, 4, 2, 49, 3, 7, 1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	// Koolnova 2.0 watches system registers 40065..40126.
	// The original test fixture only contained registers 40001..40094.
	for slaveID, state := range ms.State {
		if len(state) < 126 {
			state = append(state, make([]uint16, 126-len(state))...)
		}

		// Registers 40065..40126 belong to the Koolnova 2.0
		// system register map. Do not reuse the old V1 fixture
		// values that happened to occupy this range.
		for reg := 65; reg <= 126; reg++ {
			state[reg-1] = 0
		}

		// Representative Koolnova 2.0 system state.
		state[74-1] = 3 << 8 // 40074: EFI=3 in bits 10..8
		state[75-1] = 15     // 40075: fan + cool + heat + dry available
		state[76-1] = 12838  // 40076: maxHeat=25, minCool=19
		state[77-1] = 16713  // 40077: MAU modes + humidity threshold 73

		// AC1..AC4 airflow
		state[93-1] = 0
		state[94-1] = 0
		state[95-1] = 2 // AC3
		state[96-1] = 0

		// AC1..AC4 target temperature (x2)
		state[97-1] = 0
		state[98-1] = 0
		state[99-1] = 46 // AC3 = 23 °C
		state[100-1] = 0

		// AC1..AC4 fan mode
		state[101-1] = 4 // auto
		state[102-1] = 4 // auto
		state[103-1] = 3 // AC3 high
		state[104-1] = 4 // auto

		state[105-1] = 2               // 40105: 9600 baud, even parity
		state[106-1] = uint16(slaveID) // 40106: Modbus slave ID

		state[109-1] = 1 // 40109: global system enabled
		state[110-1] = 1 // 40110: global mode = cool

		state[111-1] = 0 // floor demand thermostats
		state[112-1] = 5 // AC3 demand thermostats

		// Connected volume AC1..AC4
		state[115-1] = 6 // AC3

		// Demand volume AC1..AC4
		state[119-1] = 1 // AC3

		// Average target temperature AC1..AC4
		state[123-1] = 46 // AC3 = 23 °C

		// EFI in MSB, AC3 speed in LSB: 0x0303
		state[125-1] = 771

		ms.State[slaveID] = state
	}

	return ms
}

func (ms *Mock) ReadRegister(slaveID byte, address uint16, quantity uint16) (results []uint16, err error) {
	state, ok := ms.State[slaveID]
	if !ok {
		return nil, errors.New("Unknown slave")
	}
	address--
	for a := address; a < address+quantity; a++ {
		results = append(results, state[a])
	}
	return results, nil
}
func (ms *Mock) WriteRegister(slaveID byte, address uint16, value uint16) (results []uint16, err error) {
	state, ok := ms.State[slaveID]
	if !ok {
		return nil, errors.New("Unknown slave")
	}
	state[address-1] = value
	return []uint16{value}, nil
}

func (ms *Mock) Close() error { return nil }
