package tfc

import (
	"fmt"
	"testing"
)

func TestCanvas(t *testing.T) {
	canvas := NewCanvas(4, 4)
	canvas.DrawLine(0, 0, 1, 1)
	canvas.DrawLine(2, 2, 1, 1)

	canvas.DrawLine(1, 1, 0, 2)
	canvas.DrawLine(2, 0, 1, 1)
	fmt.Print(canvas)
}
