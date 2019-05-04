package tfc

import (
	"fmt"
	"testing"

	"github.com/stefanprisca/strategy-code/prettyprint"
	"github.com/stretchr/testify/require"
)

func TestGenerateGameBoard(t *testing.T) {
	gb, err := NewGameBoard()
	require.NoError(t, err)
	require.NotZero(t, len(gb.Intersections),
		"expected to have intersections initialized")

	require.NotZero(t, len(gb.Edges),
		"expected to have edges initialized")

	require.NotZero(t, len(gb.Tiles),
		"expected to have tiles initialized")

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

		if E.GetTwin() != 0 {
			twin := gb.Edges[E.GetTwin()]
			require.Equal(t, E.GetId(), twin.GetTwin(),
				"expected twin edge to point back for %v:\n\t got %v", E, twin)
		}
	}

	for _, T := range gb.Tiles {
		require.NotNil(t, T.GetAttributes(),
			"expected tile to have attributes")
		require.NotZero(t, T.GetOuterComponent(),
			"expected tile to have outer component")
		outerCompID := T.GetOuterComponent()
		outerComp := gb.Edges[outerCompID]
		require.Equal(t, outerComp.IncidentTile, T.Id,
			"expected outer component to point back to tile for %v:\n\t got %v",
			T, outerComp)
	}

	boardPrettyString := prettyprint.NewTFCBoardCanvas().
		PrettyPrintTfcBoard()
	fmt.Println(boardPrettyString)

}
