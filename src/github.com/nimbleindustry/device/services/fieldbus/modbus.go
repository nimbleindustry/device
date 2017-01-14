package fieldbus

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/nimbleindustry/device/common"
	"github.com/nimbleindustry/device/define"

	"github.com/goburrow/modbus"
)

// GenericModbusService implements core modbus services that are utilized
// by both TCP and RTU implementations
type GenericModbusService struct {
	connection          common.ConnectionRecord
	client              modbus.Client
	machineIntegrations []common.ModbusEntry
}

func (svc *GenericModbusService) initConfigurations() error {
	if common.ConnectionConfig.DeviceID == "" {
		return errors.New("Connection config object not initialized")
	}
	svc.connection = common.ConnectionConfig.GetMachineConnection(define.ModbusTCP)
	if common.EquipmentConfig.Ref == "" {
		return errors.New("Equipment config object not initialized")
	}
	svc.machineIntegrations = common.EquipmentConfig.MachineIntegrations.ModbusEntries
	return nil
}

func (svc *GenericModbusService) initBusIntegration() error {
	if len(svc.machineIntegrations) == 0 {
		return errors.New("No modbus entries")
	}
	return nil
}

func (svc *GenericModbusService) readAllInputs() (map[string]interface{}, error) {
	m := make(map[string]interface{}, 1)
	err := svc.readAllDiscreteInputs(m)
	if err != nil {
		return nil, err
	}
	err = svc.readAllCoils(m)
	if err != nil {
		return nil, err
	}
	err = svc.readAllHoldingRegisters(m)
	if err != nil {
		return nil, err
	}
	err = svc.readAllInputRegisters(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (svc *GenericModbusService) readAllDiscreteInputs(m map[string]interface{}) error {
	if svc.client == nil {
		return errors.New("client nil")
	}
	entries := common.EquipmentConfig.MachineIntegrations.FindModbusEntriesByFunction([]int{modbus.FuncCodeReadDiscreteInputs})
	for _, v := range entries {
		value, err := svc.client.ReadDiscreteInputs(uint16(v.Address), 1)
		if err != nil {
			return fmt.Errorf("Error reading discrete input %s, %s", v.Desc["en"], err)
		}
		m[v.RegisterName] = value[0]
	}
	return nil
}

func (svc *GenericModbusService) readAllHoldingRegisters(m map[string]interface{}) error {
	if svc.client == nil {
		return errors.New("client nil")
	}
	entries := common.EquipmentConfig.MachineIntegrations.FindModbusEntriesByFunction([]int{modbus.FuncCodeReadHoldingRegisters})
	for _, v := range entries {
		value, err := svc.client.ReadHoldingRegisters(uint16(v.Address), 1)
		if err != nil {
			return fmt.Errorf("Error reading holding register %s, %s", v.Desc["en"], err)
		}
		m[v.RegisterName] = BytesToInt16(value)
	}
	return nil
}

func (svc *GenericModbusService) readAllCoils(m map[string]interface{}) error {
	if svc.client == nil {
		return errors.New("client nil")
	}
	entries := common.EquipmentConfig.MachineIntegrations.FindModbusEntriesByFunction([]int{modbus.FuncCodeReadCoils})
	for _, v := range entries {
		value, err := svc.client.ReadCoils(uint16(v.Address), 1)
		if err != nil {
			return fmt.Errorf("Error reading coils %s, %s", v.Desc["en"], err)
		}
		m[v.RegisterName] = value[0]
	}
	return nil
}

func (svc *GenericModbusService) readAllInputRegisters(m map[string]interface{}) error {
	if svc.client == nil {
		return errors.New("client nil")
	}
	entries := common.EquipmentConfig.MachineIntegrations.FindModbusEntriesByFunction([]int{modbus.FuncCodeReadInputRegisters})
	for _, v := range entries {
		value, err := svc.client.ReadInputRegisters(uint16(v.Address), 1)
		if err != nil {
			return fmt.Errorf("Error reading input register %s, %s", v.Desc["en"], err)
		}
		m[v.RegisterName] = BytesToInt16(value)
	}
	return nil
}

// BytesToUint16 converts the passed byte slice to an unsigned 16bit integer
func BytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// Uint16ToBytes converts a unsigned 16bit integer to a byte slice
func Uint16ToBytes(u uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, u)
	return buf
}

// BytesToInt16 converts the passed byte slice to a signed 16bit integer
func BytesToInt16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

// Int16ToBytes converts a unsigned 16bit integer to a byte slice
func Int16ToBytes(u int16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(u))
	return buf
}
