package prettyprint

import (
	"strings"
)

type Canvas struct {
	Width       int
	Height      int
	pointDist   int
	prettyRunes [][]rune
}

func NewCanvas(width, height int) *Canvas {
	return NewCanvasWithPD(width, height, 6)
}

func NewCanvasWithPD(width, height, pointDist int) *Canvas {

	width *= pointDist
	height *= pointDist

	templateRow := strings.Repeat(" ", width-1) + "\n"
	prettyRunes := make([][]rune, height)
	for y := 0; y < height; y++ {
		runeRow := []rune(templateRow)
		prettyRunes[y] = runeRow
	}

	return &Canvas{
		Width:       width,
		Height:      height,
		pointDist:   pointDist,
		prettyRunes: prettyRunes,
	}
}

func (c *Canvas) DrawPoint(x, y int) {
	x *= c.pointDist
	y *= c.pointDist
	c.prettyRunes[y][x] = '*'
}

func (c *Canvas) DrawLabel(x, y int, label string) {
	x *= c.pointDist
	y *= c.pointDist
	c.drawLabel(x, y, label)
}

func (c *Canvas) DrawLine(x0, y0, x1, y1 int, label string) {

	distX := (x0 - x1) * c.pointDist
	distY := (y0 - y1) * c.pointDist

	midX := distX / 2
	midY := distY / 2
	x, y := x0*c.pointDist, y0*c.pointDist

	for distX > 1 || distX < -1 ||
		distY > 1 || distY < -1 {

		dirX := direction(distX)
		dirY := direction(distY)

		x += dirX
		y += dirY

		if label != "" &&
			(distX == (midX-dirX) && distY == (midY-dirY)) {
			c.drawLabel(x, y, label)
		} else {
			// log.Printf("Drawing rune %v at <%v, %v> ", mark, x, y)
			// log.Printf("Remaining diff is <%v, %v>", distX, distY)
			mark := getLineMark(dirX, dirY)
			c.prettyRunes[y][x] = mark
		}

		distX = reduceDist(distX)
		distY = reduceDist(distY)
	}
}

func reduceDist(dist int) int {
	if dist < 0 {
		return dist + 1
	}
	if dist > 0 {
		return dist - 1
	}
	return dist
}

func direction(dist int) int {
	if dist < 0 {
		return 1
	}
	if dist > 0 {
		return -1
	}
	return 0
}

func getLineMark(dirX, dirY int) rune {
	switch {
	case dirY == 0:
		return '-'
	case dirX == 0:
		return '|'
	case dirX < 0 && dirY < 0:
		return '\\'
	case dirX < 0 && dirY > 0:
		return '/'
	case dirX > 0 && dirY < 0:
		return '/'
	case dirX > 0 && dirY > 0:
		return '\\'
	default:
		return '-'
	}
}

func (c *Canvas) drawLabel(x, y int, label string) {
	labelStart := x - len(label)/2
	for lx, r := range label {
		c.prettyRunes[y][labelStart+lx] = r
	}
}

func (c *Canvas) String() string {

	tempResult := make([]string, c.Height)
	for y, rRow := range c.prettyRunes {
		tempResult[y] = string(rRow)
	}

	return strings.Join(tempResult, "")
}
