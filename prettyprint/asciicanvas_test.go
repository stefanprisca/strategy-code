package prettyprint

import (
	"fmt"
	"testing"
)

func TestCanvas(t *testing.T) {
	canvas := NewCanvas(4, 4)
	canvas.DrawLine(0, 0, 1, 1, "foo")
	canvas.DrawLine(2, 2, 1, 1, "t")

	canvas.DrawLine(1, 1, 0, 2, "bar")
	canvas.DrawLine(2, 0, 1, 1, "")
	fmt.Print(canvas)
}
