package tfc

import (
	"fmt"
	"hash/crc32"

	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

const (
	N  = "NORTH"
	NE = "NORTH-EAST"
	SE = "SOUTH-EAST"
	S  = "SOUTH"
	SW = "SOUTH-WEST"
	NW = "NORTH-WEST"
)

func pointHash(c tfcPb.Coord) uint32 {
	return crc32.ChecksumIEEE([]byte(c.String()))
}

func edgeHash(c tfcPb.Coord, o string) uint32 {
	return crc32.ChecksumIEEE([]byte(c.String() + o))
}

func NewGameBoard() *tfcPb.GameBoard {
	c0 := tfcPb.Coord{X: 0, Y: 0}
	o0 := N
	gb := &tfcPb.GameBoard{
		Edges:         make(map[uint32]*tfcPb.Edge),
		Intersections: make(map[uint32]*tfcPb.Intersection),
		Tiles:         make(map[uint32]*tfcPb.Tile),
	}
	generateTile(gb, c0, o0)

	return gb
}

func generateTile(gb *tfcPb.GameBoard, c tfcPb.Coord, o string) error {
	currC := c
	currO := o
	currI, currE := newIEPair(gb, currC, currO)
	tileID := currI.Id
	currE.IncidentTile = tileID

	E0 := currE

	gb.Intersections[currI.Id] = currI
	gb.Edges[currE.Id] = currE

	for k := 0; k < 5; k++ {
		nextC, nextO, err := nextCoord(currC, currO)
		if err != nil {
			return fmt.Errorf("could not generate game board: %v", err)
		}

		nextI, nextE := newIEPair(gb, nextC, nextO)
		nextE.IncidentTile = tileID
		currE.Next, nextE.Prev = nextE.Id, currE.Id

		gb.Intersections[nextI.Id] = nextI
		gb.Edges[nextE.Id] = nextE

		currE, currI, currC, currO = nextE, nextI, nextC, nextO
	}

	currE.Next, E0.Prev = E0.Id, currE.Id

	return nil
}

func newIEPair(gb *tfcPb.GameBoard, c tfcPb.Coord, o string) (
	I *tfcPb.Intersection, E *tfcPb.Edge) {

	iID := pointHash(c)

	I, ok := gb.Intersections[iID]
	if !ok {
		I = &tfcPb.Intersection{
			Attributes: &tfcPb.IntersectionAttributes{
				Settlement: tfcPb.Settlement_NOSETTLE,
			},
			Coordinates: &c,
			Id:          iID,
		}
	}

	E = &tfcPb.Edge{
		Id: edgeHash(c, o),
		Attributes: &tfcPb.EdgeAttributes{
			Road: tfcPb.Road_NOROAD,
		},
		Origin: iID,
	}

	return I, E
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
