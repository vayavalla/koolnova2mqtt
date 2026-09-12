[![Build and test](https://github.com/vayavalla/koolnova2mqtt/actions/workflows/build-test.yml/badge.svg?branch=koolnova2mqtt-2.0)](https://github.com/vayavalla/koolnova2mqtt/actions/workflows/build-test.yml)

# koolnova2mqtt bridge

![homeassistant](gr/ha.png)

![openhab](gr/openhab.png)

![koolnova](gr/koolnova.png)

***DISCLAIMER: This is not a Koolnova official product. Use at your own risk.***


## Koolnova 2.0 support

This fork adds support for the **Koolnova 2.0 Modbus register map** and has
been tested against a real Koolnova 2.0 installation.

The original project by [Javier Peletier](https://github.com/jpeletier/koolnova2mqtt)
was used as the base for this work.

Main changes in this fork include:

* Koolnova 2.0 system register mapping.
* Global system ON/OFF control.
* Global HVAC mode control.
* Independent ON/OFF control for each zone.
* Additional Koolnova 2.0 diagnostic and system-state MQTT topics.
* Support for `cool`, `heat`, `dry` and `fan_only` global HVAC modes.
* Updated test fixtures for the Koolnova 2.0 register map.
* Regression tests to ensure zone commands cannot accidentally change the
  global HVAC mode.

### Important safety note

Registers **40080 and 40081**, related to APZ opening-angle configuration,
are currently treated as unsupported for writing.

Their exact behaviour has not been sufficiently verified on real hardware,
so this fork deliberately does **not** write to them.


**koolnova2mqtt** is a bidirectional bridge between MQTT and Koolnova's [100-CPND00](https://koolnova.com/catalogo/novaplus-domotica-100-cpnd00/) family of controllers, featuring auto-discovery for [Home Assistant](https://www.home-assistant.io), and compatible with [OpenHAB](https://www.openhab.org/) and any other automation platform that supports MQTT.

## Features

* Bi-directional synchronization between MQTT topics and Koolnova thermostats.
* Home Assistant auto-discovery as `climate` component thermostats and `sensor` reporting current temperature
* Written in go, cross-platform.

When connected, **koolnova2mqtt** reads all configuration parameters and exports to your MQTT server a topic structure and configuration parameters for Home Assistant. This allows to use thermostats and temperature sensor cards like this one:
 ![ha thermostat](gr/ha-thermostat.png)

Koolnova thermostat modes can be changed by entering the thermostat settings menu (three dots in the top right) and selecting various modes:

![ha thermostat settings](gr/ha-thermostat-settings.png)

# Getting Started

## Requirements

* A Raspberry Pi or a PC 
* A RS4835 USB dongle, such as [this one](https://es.aliexpress.com/item/32978270588.html):![dongle](gr/dongle.png)

## Connecting

Wire the controller's D+ and D- ports as follows:

![Schematic](gr/schematic.png)

* Controller D+ to USB dongle A
* Controller D- to USB dongle B
* Controller GND to USB dongle GND

## Download

Check the [Releases](https://github.com/vayavalla/koolnova2mqtt/releases) page and download the appropriate binary for your platform.

## Running koolnova2mqtt

Before using `koolnova2mqtt`, test your connection with a tool such as [modpoll](https://www.modbusdriver.com/modpoll.html) to ensure your PC and controller can communicate. For example, the following command should return the controller's slave id (default 49):

```bash
./modpoll -b 9600 -p even -d 8 -s 1 -m rtu -t 4 -a 49 -1 -r 106 -4 100 /dev/ttyUSB0
```

Where:

* `-a 49` means send the message to slave with id `49`
* `-r 106` means read register #40106, which contains the slave ID on Koolnova 2.0
* `/dev/ttyUSB0` is the USB dongle's serial port. In Windows it would look like `COM3` or similar. Look in `/dev/ttyUSB*` or in Windows Device Manager to locate the port for your dongle.

Check [modpoll](https://www.modbusdriver.com/modpoll.html) documentation for further information.

## Command line reference
Once you are certain your dongle is properly connected, you can launch `koolnova2mqtt` with these parameters:

```
SYNTAX: 

koolnova2mqtt [options]

options:
  --clientid string
    	A clientid for the connection (default "your hostname")
  --hassPrefix string
    	Home assistant discovery prefix (default "homeassistant")
  --modbusDataBits int
    	Modbus port data bits (default 8)
  --modbusParity string
    	N - None, E - Even, O - Odd (default E) (The use of no parity requires 2 stop bits.) (default "E")
  --modbusPort string
    	Serial port where modbus hardware is connected (default "/dev/ttyUSB0")
  --modbusRate int
    	Modbus port data rate (default 9600)
  --modbusSlaveIDs string
    	Comma-separated list of modbus slave IDs to manage (default "49")
  --modbusSlaveNames string
    	Comma-separated list of modbus slave names. Defaults to 'slave#'
  --modbusStopBits int
    	Modbus port stop bits (default 1)
  --password string
    	Password to match MQTT username
  --prefix string
    	MQTT topic root where to publish/read topics (default "koolnova2mqtt")
  --server string
    	The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883 (default "tcp://127.0.0.1:1883")
  --username string
    	A username to authenticate to the MQTT server
```

### Example:

```
koolnova2mqtt --server tcp://192.168.1.1:1883 --modbusPort '/dev/ttyUSB1' --modbusSlaveIDs '49,50' --modbusSlaveNames 'firstFloor,secondFloor'
```

## MQTT topic structure

The generated structure in MQTT looks as follows:

```
koolnova2mqtt
└── firstFloor
    ├── zone1
    │   ├── enabled = true
    │   ├── enabled
    │   │   └── set
    │   ├── fanMode = auto
    │   ├── fanMode
    │   │   └── set
    │   ├── targetTemp = 20.5
    │   ├── targetTemp
    │   │   └── set
    │   ├── currentTemp = 21
    │   ├── hvacMode = heat
    │   └── hvacMode
    │       └── set        # legacy compatibility
    │
    ├── zone2
    │   ├── enabled = true
    │   ├── fanMode = low
    │   ├── targetTemp = 21
    │   ├── currentTemp = 20
    │   └── hvacMode = heat
    │
    ├── ...
    │
    └── sys
        ├── enabled = true
        │   └── set
        ├── hvacMode = heat
        │   └── set
        ├── efficiency = 3
        ├── holdMode = underfloor and fan
        │   └── set
        ├── serialBaud = 9600
        ├── serialParity = even
        ├── slaveId = 49
        │
        ├── activeModes
        │   ├── fan
        │   ├── cool
        │   ├── heat
        │   └── dry
        │
        ├── temperatureLimits
        │   ├── maxHeat
        │   └── minCool
        │
        ├── autoChange
        │   └── ...
        │
        ├── humidityControl
        │   └── threshold
        │
        ├── waterTemperature
        │   └── ...
        ├── outdoorTemperature
        │   └── ...
        ├── auxTemperature
        │   └── ...
        │
        ├── demand
        │   ├── floor
        │   └── ac3
        │
        ├── ac1
        │   ├── airflow
        │   ├── targetTemp
        │   ├── fanMode
        │   ├── connectedVolume
        │   ├── demandVolume
        │   └── averageTargetTemp
        │
        ├── ac2
        │   ├── airflow
        │   ├── targetTemp
        │   ├── fanMode
        │   ├── connectedVolume
        │   ├── demandVolume
        │   └── averageTargetTemp
        │
        ├── ac3
        │   ├── airflow
        │   ├── targetTemp
        │   ├── fanMode
        │   ├── connectedVolume
        │   ├── demandVolume
        │   ├── averageTargetTemp
        │   ├── efficiency
        │   └── speed
        │
        └── ac4
            ├── airflow
            ├── targetTemp
            ├── fanMode
            ├── connectedVolume
            ├── demandVolume
            └── averageTargetTemp
```
Most of the topics have a child `set` topic that allow you to modify that value. Thus, to change the target temperature of zone2 to 20.5ºC, write the string `20.5` to `koolnova2mqtt/firstFloor/zone2/targetTemp/set` topic. With the tool `mosquitto_pub`:

```bash
mosquitto_pub -t "koolnova2mqtt/firstFloor/zone2/targetTemp/set" -m "20.5"
```

If the operation is successful, the topic `"koolnova2mqtt/firstFloor/zone2/targetTemp"` (without `set`) will be updated with the new target temperature, and the thermostat will show the new value.


## Koolnova 2.0 control topics

In addition to the original topics, this fork provides explicit global and
per-zone control topics.

### Global system control

Global ON/OFF:

```text
koolnova2mqtt/<module>/sys/enabled
koolnova2mqtt/<module>/sys/enabled/set
```

Accepted values for `sys/enabled/set` include:

```text
true
false
ON
OFF
1
0
```

The command writes only the Koolnova 2.0 global system-state register
**40109**.

Global HVAC mode:

```text
koolnova2mqtt/<module>/sys/hvacMode
koolnova2mqtt/<module>/sys/hvacMode/set
```

Supported commands are:

```text
fan_only
cool
heat
dry
```

They correspond to Koolnova 2.0 register **40110**.

`sys/hvacMode/set` does not switch the system ON or OFF. To shut the complete
system down, use:

```text
sys/enabled/set = false
```

This separation is intentional: a global power command writes only 40109,
while a global HVAC-mode command writes only 40110.

### Per-zone control

Each zone publishes its individual enabled state:

```text
koolnova2mqtt/<module>/zone1/enabled
koolnova2mqtt/<module>/zone1/enabled/set
```

and equivalently for the remaining zones.

For example:

```bash
mosquitto_pub \
  -t "koolnova2mqtt/firstFloor/zone1/enabled/set" \
  -m "false"
```

switches only zone 1 off.

The legacy topic:

```text
zoneN/hvacMode/set
```

is retained for compatibility, but in this fork it **never changes the global
HVAC mode**.

* `off` switches only that zone off.
* `cool`, `heat`, `dry` or `fan_only` switch only that zone on.
* The actual global mode must be changed through `sys/hvacMode/set`.

This prevents a command sent to one thermostat from unexpectedly changing the
operating mode of the entire installation.

### Additional Koolnova 2.0 state topics

Depending on the installed equipment, the bridge also publishes information
such as:

```text
sys/activeModes/...
sys/temperatureLimits/...
sys/autoChange/...
sys/humidityControl/threshold
sys/waterTemperature/...
sys/outdoorTemperature/...
sys/auxTemperature/...
sys/demand/...
sys/ac1/connectedVolume
sys/ac1/demandVolume
sys/ac1/averageTargetTemp
...
sys/ac3/efficiency
sys/ac3/speed
```

The same structure is provided for AC1 through AC4 where applicable.

### Hardware-verified 40121-40125 mapping

On the Koolnova 2.0 hardware used to develop and test this fork, the following
mapping was observed:

```text
40121  average target temperature AC1
40122  average target temperature AC2
40123  average target temperature AC3
40124  average target temperature AC4
40125  EFI in the high byte / AC3 speed in the low byte
```

This differs from some interpretations of the available Koolnova 2.0
documentation, so this mapping should be considered **hardware/firmware
verified for the tested installation**, rather than assumed to be universal.

## Author(s)

The original project was written by Javier Peletier ([@jpeletier](https://github.com/jpeletier)). This Koolnova 2.0 fork is adapted and maintained by [@vayavalla](https://github.com/vayavalla).
