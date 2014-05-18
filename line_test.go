package thermal

import (
	"testing"
)

func TestGenerateLineId(t *testing.T) {
	id := generateLineId()
	t.Logf("Line ID: %s", id)
	if len(id) != 16 {
		t.Log("Line ID is not 16 char!")
		t.Logf("Line ID: %s", id)
		t.Fail()
	}
}
