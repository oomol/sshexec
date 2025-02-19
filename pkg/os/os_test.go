package os

import "testing"

func TestOSInfo(t *testing.T) {
	if IsSequoia() {
		t.Log("This is a Sequoia machine.")
	} else {
		t.Log("This is not a Sequoia machine.")
	}

	if IsVentura() {
		t.Log("This is a Ventura machine.")
	} else {
		t.Log("This is not a Ventura machine.")
	}

	if IsSonoma() {
		t.Log("This is a Sonoma machine.")
	} else {
		t.Log("This is not a Sonoma machine.")
	}
}
