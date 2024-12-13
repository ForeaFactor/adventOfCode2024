package day_06

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func Main() {
	input := readInput("./day_06/sample.txt")
	lab := newGridFromText(input)

	fmt.Printf("\n====== DAY 06 ======\n")
	fmt.Printf("%d = number of positions the guard did not pass\n", 0)
	fmt.Printf(string(lab.exportGridToText()))

}

//---------structs declaration---------

type grid struct {
	pois   map[cord]poi //PointOfInterests
	height int
	width  int
	// poiIsAllowed(p poi) bool

}

type cord struct {
	x int
	y int
}

type poi struct {
	icon      byte
	id        int // unused but reasonable
	isMovable bool
	pos       cord
	facing    cord
}

//---------methods declaration---------

func (g *grid) getPoiByPos(pos cord) (*poi, error) {
	if !g.containsCord(pos) {
		return new(poi), errors.New("no such position in grid")
	}
	// check if key in map (poicord) exists
	point, ok := g.pois[pos]
	if !ok {
		// if empty -> i thought of an empty space poi
		point = *new(poi)
		point.icon = byte('.')
		point.pos = pos
	}
	return &point, nil
}

func (g *grid) exportGridToText() []byte {
	txt := make([]byte, (g.width+1)*g.height)
	var txtPos int = 0
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			gridPoint, err := g.getPoiByPos(cord{x, y})
			if err != nil {
				log.Fatal(err)
			}
			icon := gridPoint.icon
			txt[txtPos] = icon
			if x == g.width-1 {
				txtPos++
				txt[txtPos] = byte('\n')
			}
			txtPos++
		}
	}
	return txt
}

func (g *grid) containsCord(c cord) bool {
	// method to check if coordinates are in bound of the grid
	if isIntBetween(0, g.width, c.x) && isIntBetween(0, g.height, c.y) {
		// including lower and excluding upper
		return true
	}
	return false
}

func (p *poi) setPos(c cord) {
	// TODO: this method seems unnecessary
	p.facing = c
}

func (p *poi) rotate(dir string) error {
	var directions = [4]cord{
		{0, -1}, {1, 0}, {0, 1}, {-1, 0},
	}
	var facing = -1 // dummy value
	for dirIdx, comparedir := range directions {
		if p.facing == comparedir {
			facing = dirIdx
			break
		}
	}
	if facing == -1 {
		return errors.New("poi has invalid facing")
	}
	var err error

	switch dir {
	case "right":
		if facing == 3 {
			p.facing = directions[0]
		} else {
			p.facing = directions[facing+1]
		}
		err = nil
	case "left":
		if facing == 0 {
			p.facing = directions[3]
		} else {
			p.facing = directions[facing-1]
		}
		err = nil
	default:
		err = errors.New("no such direction to turn to")
	}
	if err != nil {
		// TODO Add logic to assign new Icon to rotated guards
	}
	return err
}

func newGridFromText(input []byte) grid {
	// this is a constructor for grid
	// assume ascii Encoding // assume EACH Line ends with a LF
	var out grid
	out.pois = make(map[cord]poi)
	in := strings.Split(string(input), "\n")
	out.width = len(in[0])   // TODO: add check if each line has the same length
	out.height = len(in) - 1 // account for linefeed at the last line - Split generates extra Slice for that
	for yCord, line := range in {
		for xCord, icon := range line {
			pos := cord{x: xCord, y: yCord} // no need for pos consistency check
			if point, err := newPoiFromIcon(byte(icon), pos); err != nil {
				log.Fatal(err)
			} else {
				out.pois[pos] = point // yea should use a setter function
			}
		}
	}
	return out
}

func newPoiFromIcon(b byte, pos cord) (poi, error) {
	// constructor for poi
	// TODO: add implement icon check, when the facing of a poi changes -> func poi.rotate( pos )
	// TODO: do something with the error
	var point = poi{
		icon:      b,
		pos:       pos,
		facing:    cord{}, // if all pois face to {0,0} by standart
		isMovable: false,
		id:        0, // TODO: default not used
	}
	switch b {
	case byte('^'):
		point.isMovable = true
		point.facing = cord{0, -1}
	case byte('<'):
		point.isMovable = true
		point.facing = cord{-1, 0}
	case byte('>'):
		point.isMovable = true
		point.facing = cord{1, 0}
	case byte('v'):
		point.isMovable = true
		point.facing = cord{1, 0}
	}
	return point, nil
}

//---------functions declaration---------

func isIntBetween(lower int, upper int, x int) bool {
	// inlcuding lower and excluding upper
	return x >= lower && x < upper
}

func readInput(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
