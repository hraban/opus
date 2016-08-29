// Copyright © 2015, 2016 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import "testing"

func TestUseDTX(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_VOIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []bool{true, false}
	for _, dtx := range vals {
		enc.UseDTX(dtx)
		gotv := enc.GetDTX()
		if gotv != dtx {
			t.Errorf("Error set dtx: expect dtx=%v, got dtx=%v", dtx, gotv)
		}
	}
}
