package tfc

import (
	"fmt"
	"hash/crc32"
	"math/rand"

	"github.com/golang-collections/collections/stack"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

const (
	N  = "NORTH"
	NE = "NORTH-EAST"
	SE = "SOUTH-EAST"
	S  = "SOUTH"
	SW = "SOUTH-WEST"
	NW = "NORTH-WEST"

	SIZE = 2
)

func pointHash(c tfcPb.Coord) uint32 {
	return crc32.ChecksumIEEE([]byte(c.String()))
}

func edgeHash(c tfcPb.Coord, o string) uint32 {
	return crc32.ChecksumIEEE([]byte(c.String() + o))
}

type tileAttributeStacks struct {
	resourceStack *stack.Stack
	rollStack     *stack.Stack
}

func NewGameBoard() (*tfcPb.GameBoard, error) {
	c0 := tfcPb.Coord{X: 0, Y: 0}
	o0 := N
	gb := &tfcPb.GameBoard{
		Edges:         make(map[uint32]*tfcPb.Edge),
		Intersections: make(map[uint32]*tfcPb.Intersection),
		Tiles:         make(map[uint32]*tfcPb.Tile),
	}

	resourceStack := newResourceStack()
	rollStack := newRollStack()
	tileAttrStacks := tileAttributeStacks{resourceStack, rollStack}

	err := generateTile(gb, c0, o0, tileAttrStacks)
	if err != nil {
		return nil, fmt.Errorf("Could not generate initial tile, %s", err)
	}
	for l := 0; l < SIZE; l++ {
		err := expandGameBoard(gb, tileAttrStacks)
		if err != nil {
			return nil, fmt.Errorf("could not expand gb: %s", err)
		}
	}

	return gb, nil
}

func newResourceStack() *stack.Stack {
	resStack := stack.New()
	for i := 0; i < 10; i++ {
		resStack.Push(tfcPb.Resource_CAMP)
		resStack.Push(tfcPb.Resource_FIELD)
		resStack.Push(tfcPb.Resource_FOREST)
		resStack.Push(tfcPb.Resource_HILL)
		resStack.Push(tfcPb.Resource_MOUNTAIN)
		resStack.Push(tfcPb.Resource_PASTURE)
	}
	return resStack
}

func newRollStack() *stack.Stack {
	rollStack := stack.New()
	for i := 0; i < 5; i++ {
		for _, rn := range rand.Perm(10) {
			rollStack.Push(int32(rn + 2))
		}
	}
	return rollStack
}

func expandGameBoard(gb *tfcPb.GameBoard, tileAttrStacks tileAttributeStacks) error {
	edgeIDs := []uint32{}
	for eID, E := range gb.Edges {
		if E.Twin == 0 {
			edgeIDs = append(edgeIDs, eID)
		}
	}

	// log.Printf("####\n\n Expanding gb \n\n")

	for _, eID := range edgeIDs {
		E := gb.Edges[eID]

		// The twin might have been set on a previous expand step
		// Continute to next in this case.
		if E.Twin != 0 {
			continue
		}

		originId := gb.Edges[E.Next].Origin
		c := *gb.Intersections[originId].Coordinates
		o := twinOrientation(E.Orientation)
		err := generateTile(gb, c, o, tileAttrStacks)
		if err != nil {
			return fmt.Errorf("could not expand on %v: %s", E, err)
		}

		// log.Printf("Expanded on edge %s", E)
	}
	return nil
}

func generateTile(gb *tfcPb.GameBoard, c tfcPb.Coord, o string, attrStacks tileAttributeStacks) error {
	currC := c
	currO := o
	currI, currE := newIEPair(gb, currC, currO)
	tileID := currE.Id

	gb.Intersections[currI.Id] = currI
	currE.IncidentTile = tileID
	// log.Printf("Setting incident tile for %v (%v,%v) to %v", currE.Id, currC, currO, tileID)
	gb.Edges[currE.Id] = currE

	gb.Tiles[tileID] = &tfcPb.Tile{
		Id:             tileID,
		OuterComponent: currE.Id,
		Attributes: &tfcPb.TileAttributes{
			Resource:   attrStacks.resourceStack.Pop().(tfcPb.Resource),
			RollNumber: attrStacks.rollStack.Pop().(int32),
		},
	}

	E0 := currE
	for k := 0; k < 5; k++ {
		nextC, nextO, err := nextCoord(currC, currO)
		if err != nil {
			return fmt.Errorf("could not generate game board: %v", err)
		}

		nextI, nextE := newIEPair(gb, nextC, nextO)
		nextE.IncidentTile = tileID
		// log.Printf("Setting incident tile for %v (%v,%v) to %v", nextE.Id, nextC, nextO, tileID)
		currE.Next, nextE.Prev = nextE.Id, currE.Id

		gb.Intersections[nextI.Id] = nextI
		gb.Edges[nextE.Id] = nextE

		currE, currI, currC, currO = nextE, nextI, nextC, nextO
	}

	currE.Next, E0.Prev = E0.Id, currE.Id

	updateTwins(gb, E0)

	return nil
}

func newIEPair(gb *tfcPb.GameBoard, c tfcPb.Coord, o string) (
	I *tfcPb.Intersection, E *tfcPb.Edge) {

	iID := pointHash(c)
	eID := edgeHash(c, o)
	I, ok := gb.Intersections[iID]
	if !ok {
		I = &tfcPb.Intersection{
			Attributes: &tfcPb.IntersectionAttributes{
				Settlement: tfcPb.Settlement_NOSETTLE,
			},
			Coordinates:  &c,
			Id:           iID,
			IncidentEdge: eID,
		}
	}
	E = &tfcPb.Edge{
		Id: eID,
		Attributes: &tfcPb.EdgeAttributes{
			Road: tfcPb.Road_NOROAD,
		},
		Origin:      iID,
		Orientation: o,
	}

	return I, E
}

func updateTwins(gb *tfcPb.GameBoard, e0 *tfcPb.Edge) {
	nextE := gb.Edges[e0.Next]
	currE := e0

	for {
		twinC := *gb.Intersections[nextE.Origin].Coordinates
		twinO := twinOrientation(currE.Orientation)
		twinID := edgeHash(twinC, twinO)
		if twin, ok := gb.Edges[twinID]; ok {
			currE.Twin = twinID
			twin.Twin = currE.Id
		}
		currE = nextE
		nextE = gb.Edges[nextE.Next]

		if currE.Id == e0.Id {
			break
		}
	}

}

func twinOrientation(o string) string {
	switch o {
	case N:
		return S
	case NW:
		return SE
	case SW:
		return NE
	case S:
		return N
	case SE:
		return NW
	case NE:
		return SW
	default:
		return o
	}
}

func nextCoord(c tfcPb.Coord, o string) (tfcPb.Coord, string, error) {
	switch o {
	case N:
		return tfcPb.Coord{X: c.X - 1, Y: c.Y - 1}, NW, nil
	case NW:
		return tfcPb.Coord{X: c.X, Y: c.Y - 1}, SW, nil
	case SW:
		return tfcPb.Coord{X: c.X + 1, Y: c.Y - 1}, S, nil
	case S:
		return tfcPb.Coord{X: c.X + 1, Y: c.Y + 1}, SE, nil
	case SE:
		return tfcPb.Coord{X: c.X, Y: c.Y + 1}, NE, nil
	case NE:
		return tfcPb.Coord{X: c.X - 1, Y: c.Y + 1}, N, nil
	default:
		return c, o, fmt.Errorf("unkown orientation")
	}
}
