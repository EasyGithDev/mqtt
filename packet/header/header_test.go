package header

import (
	"testing"
)

func TestRemaingLengthEncode(t *testing.T) {

	mp := NewMqttHeader()
	remainingLength := mp.RemainingLengthEncode(321)

	if remainingLength[0] != 193 {
		t.Errorf("remainingLength error  found %d; want 193", remainingLength[0])
	}

	if remainingLength[1] != 2 {
		t.Errorf("remainingLength error  found %d; want 2", remainingLength[1])
	}
}

func TestRemaingLengthDecode(t *testing.T) {

	remainingLength := []byte{193, 2}

	mp := NewMqttHeader()
	res := mp.RemaingLengthDecode(remainingLength)

	if res != 321 {
		t.Errorf("remainingLength error  found %d; want 321", res)
	}

}
