package common

// Asset uniquely identifies a machine/piece of equipment
type Asset struct {
	MachineID  string `json:"machineId"`
	Serial     string `json:"serial,omitempty"`
	Type       string `json:"type,omitempty"`
	Entity     string `json:"entity,omitempty"`
	Location   string `json:"location,omitempty"`
	Group      string `json:"group,omitempty"`
	Line       string `json:"line,omitempty"`
	WorkCenter string `json:"workCenter,omitempty"`
}
