package util

import "testing"

func TestCompareNumInStr(t *testing.T) {
	var prev, next string

	prev = "2"
	next = "3"
	if CompareNumInStr(prev, next) != true {
		t.Error(prev, next)
	}

	prev = "2"
	next = "20"
	if CompareNumInStr(prev, next) != true {
		t.Error(prev, next)
	}

	prev = "1.1"
	next = "1.10"
	if CompareNumInStr(prev, next) != true {
		t.Error(prev, next)
	}

	prev = "1.2"
	next = "1.10"
	if CompareNumInStr(prev, next) != true {
		t.Error(prev, next)
	}

	prev = "1.3.2"
	next = "1.3.10"
	if CompareNumInStr(prev, next) != true {
		t.Error(prev, next)
	}
}