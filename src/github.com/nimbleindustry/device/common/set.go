package common

import "encoding/json"

// Set conveniently typedefs map[interface{}]bool
type Set map[interface{}]bool

// NewSet constructs a new Set
func NewSet() *Set {
	set := make(Set)
	return &set
}

// Add the supplied argument to the set
func (set *Set) Add(x interface{}) {
	(*set)[x] = true
}

// Remove the supplied argument to the set
func (set *Set) Remove(x interface{}) {
	delete((*set), x)
}

// Reset the set to empty
func (set *Set) Reset() {
	(*set) = make(map[interface{}]bool)
}

// Has tests the set for the supplied argument
func (set *Set) Has(x interface{}) bool {
	_, found := (*set)[x]
	return found
}

// Do runs the supplied function across all members of set
func (set *Set) Do(f func(interface{})) {
	for k := range *set {
		f(k)
	}
}

// Len returns the number of elements in the Set
func (set *Set) Len() int {
	return len(*set)
}

// Map applies the results of the supplied function, using the set's members as arguments,
// and returns those results in a new Set.
func (set *Set) Map(f func(interface{}) interface{}) *Set {
	newSet := NewSet()
	for k := range *set {
		newSet.Add(f(k))
	}
	return newSet
}

// MarshalJSON returns a JSON encoded representation of the Set
func (set *Set) MarshalJSON() ([]byte, error) {
	var newSet []interface{}
	for k := range *set {
		newSet = append(newSet, k)
	}
	return json.Marshal(newSet)
}

// UnmarshalJSON adds passed JSON-encoded body to the Set.
func (set *Set) UnmarshalJSON(body []byte) (err error) {
	var newSet []interface{}
	err = json.Unmarshal(body, &newSet)
	if err != nil {
		return
	}
	for _, v := range newSet {
		set.Add(v)
	}
	return
}
