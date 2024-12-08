package day_08

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func Main() {
	input := readInput("./day_08/input.txt")
	locationsMap := genAntennasGridFromText(input)

	fmt.Printf("\n====== DAY 04 ======\n")
	fmt.Printf("%d = Number of Antinodes\n", countAntinodes(locationsMap))

	print(string(locationsMap.exportGridToText()))
}

//---------structs declaration---------

type poi interface {
	getIcon() byte // uses char as type specifier
	getPos() cord
}

type grid struct {
	pois   map[cord][]poi //PointOfInterests
	height int
	width  int
}

type cord struct {
	x int
	y int
}

type poiAntenna struct {
	icon byte
	pos  cord
	freq byte // freqencies identified by the icon for now (love redundancy :))
}

type poiAntinode struct {
	icon    byte
	pos     cord
	sources [2]poiAntenna
}

type vector struct {
	xShift int
	yShift int
}

func (g *grid) exportGridToText() []byte {
	txt := make([]byte, (g.width+1)*g.height)
	var txtPos int
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			gridPoint := *g.getPoiByPos(cord{x, y})
			icon := byte(' ')
			if gridPoint != nil {
				icon = gridPoint.getIcon()
			}
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

//---------methods declaration---------

func (g *grid) getPoiByPos(c cord) *poi {
	// only the first will  be shown - to lazy
	point := poi(nil)
	if g.pois[c] != nil {
		point = g.pois[c][0]
	}
	return &point
}

func (a poiAntenna) getPos() cord {
	return a.pos
}

func (a poiAntenna) getIcon() byte {
	return a.icon
}

func (a poiAntinode) getPos() cord {
	return a.pos
}

func (a poiAntinode) getIcon() byte {
	return a.icon
}

// poiAntenna Constructor
func newPoiAntenna(freq byte, pos cord) *poiAntenna {
	r := regexp.MustCompile("[0-9A-Za-z_]")
	if r.Match([]byte{freq}) == false {
		return nil // TODO: add meaningful error handling
	}
	return &poiAntenna{
		icon: freq,
		pos: cord{
			x: pos.x,
			y: pos.y,
		},
		freq: freq,
	}

}

func readInput(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//---------functions declaration---------

func genAntennasGridFromText(input []byte) grid {
	// assume ascii Encoding
	var out grid
	out.pois = make(map[cord][]poi)
	in := strings.Split(string(input), "\n")
	out.width = len(in[0]) // TODO: add check if each line has the same length
	out.height = len(in)

	for yCord, line := range in {
		for xCord, icon := range line {
			pos := cord{x: xCord, y: yCord}
			antenna_p := newPoiAntenna(byte(icon), pos)
			if antenna_p != nil {
				out.pois[pos] = append(out.pois[pos], *antenna_p)
			}
		}
	}

	return out
}

func countAntinodes(g grid) uint {
	// TODO: implement
	return 0
}
