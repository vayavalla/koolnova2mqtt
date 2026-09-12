package kn

import (
	"errors"
)

const NUM_ZONES = 16

const REG_PER_ZONE = 4
const REG_ENABLED = 1
const REG_MODE = 2
const REG_TARGET_TEMP = 3
const REG_CURRENT_TEMP = 4

// Koolnova Modbus 2.0
//
// Zone registers 40001..40064 are unchanged.
// System registers move in Modbus 2.0.

const REG_AIRFLOW = 93
const REG_AC_TARGET_TEMP = 97
const REG_AC_TARGET_FAN_MODE = 101
const REG_SERIAL_CONFIG = 105
const REG_SLAVE_ID = 106
const REG_EFFICIENCY = 74
const REG_SYSTEM_ENABLED = 109
const REG_SYS_KN_MODE = 110

// Koolnova 2.0 advanced system registers.
const REG_ACTIVE_MODES = 75
const REG_TEMP_LIMITS = 76
const REG_AUTO_CHANGE = 77
const REG_WATER_TEMP = 82
const REG_OUTDOOR_TEMP = 83
const REG_AUX_TEMP = 84
const REG_FLOOR_DEMAND = 111
const REG_AC3_DEMAND = 112
const REG_CONNECTED_VOLUME = 113
const REG_DEMAND_VOLUME = 117
const REG_AVG_TARGET_AC1 = 121
const REG_AVG_TARGET_AC2 = 122
const REG_AVG_TARGET_AC3 = 123
const REG_AVG_TARGET_AC4 = 124 // verified on real Koolnova 2.0 hardware
const REG_EFI_SPEED_AC3 = 125  // verified on real Koolnova 2.0 hardware

const FIRST_ZONE_REGISTER = REG_ENABLED
const TOTAL_ZONE_REGISTERS = NUM_ZONES * REG_PER_ZONE
const FIRST_SYS_REGISTER = 65
const TOTAL_SYS_REGISTERS = 62 // 40065..40126

type FanMode byte

const FAN_OFF FanMode = 0
const FAN_LOW FanMode = 1
const FAN_MED FanMode = 2
const FAN_HIGH FanMode = 3
const FAN_AUTO FanMode = 4

type KnMode byte

const MODE_AIR_VENTILATION KnMode = 0x00
const MODE_AIR_COOLING KnMode = 0x01
const MODE_AIR_HEATING KnMode = 0x02
const MODE_DEHUMIDIFICATION KnMode = 0x03
const MODE_UNDERFLOOR_HEATING KnMode = 0x04
const MODE_UNDERFLOOR_AIR_COOLING KnMode = 0x05
const MODE_UNDERFLOOR_AIR_HEATING KnMode = 0x06

const HOLD_MODE_UNDERFLOOR_ONLY = "underfloor"
const HOLD_MODE_FAN_ONLY = "fan"
const HOLD_MODE_UNDERFLOOR_AND_FAN = "underfloor and fan"

const HVAC_MODE_OFF = "off"
const HVAC_MODE_COOL = "cool"
const HVAC_MODE_HEAT = "heat"
const HVAC_MODE_DRY = "dry"
const HVAC_MODE_FAN_ONLY = "fan_only"

type ACMachine int

const ACMachines = 4

const AC1 ACMachine = 1
const AC2 ACMachine = 2
const AC3 ACMachine = 3
const AC4 ACMachine = 4

const HA_COMPONENT_SENSOR = "sensor"
const HA_COMPONENT_CLIMATE = "climate"

func FanMode2Str(fm FanMode) string {
	switch fm {
	case FAN_OFF:
		return "off"
	case FAN_LOW:
		return "low"
	case FAN_MED:
		return "medium"
	case FAN_HIGH:
		return "high"
	case FAN_AUTO:
		return "auto"
	default:
		return "unknown"
	}
}

func Str2FanMode(st string) (FanMode, error) {
	switch st {
	case "off":
		return FAN_OFF, nil
	case "low":
		return FAN_LOW, nil
	case "medium":
		return FAN_MED, nil
	case "high":
		return FAN_HIGH, nil
	case "auto":
		return FAN_AUTO, nil
	default:
		return FAN_OFF, errors.New("Unknown fan mode")
	}
}

func ApplyHvacMode(knMode KnMode, hvacMode string) KnMode {
	switch hvacMode {
	case HVAC_MODE_FAN_ONLY:
		return MODE_AIR_VENTILATION
	case HVAC_MODE_DRY:
		return MODE_DEHUMIDIFICATION
	case HVAC_MODE_COOL:
		switch knMode {
		case MODE_UNDERFLOOR_HEATING, MODE_UNDERFLOOR_AIR_HEATING, MODE_UNDERFLOOR_AIR_COOLING:
			return MODE_UNDERFLOOR_AIR_COOLING
		default:
			return MODE_AIR_COOLING
		}
	case HVAC_MODE_HEAT:
		switch knMode {
		case MODE_UNDERFLOOR_HEATING:
			return MODE_UNDERFLOOR_HEATING
		case MODE_UNDERFLOOR_AIR_COOLING, MODE_UNDERFLOOR_AIR_HEATING:
			return MODE_UNDERFLOOR_AIR_HEATING
		default:
			return MODE_AIR_HEATING
		}
	default:
		return knMode
	}
}

func KnMode2Str(knMode KnMode) string {
	switch knMode {
	case MODE_AIR_VENTILATION:
		return "fan"
	case MODE_AIR_COOLING:
		return "cool"
	case MODE_AIR_HEATING:
		return "heat"
	case MODE_DEHUMIDIFICATION:
		return "dry"
	case MODE_UNDERFLOOR_HEATING:
		return "underfloor"
	case MODE_UNDERFLOOR_AIR_COOLING:
		return "underfloorCool"
	case MODE_UNDERFLOOR_AIR_HEATING:
		return "underfloorHeat"
	default:
		return "unknown"
	}
}

func ApplyHoldMode(knMode KnMode, holdMode string) KnMode {
	cool := knMode == MODE_AIR_COOLING || knMode == MODE_UNDERFLOOR_AIR_COOLING
	switch holdMode {
	case HOLD_MODE_FAN_ONLY:
		if cool {
			return MODE_AIR_COOLING
		}
		return MODE_AIR_HEATING
	case HOLD_MODE_UNDERFLOOR_ONLY:
		if cool {
			return MODE_UNDERFLOOR_AIR_COOLING
		}
		return MODE_UNDERFLOOR_HEATING
	case HOLD_MODE_UNDERFLOOR_AND_FAN:
		if cool {
			return MODE_UNDERFLOOR_AIR_COOLING
		}
		return MODE_UNDERFLOOR_AIR_HEATING
	}
	return knMode
}
