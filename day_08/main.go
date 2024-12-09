package day_08

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func Main() {
	input := readInput("./day_08/sample.txt")
	gridTsk1 := generateAntennasGridFromText(input)
	antinodes := generateAntinodesModelTsk1(&gridTsk1)
	gridTsk1.addPoiAntinodes(antinodes)

	gridTsk2 := generateAntennasGridFromText(input)
	antinodes = generateAntinodesModelTsk2(&gridTsk2)
	gridTsk2.addPoiAntinodes(antinodes)

	fmt.Printf("\n====== DAY 04 ======\n")
	fmt.Printf("%d = Number of Antinodes\n", countAntinodes(gridTsk1))
	fmt.Printf("%d = Number of Antinodes corrected Model (Tsk2)\n", countAntinodes(gridTsk2))

	fmt.Printf("%s", gridTsk2.exportGridToText())
}

type poi interface {
	getIcon() byte // uses char as type specifier
	getPos() cord
	isAllowedToExistIn(g *grid) bool
}

type grid struct {
	pois   map[cord][]poi //PointOfInterests
	height int
	width  int
}

//---------structs declaration---------

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
	sources [2]*poi // these should be antennas - but is not enforced
}

type vector struct {
	xShift int
	yShift int
}

func (g *grid) addPois(input []poi) {
	/*	List of Points Of Interrest can contain
		- POIs outside the gridSize
		- POIs at positions, where POI(s) already exist */
	for _, point := range input {
		if point.isAllowedToExistIn(g) {
			g.pois[point.getPos()] = append(g.pois[point.getPos()], point)
		}
	}
}

func (g *grid) addPoiAntinodes(input []poiAntinode) {
	for _, point := range input {
		if point.isAllowedToExistIn(g) {
			g.pois[point.getPos()] = append(g.pois[point.getPos()], point)
		}
	}
}

func (g *grid) exportGridToText() []byte {
	txt := make([]byte, (g.width+1)*g.height)
	var txtPos int
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			gridPoint := *g.getPoiByPos(cord{x, y})
			icon := byte('.')
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

func (g *grid) getPoiByPos(c cord) *poi {
	// only the first will  be shown - to lazy ¯\_(ツ)_/¯
	point := poi(nil)
	if g.pois[c] != nil {
		point = g.pois[c][0]
	}
	return &point
}

func (g *grid) containsCord(c cord) bool {
	// method to check if coordinates are in bound of the grid
	if isIntBetween(0, g.width, c.x) && isIntBetween(0, g.height, c.y) {
		// inlcuding lower and excluding upper
		return true
	}
	return false
}

func (a poiAntenna) getPos() cord {
	return a.pos
}

//---------methods declaration---------

func (a poiAntenna) getIcon() byte {
	return a.icon
}

func (a poiAntenna) isAllowedToExistIn(g *grid) bool {
	// multiple antennas at same Coordinates are allowed
	return g.containsCord(a.pos)
}

func (a poiAntinode) getPos() cord {
	return a.pos
}

func (a poiAntinode) getIcon() byte {
	return a.icon
}

func (a poiAntinode) isAllowedToExistIn(g *grid) bool {
	// [ ] multiple antinodes of the same freqency are not allowed
	// [x] muliple antinodes per position are not allowed at all
	if false == g.containsCord(a.pos) {
		return false
	}
	existingPois := g.pois[a.pos]
	if len(existingPois) == 0 {
		return true // obviously -- I mean -- there are no others
	}
	for _, point := range existingPois {
		// TODO: access Antinode frequency by performing a 'type switch' on the poi interface
		if point.getIcon() == a.icon {
			return false
		}
	}
	return true
}

func newPoiAntenna(freq byte, pos cord) *poiAntenna {
	// poiAntenna Constructor
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

func newPoiAntinode(pos cord, sources [2]*poi) *poiAntinode {
	// newPoiAntinode Constructor
	// icon is always '#'
	// TODO: add type switch as check weather the sources are poiAntennas
	return &poiAntinode{
		icon:    byte('#'),
		pos:     pos,
		sources: sources,
	}
}

func generateAntinodesModelTsk1(g *grid) []poiAntinode {
	// identify all unique frequencies
	freqSet := make(map[byte]struct{})
	for _, points := range g.pois {
		for _, point := range points {
			switch i := point.(type) {
			case poiAntenna:
				antenna := poiAntenna(i)
				freqSet[antenna.freq] = struct{}{}
			}
		}
	}
	var newAntinodes []poiAntinode
	//it's dirty to assume the freq is same as Icon but, need to change to getPoiAntennasByFreq()
	for freq, _ := range freqSet { // only antennas of same frequency can generate antinodes
		sameFreqAntennas := g.getPoisByIcon(freq)
		for idx, antennaOne := range sameFreqAntennas {
			for _, antennaTwo := range sameFreqAntennas[idx+1:] {
				// In particular, an antinode occurs at any point that is perfectly in line with two antennas of the same frequency
				// but only when one of the antennas is twice as far away as the other
				// hate these type switches -- so ugly -- only the antPosLeft assignment is importend
				var antinodeOne *poiAntinode = nil
				var antinodeTwo *poiAntinode = nil
				switch one := (*antennaOne).(type) {
				case poiAntenna:
					switch two := (*antennaTwo).(type) {
					case poiAntenna:
						var antPosLeft cord = calcPosOfAntinode(one, two)
						var antPosRight cord = calcPosOfAntinode(two, one)
						antinodeOne = newPoiAntinode(antPosLeft, [2]*poi{antennaOne, antennaTwo})
						antinodeTwo = newPoiAntinode(antPosRight, [2]*poi{antennaTwo, antennaOne})
					}
				}
				if antinodeOne != nil {
					newAntinodes = append(newAntinodes, *antinodeOne)
				}
				if antinodeOne != nil {
					newAntinodes = append(newAntinodes, *antinodeTwo)
				}
			}
		}
	}
	return newAntinodes
}

func generateAntinodesModelTsk2(g *grid) []poiAntinode {
	// identify all unique frequencies
	freqSet := make(map[byte]struct{})
	for _, points := range g.pois {
		for _, point := range points {
			switch i := point.(type) {
			case poiAntenna:
				antenna := poiAntenna(i)
				freqSet[antenna.freq] = struct{}{}
			}
		}
	}
	var newAntinodes []poiAntinode
	//it's dirty to assume the freq is same as Icon but, need to change to getPoiAntennasByFreq()
	for freq, _ := range freqSet { // only antennas of same frequency can generate antinodes
		sameFreqAntennas := g.getPoisByIcon(freq)
		for idx, antennaOne := range sameFreqAntennas {
			for _, antennaTwo := range sameFreqAntennas[idx+1:] {
				// In particular, an antinode occurs at any point that is perfectly in line with two antennas of the same frequency
				// but only when one of the antennas is twice as far away as the other
				// hate these type switches -- so ugly -- only the antPosLeft assignment is importend
				switch one := (*antennaOne).(type) {
				case poiAntenna:
					switch two := (*antennaTwo).(type) {
					case poiAntenna:
						// Direciton [ant1]----[ant2]--->[node]--->[node]--->[node]
						var antPositions []cord = calcPosOfAntinodes(one, two, *g)
						for _, position := range antPositions {
							newAntinodes = append(newAntinodes, *newPoiAntinode(position, [2]*poi{antennaOne, antennaTwo}))
						}
						// Direction [node]<---[node]<---[node]<---[ant1]----[ant2]
						antPositions = calcPosOfAntinodes(two, one, *g)
						for _, position := range antPositions {
							newAntinodes = append(newAntinodes, *newPoiAntinode(position, [2]*poi{antennaTwo, antennaOne}))
						}
					}
				}
			}
		}
	}
	return newAntinodes
}

func (g *grid) getPoisByIcon(icon byte) []*poi {
	result := make([]*poi, 0)
	for _, pois := range g.pois {
		for _, point := range pois {
			if point.getIcon() == icon {
				result = append(result, &point)
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func readInput(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func generateAntennasGridFromText(input []byte) grid {
	// assume ascii Encoding
	var out grid
	out.pois = make(map[cord][]poi)
	in := strings.Split(string(input), "\n")
	out.width = len(in[0])   // TODO: add check if each line has the same length
	out.height = len(in) - 1 // account for linefeed at the last line - Split generates extra Slice for that

	for yCord, line := range in {
		for xCord, icon := range line {
			pos := cord{x: xCord, y: yCord}
			antennaP := newPoiAntenna(byte(icon), pos)
			if antennaP != nil {
				out.pois[pos] = append(out.pois[pos], *antennaP)
			}
		}
	}
	return out
}

//---------functions declaration---------

func countAntinodes(g grid) uint {
	var count uint = 0
	for _, pois := range g.pois {
		for _, point := range pois {
			// to much?? maby just the conditional statement. Should have made clear, what defines the type of a POI
			if point.getIcon() == byte('#') {
				count++
			}
		}
	}
	return count
}

func isIntBetween(lower int, upper int, x int) bool {
	// inlcuding lower and excluding upper
	return x >= lower && x < upper
}

func calcPosOfAntinode(one poiAntenna, two poiAntenna) cord {
	// Task One Antinode position calc
	// [ant]----[ant]--->[node]
	//  one	     two	      newAntinode
	distanceOfSources := manhattanDistance(one.pos, two.pos)
	return cord{
		x: two.pos.x + distanceOfSources.x,
		y: two.pos.y + distanceOfSources.y,
	}
}

func calcPosOfAntinodes(one poiAntenna, two poiAntenna, g grid) []cord {
	// Task Two Antinode position calc TODO: make the naming difference between the calcPos..() functions more clear
	// [ant]----[ant]--->[node]--->[node]--->[node]--->[node]--->[node]--->[node]--->[node]--->||Border
	//  one	     two	  newAntinodes
	// need grid to determine the boundaries
	distanceOfSources := manhattanDistance(one.pos, two.pos)
	newCordsList := make([]cord, 0)
	distanceMultiplier := 0
	for {
		nextCord := cord{
			x: two.pos.x + distanceOfSources.x*distanceMultiplier,
			y: two.pos.y + distanceOfSources.y*distanceMultiplier,
		}
		if !g.containsCord(nextCord) {
			break
		}
		distanceMultiplier++
		newCordsList = append(newCordsList, nextCord)
	}
	return newCordsList
}

func manhattanDistance(posOne cord, posTwo cord) cord {
	// calculates a vector rather than an always positive distance
	// the vector points from One to Two
	// TODO: add datatype int out of bound handling
	return cord{
		x: posTwo.x - posOne.x,
		y: posTwo.y - posOne.y,
	}
}
