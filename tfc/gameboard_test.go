package tfc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

func TestGenerateTile(t *testing.T) {
	gb, err := NewGameBoard()
	require.NoError(t, err)

	for _, I := range gb.Intersections {
		require.NotNil(t, I.Id,
			"expected intersection ID for %v", I)
		require.NotNil(t, I.Attributes,
			"expected intersection attributes for %v", I)
		require.NotNil(t, I.Coordinates,
			"expected intersection coordinates for %v", I)
		require.NotNil(t, I.IncidentEdge,
			"expected intersection incident edge for %v", I)
	}

	for _, E := range gb.Edges {
		require.NotNil(t, E.Id,
			"expected edge ID for %v", E)
		require.NotNil(t, E.Attributes,
			"expected edge coordinates for %v", E)
		require.NotNil(t, E.Origin,
			"expected edge origin for %v", E)
		require.NotNil(t, E.Next,
			"expected edge next pointer for %v", E)
		require.NotNil(t, E.Prev,
			"expected edge prev pointer for %v", E)
		require.NotNil(t, E.IncidentTile,
			"expected edge tile for %v", E)

		// if E.GetTwin() != 0 {
		// 	twin := gb.Edges[E.GetTwin()]
		// 	require.Equal(t, E.Id, twin.GetTwin(),
		// 		"expected twin edge to point back for %v:\n\t got %v", E, twin)
		// }
	}

	fmt.Println(prettyPrintGb(*gb))

}

func prettyPrintGb(gb tfcPb.GameBoard) string {
	canvas := NewCanvas(16, 12)

	var xOffset, yOffset int32 = 8, 7

	for _, I := range gb.Intersections {
		x := int(I.Coordinates.X + xOffset)
		y := int(I.Coordinates.Y + yOffset)
		// log.Printf("Printing intersection %v at position <%d, %d>", I, x, y)
		canvas.DrawPoint(x, y)
	}

	for _, E := range gb.Edges {
		origin := gb.Intersections[E.Origin]
		dest := gb.Intersections[gb.Edges[E.Next].Origin]
		x0 := int(origin.Coordinates.X + xOffset)
		y0 := int(origin.Coordinates.Y + yOffset)
		x1 := int(dest.Coordinates.X + xOffset)
		y1 := int(dest.Coordinates.Y + yOffset)

		// log.Printf("Drawing edge (%v, %v) - (%v, %v)", x0, y0, x1, y1)

		canvas.DrawLine(x0, y0, x1, y1)
	}

	return canvas.String()
}
