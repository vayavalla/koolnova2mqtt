package kn

import "errors"

type SysConfig struct {
	Watcher Watcher
}

// SysDriver watches system registers and allows
// to read/change the module configuration
// notifies callbacks when specific system registers change
type SysDriver struct {
	SysConfig
	OnACAirflowChange          func(ac ACMachine)
	OnACTargetTempChange       func(ac ACMachine)
	OnACTargetFanModeChange    func(ac ACMachine)
	OnEfficiencyChange         func()
	OnSystemEnabledChange      func()
	OnKnModeChange             func()
	OnActiveModesChange        func()
	OnTemperatureLimitsChange  func()
	OnAutoChangeChange         func()
	OnWaterTemperatureChange   func()
	OnOutdoorTemperatureChange func()
	OnAuxTemperatureChange     func()
	OnDemandChange             func()
	OnACConnectedVolumeChange  func(ac ACMachine)
	OnACDemandVolumeChange     func(ac ACMachine)
	OnACAverageTargetChange    func(ac ACMachine)
	OnAC3StatusChange          func()
}

var ErrUnknownSerialConfig = errors.New("Uknown serial configuration")

func NewSys(config *SysConfig) *SysDriver {
	s := &SysDriver{
		SysConfig: *config,
	}
	for n := byte(0); n < ACMachines; n++ {
		func(ac ACMachine) {
			s.Watcher.RegisterCallback(uint16(REG_AIRFLOW+int(ac)-1), func(address uint16) {
				if s.OnACAirflowChange != nil {
					s.OnACAirflowChange(ac)
				}
			})
			s.Watcher.RegisterCallback(uint16(REG_AC_TARGET_TEMP+int(ac)-1), func(address uint16) {
				if s.OnACTargetTempChange != nil {
					s.OnACTargetTempChange(ac)
				}
			})
			s.Watcher.RegisterCallback(uint16(REG_AC_TARGET_FAN_MODE+int(ac)-1), func(address uint16) {
				if s.OnACTargetFanModeChange != nil {
					s.OnACTargetFanModeChange(ac)
				}
			})
		}(ACMachine(n + 1))
	}

	s.Watcher.RegisterCallback(REG_EFFICIENCY, func(address uint16) {
		if s.OnEfficiencyChange != nil {
			s.OnEfficiencyChange()
		}
	})

	s.Watcher.RegisterCallback(REG_SYSTEM_ENABLED, func(address uint16) {
		if s.OnSystemEnabledChange != nil {
			s.OnSystemEnabledChange()
		}
	})

	s.Watcher.RegisterCallback(REG_SYS_KN_MODE, func(address uint16) {
		if s.OnKnModeChange != nil {
			s.OnKnModeChange()
		}
	})

	s.Watcher.RegisterCallback(REG_ACTIVE_MODES, func(address uint16) {
		if s.OnActiveModesChange != nil {
			s.OnActiveModesChange()
		}
	})
	s.Watcher.RegisterCallback(REG_TEMP_LIMITS, func(address uint16) {
		if s.OnTemperatureLimitsChange != nil {
			s.OnTemperatureLimitsChange()
		}
	})
	s.Watcher.RegisterCallback(REG_AUTO_CHANGE, func(address uint16) {
		if s.OnAutoChangeChange != nil {
			s.OnAutoChangeChange()
		}
	})
	s.Watcher.RegisterCallback(REG_WATER_TEMP, func(address uint16) {
		if s.OnWaterTemperatureChange != nil {
			s.OnWaterTemperatureChange()
		}
	})
	s.Watcher.RegisterCallback(REG_OUTDOOR_TEMP, func(address uint16) {
		if s.OnOutdoorTemperatureChange != nil {
			s.OnOutdoorTemperatureChange()
		}
	})
	s.Watcher.RegisterCallback(REG_AUX_TEMP, func(address uint16) {
		if s.OnAuxTemperatureChange != nil {
			s.OnAuxTemperatureChange()
		}
	})
	s.Watcher.RegisterCallback(REG_FLOOR_DEMAND, func(address uint16) {
		if s.OnDemandChange != nil {
			s.OnDemandChange()
		}
	})
	s.Watcher.RegisterCallback(REG_AC3_DEMAND, func(address uint16) {
		if s.OnDemandChange != nil {
			s.OnDemandChange()
		}
	})

	for n := byte(0); n < ACMachines; n++ {
		func(ac ACMachine) {
			s.Watcher.RegisterCallback(uint16(REG_CONNECTED_VOLUME+int(ac)-1), func(address uint16) {
				if s.OnACConnectedVolumeChange != nil {
					s.OnACConnectedVolumeChange(ac)
				}
			})
			s.Watcher.RegisterCallback(uint16(REG_DEMAND_VOLUME+int(ac)-1), func(address uint16) {
				if s.OnACDemandVolumeChange != nil {
					s.OnACDemandVolumeChange(ac)
				}
			})
			s.Watcher.RegisterCallback(uint16(s.averageTargetRegister(ac)), func(address uint16) {
				if s.OnACAverageTargetChange != nil {
					s.OnACAverageTargetChange(ac)
				}
			})
		}(ACMachine(n + 1))
	}

	s.Watcher.RegisterCallback(REG_EFI_SPEED_AC3, func(address uint16) {
		if s.OnAC3StatusChange != nil {
			s.OnAC3StatusChange()
		}
	})

	return s
}

func (s *SysDriver) ReadRegister(n int) int {
	r := s.Watcher.ReadRegister(uint16(n))
	return int(r)
}

func (s *SysDriver) WriteRegister(n int, value uint16) error {
	return s.Watcher.WriteRegister(uint16(n), value)
}

func (s *SysDriver) GetAirflow(ac ACMachine) int {
	r := s.ReadRegister(REG_AIRFLOW + int(ac) - 1)
	return r
}

func (s *SysDriver) GetMachineTargetTemp(ac ACMachine) float32 {
	r := s.ReadRegister(REG_AC_TARGET_TEMP + int(ac) - 1)
	return reg2temp(uint16(r))
}

func (s *SysDriver) GetTargetFanMode(ac ACMachine) FanMode {
	r := s.ReadRegister(REG_AC_TARGET_FAN_MODE + int(ac) - 1)
	return FanMode(r)
}

func (s *SysDriver) GetActiveModesRaw() uint16       { return uint16(s.ReadRegister(REG_ACTIVE_MODES)) }
func (s *SysDriver) ActiveModeEnabled(bit uint) bool { return s.GetActiveModesRaw()&(1<<bit) != 0 }

func (s *SysDriver) GetTemperatureLimitsRaw() uint16 { return uint16(s.ReadRegister(REG_TEMP_LIMITS)) }
func (s *SysDriver) GetMaxHeatTemperature() float32 {
	return float32((s.GetTemperatureLimitsRaw()>>8)&0xff) / 2.0
}
func (s *SysDriver) GetMinCoolTemperature() float32 {
	return float32(s.GetTemperatureLimitsRaw()&0xff) / 2.0
}

func (s *SysDriver) GetAutoChangeRaw() uint16  { return uint16(s.ReadRegister(REG_AUTO_CHANGE)) }
func (s *SysDriver) GetAutoAboveMode() KnMode  { return KnMode((s.GetAutoChangeRaw() >> 12) & 0x0f) }
func (s *SysDriver) GetAutoBelowMode() KnMode  { return KnMode((s.GetAutoChangeRaw() >> 8) & 0x0f) }
func (s *SysDriver) GetHumidityThreshold() int { return int(s.GetAutoChangeRaw() & 0xff) }

func (s *SysDriver) GetWaterTemperatureRaw() uint16 { return uint16(s.ReadRegister(REG_WATER_TEMP)) }
func (s *SysDriver) GetWaterTemperature() float32   { return float32(s.GetWaterTemperatureRaw()) / 10.0 }
func (s *SysDriver) GetOutdoorTemperatureRaw() uint16 {
	return uint16(s.ReadRegister(REG_OUTDOOR_TEMP))
}
func (s *SysDriver) GetOutdoorTemperature() float32 {
	return float32(int16(s.GetOutdoorTemperatureRaw())) / 10.0
}
func (s *SysDriver) GetAuxTemperatureRaw() uint16 { return uint16(s.ReadRegister(REG_AUX_TEMP)) }
func (s *SysDriver) GetAuxTemperature() float32 {
	return float32(int16(s.GetAuxTemperatureRaw())) / 10.0
}

func (s *SysDriver) GetFloorDemand() int { return s.ReadRegister(REG_FLOOR_DEMAND) }
func (s *SysDriver) GetAC3Demand() int   { return s.ReadRegister(REG_AC3_DEMAND) }
func (s *SysDriver) GetConnectedVolume(ac ACMachine) int {
	return s.ReadRegister(REG_CONNECTED_VOLUME + int(ac) - 1)
}
func (s *SysDriver) GetDemandVolume(ac ACMachine) int {
	return s.ReadRegister(REG_DEMAND_VOLUME + int(ac) - 1)
}

func (s *SysDriver) averageTargetRegister(ac ACMachine) int {
	switch ac {
	case AC1:
		return REG_AVG_TARGET_AC1
	case AC2:
		return REG_AVG_TARGET_AC2
	case AC3:
		return REG_AVG_TARGET_AC3
	case AC4:
		return REG_AVG_TARGET_AC4
	default:
		return REG_AVG_TARGET_AC1
	}
}
func (s *SysDriver) GetAverageTargetRaw(ac ACMachine) uint16 {
	return uint16(s.ReadRegister(s.averageTargetRegister(ac)))
}
func (s *SysDriver) GetAverageTargetTemperature(ac ACMachine) float32 {
	return float32(s.GetAverageTargetRaw(ac)) / 2.0
}
func (s *SysDriver) GetAC3Efficiency() int {
	return int((uint16(s.ReadRegister(REG_EFI_SPEED_AC3)) >> 8) & 0xff)
}
func (s *SysDriver) GetAC3Speed() int { return int(uint16(s.ReadRegister(REG_EFI_SPEED_AC3)) & 0xff) }

func (s *SysDriver) GetBaudRate() int {
	r := s.ReadRegister(REG_SERIAL_CONFIG)
	switch r {
	case 2, 6:
		return 9600
	case 3, 7:
		return 19200
	}
	return 0
}

func (s *SysDriver) GetParity() string {
	r := s.ReadRegister(REG_SERIAL_CONFIG)
	switch r {
	case 2, 3:
		return "even"
	case 6, 7:
		return "none"
	}
	return "unknown"
}

func (s *SysDriver) GetSlaveID() int {
	r := s.ReadRegister(REG_SLAVE_ID)
	return r
}

func (s *SysDriver) GetEfficiency() int {
	r := s.ReadRegister(REG_EFFICIENCY)

	// Koolnova 2.0 register 40074:
	// bits 10..8 = EFI (efficiency level)
	return int((uint16(r) >> 8) & 0x07)
}

func (s *SysDriver) GetSystemEnabled() bool {
	r := s.ReadRegister(REG_SYSTEM_ENABLED)
	return r != 0
}

// SetSystemEnabled changes the global system state.
// Koolnova 2.0 register 40109:
// true = system enabled
// false = global lock / machines off
func (s *SysDriver) SetSystemEnabled(enabled bool) error {
	var value uint16
	if enabled {
		value = 1
	}
	return s.WriteRegister(REG_SYSTEM_ENABLED, value)
}

func (s *SysDriver) GetSystemKNMode() KnMode {
	r := s.ReadRegister(REG_SYS_KN_MODE)
	return KnMode(r)
}

func (s *SysDriver) SetSystemKNMode(knMode KnMode) error {
	return s.WriteRegister(REG_SYS_KN_MODE, uint16(knMode))
}

// HVACMode returns the HA HVAC mode based on the
// module state
func (s *SysDriver) HVACMode() string {
	if !s.GetSystemEnabled() {
		return HVAC_MODE_OFF
	}
	switch s.GetSystemKNMode() {
	case MODE_AIR_VENTILATION:
		return HVAC_MODE_FAN_ONLY
	case MODE_AIR_COOLING, MODE_UNDERFLOOR_AIR_COOLING:
		return HVAC_MODE_COOL
	case MODE_AIR_HEATING, MODE_UNDERFLOOR_HEATING, MODE_UNDERFLOOR_AIR_HEATING:
		return HVAC_MODE_HEAT
	case MODE_DEHUMIDIFICATION:
		return HVAC_MODE_DRY
	}
	return "unknown"
}

// HoldMode returns the HA Hold Mode based on the
// module state
func (s *SysDriver) HoldMode() string {
	switch s.GetSystemKNMode() {
	case MODE_AIR_COOLING, MODE_AIR_HEATING:
		return HOLD_MODE_FAN_ONLY
	case MODE_UNDERFLOOR_HEATING:
		return HOLD_MODE_UNDERFLOOR_ONLY
	case MODE_UNDERFLOOR_AIR_COOLING, MODE_UNDERFLOOR_AIR_HEATING:
		return HOLD_MODE_UNDERFLOOR_AND_FAN
	}
	return "unknown"
}
