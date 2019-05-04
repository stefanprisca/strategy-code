package prettyprint

import (
	"strconv"

	tfc "github.com/stefanprisca/strategy-code/tfc"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

func NewTFCBoardCanvas() *Canvas {
	width, height := 16, 13
	return NewCanvasWithPD(width, height, 10)
}

func (c *Canvas) PrettyPrintTfcBoard(gb tfcPb.GameBoard) *Canvas {

	var xOffset, yOffset int = 8, 8

	// Draw the scale of the map

	for _, I := range gb.Intersections {
		x := int(I.Coordinates.X) + xOffset
		y := int(I.Coordinates.Y) + yOffset
		// log.Printf("Printing intersection %v at position <%d, %d>", I, x, y)

		xLabel := strconv.Itoa(x - xOffset)
		yLabel := strconv.Itoa(y - yOffset)

		c.DrawLabel(x, 0, xLabel)
		c.DrawLabel(1, y, yLabel)
	}

	for _, I := range gb.Intersections {
		x := int(I.Coordinates.X) + xOffset
		y := int(I.Coordinates.Y) + yOffset
		// log.Printf("Printing intersection %v at position <%d, %d>", I, x, y)
		iLabel := I.Attributes.Settlement.String()
		c.DrawLabel(x, y, iLabel)
	}

	for _, E := range gb.Edges {
		origin := gb.Intersections[E.Origin]
		dest := gb.Intersections[gb.Edges[E.Next].Origin]
		x0 := int(origin.Coordinates.X) + xOffset
		y0 := int(origin.Coordinates.Y) + yOffset
		x1 := int(dest.Coordinates.X) + xOffset
		y1 := int(dest.Coordinates.Y) + yOffset

		// log.Printf("Drawing edge (%v, %v) - (%v, %v)", x0, y0, x1, y1)

		eLabel := E.Attributes.Road.String()
		c.DrawLine(x0, y0, x1, y1, eLabel)

		if E.Orientation == tfc.N {
			tile := gb.Tiles[E.IncidentTile]
			resLabel := tile.Attributes.Resource.String()
			c.DrawLabel(x0, y0-2, resLabel)

			rollLabel := strconv.Itoa(int(tile.Attributes.RollNumber))
			c.DrawLabel(x0, y0-1, rollLabel)
		}
	}

	return c
}
