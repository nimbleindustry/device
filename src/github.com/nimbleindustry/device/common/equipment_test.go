package common

import (
	"encoding/json"
	"testing"

	"github.com/goburrow/modbus"
	"github.com/stretchr/testify/assert"
)

func TestLoadModbusDefinitions(t *testing.T) {
	var equipment Equipment
	err := json.Unmarshal([]byte(equipmentFixture1), &equipment)
	assert.Nil(t, err, "unmarshall failed")
	integrations := equipment.MachineIntegrations.ModbusEntries
	assert.Equal(t, 5, len(integrations), "expected number of modbus integrations to equal 5")
	descriptionMap := equipment.Desc
	assert.Equal(t, 2, len(descriptionMap), "expected number of description entries to equal 2")
	jpEntry := descriptionMap["jp"]
	assert.Equal(t, "テストのために軽快な産業で使用される機器", jpEntry, "expected テストのために軽快な産業で使用される機器 from description map")
}

func TestFindModbusReadDefinitions(t *testing.T) {
	var equipment Equipment
	err := json.Unmarshal([]byte(equipmentFixture1), &equipment)
	assert.Nil(t, err, "unmarshall failed")
	integrations := equipment.MachineIntegrations.ModbusEntries
	assert.Equal(t, 5, len(integrations), "expected number of modbus integrations to equal 5")
	entries := equipment.MachineIntegrations.FindModbusEntriesByFunction([]int{2})
	assert.Equal(t, 2, len(entries), "expected number of modbus integrations with function 2 to equal 2")
	entries = equipment.MachineIntegrations.FindModbusEntriesByFunction([]int{modbus.FuncCodeReadInputRegisters, modbus.FuncCodeWriteSingleRegister})
	assert.Equal(t, 2, len(entries), "expected number of modbus integrations with function 4&6 to equal 2")
	entries = equipment.MachineIntegrations.FindModbusEntriesByFunction([]int{modbus.FuncCodeReadCoils})
	assert.Equal(t, 1, len(entries), "expected number of modbus integrations with function 1 to equal 1")
	name := entries[0].RegisterName
	assert.Equal(t, "CommandPump", name, "expected name to equal CommandPump")
	entries = equipment.MachineIntegrations.FindModbusEntriesByClass("state")
	assert.Equal(t, 2, len(entries), "expected number of modbus integrations with class of 'state' to equal 2")

}

const equipmentFixture1 = `{
  "ref": "http://machineconfig.com/nimbleindustry.com/test-equipment-a-1.json",
  "entity": "nimbleindustry.com",
  "uri": "http://nimbleindustry.com/productpage/testequipment",
  "desc": {
    "en": "Equipment used by Nimble Industry for testing",
    "jp": "テストのために軽快な産業で使用される機器"
  },
  "models": [
    "A"
  ],
  "versions": [
    "1"
  ],
  "machineIntegration": {
    "modbus": [
      {
        "registerName": "TankFull",
        "functions": [2],
        "address": 0,
        "class": "state",
        "desc": {
          "en": "Tank full",
          "jp": "満タン"
        }
      },
      {
        "registerName": "TankEmpty",
        "functions": [2],
        "address": 1,
        "class": "state",
        "desc": {
          "en": "Tank empty",
          "jp": "空タンク"
        }
      },
      {
        "registerName": "CommandPump",
        "functions": [1, 5],
        "address": 0,
        "class": "control",
        "desc": {
          "en": "Command pump",
          "jp": "コマンドポンプ"
        }
      },
      {
        "registerName": "SystemRun",
        "functions": [3, 6],
        "address": 0,
        "class": "control",
        "desc": {
          "en": "System run control",
          "jp": "システムメイン制御"
        }
      },
      {
        "registerName": "LiquidTemp",
        "functions": [4],
        "address": 0,
        "class": "telemetry",
        "desc": {
          "en": "Liquid Temp",
          "jp": "タンク液温"
        }
      }
    ]
  },
  "errorCodes": [
    {
      "code": "XXX",
      "desc": {
        "en": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
        "jp": "Loremのイプサムの嘆き、AMET consecteturのadipiscingのELIT座ります。"
      }
    }
  ],
  "maintenance": [
    {
      "id": "XXX",
      "interval": "30d",
      "title": {
        "en": "Clean filters",
        "jp": "クリーンフィルタ"
      },
      "desc": {
        "en": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
        "jp": "Loremのイプサムの嘆き、AMET consecteturのadipiscingのELIT座ります。"
      },
      "uri": "http://youtube.com/watch/XXXXXXXXX"
    }
  ]
}`
