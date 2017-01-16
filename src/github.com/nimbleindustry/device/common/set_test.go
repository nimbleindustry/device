package common

import (
	"encoding/json"
	"testing"
)

func compare(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("value differs. Expected [%v], actual [%v]", expected, actual)
	}
}

func TestSetAddRemoveContains(t *testing.T) {
	set := NewSet()
	set.Add(10)
	compare(t, true, set.Has(10))
	set.Remove(10)
	compare(t, false, set.Has(10))
	compare(t, false, set.Has("shie"))
}

func TestSetMap(t *testing.T) {
	set := NewSet()
	set.Add(10)
	set.Add(12)
	set.Add(14)
	compare(t, true, set.Has(10))
	compare(t, false, set.Has(11))
	set = set.Map(func(x interface{}) interface{} { return x.(int) + 1 })
	compare(t, false, set.Has(10))
	compare(t, true, set.Has(11))
}

func TestSetMarshalAndUnMarshal(t *testing.T) {
	set := NewSet()
	set.Add("hello")
	set.Add("there")
	set.Add("you")
	compare(t, true, set.Has("you"))
	marshaled, err := json.Marshal(set)
	if err != nil {
		t.Fail()
	}
	unmarshaled := NewSet()
	err = json.Unmarshal(marshaled, unmarshaled)
	if err != nil {
		t.Fail()
	}
	compare(t, true, unmarshaled.Has("you"))
}
