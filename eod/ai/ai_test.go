package ai

import (
	"fmt"
	"testing"
)

func TestAI(t *testing.T) {
	a := NewAI()
	a.AddCombo("1+1", true)
	a.AddCombo("2+1", true)
	//a.AddCombo("2+1", true)
	fmt.Println(a.PredictCombo())
}
