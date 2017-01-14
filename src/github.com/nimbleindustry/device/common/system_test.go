package common

import "testing"

func TestDiskConsumed(t *testing.T) {
	consumed := SystemDiskConsumed()
	if consumed == 0.0 {
		t.Fail()
	}
}

func TestMemoryConsumed(t *testing.T) {
	consumed := SystemMemoryConsumed()
	if consumed == 0.0 {
		t.Fail()
	}
}

func TestSystemLoad(t *testing.T) {
	load := SystemLoadAverage()
	if load == 0.0 {
		t.Fail()
	}
}
