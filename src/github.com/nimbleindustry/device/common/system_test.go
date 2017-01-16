package common

import "testing"

func TestDiskConsumed(t *testing.T) {
	consumed := systemDiskConsumed()
	if consumed == 0.0 {
		t.Fail()
	}
}

func TestMemoryConsumed(t *testing.T) {
	consumed := systemMemoryConsumed()
	if consumed == 0.0 {
		t.Fail()
	}
}

func TestSystemLoad(t *testing.T) {
	load := systemLoadAverage()
	if load == 0.0 {
		t.Fail()
	}
}
