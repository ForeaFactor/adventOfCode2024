package day_09

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func Main() {
	input := readInput("./day_09/input.txt")
	diskMap := interpretBytesAsNumerals(input)
	disk := newStorageFromDiskMap(diskMap)
	disk.defragment()

	fmt.Printf("\n====== DAY 8 ======\n")
	fmt.Printf("%d = Checksum Of File After Compacting\n", disk.checksum())
	fmt.Println(string(disk.exportDataToText()))

}

//---------structs declaration---------

type storage struct {
	data []datapoint
}

type datapoint struct {
	// typ 0 = free space | typ 1 = file
	id  int
	typ byte
}

func newStorageFromDiskMap(in []byte) storage {
	// storage Constructor
	// worstCase: each space is 9 Bytes long --> 9*len/(
	out := make([]datapoint, 9*len(in))
	//var out = make([]datapoint, 0)
	var idxOut int = 0    // used as iterator
	var dataType byte = 1 // start with file bytes
	var fileIdx int = 0

	for _, segSize := range in {
		var dp datapoint
		switch dataType {
		case 0:
			dp = datapoint{id: 0, typ: dataType}
		case 1:
			dp = datapoint{id: fileIdx, typ: dataType}
			fileIdx++
		default:
			dp = datapoint{id: 0, typ: dataType}
		}
		for i := 0; byte(i) < segSize; i++ {
			out[idxOut] = dp
			idxOut++
		}
		dataType = dataType ^ 1 // alternating pattern of file and space bytes
	}
	return storage{data: out[:idxOut]}
}

func (s *storage) defragment() {
	var iterFront int = -1 // assign outOfBound Indicies cause nextIter function increments directly
	var iterBack int = len(s.data)
	for iterFront < iterBack {
		// next emptyStoragePoint typ 0 - iterFront is advanced
		if _, err := s.nextDatapoint(&iterFront, 0); err == io.EOF {
			panic(err)
		}
		// next fileDataPoint typ 1 - iterBack is advanced
		if _, err := s.prevDatapoint(&iterBack, 1); err == io.EOF {
			panic(err)
		}
		if iterFront < iterBack /* Condition because then all spaces have been filled*/ {
			s.swapDatapoints(iterFront, iterBack)
		}
		//fmt.Printf("%s\n", s.exportDataToText())
	}
}

//---------methods declaration---------

func (s *storage) exportDataToText() string {
	var out string
	for _, dataPoint := range s.data {
		switch dataPoint.typ {
		case 0:
			out += "."
		case 1:
			out += "[" + strconv.Itoa(dataPoint.id) + "]" // this results in wierd characters
		}
	}
	return out
}

func (s *storage) nextDatapoint(iter *int, typ byte) (datapoint, error) {
	//scan storage for the next dataPoint of type typ
	// returns empty datapoint when eof
	for {
		*iter = *iter + 1
		if *iter >= len(s.data) || *iter < 0 {
			*iter = len(s.data) // stop iter progress here
			return *new(datapoint), io.EOF
		}
		if s.data[*iter].typ == typ {
			return s.data[*iter], nil
		}
	}
}

func (s *storage) prevDatapoint(iter *int, typ byte) (datapoint, error) {
	// scan storage for the prev dataPoint of type typ
	// returns empty datapoint when eof
	for {
		*iter = *iter - 1
		if *iter >= len(s.data) || *iter < 0 {
			*iter = 0 // stop iter progress here
			return *new(datapoint), io.EOF
		}
		if s.data[*iter].typ == typ {
			return s.data[*iter], nil
		}
	}
}

func (s *storage) swapDatapoints(front int, back int) {
	tmp := s.data[front]
	s.data[front] = s.data[back]
	s.data[back] = tmp
}

func (s *storage) checksum() uint64 {
	var checksum uint64 = 0
	var pointPos int = 0
	for _, dataPoint := range s.data {
		if dataPoint.typ != 0 {
			checksum += uint64(pointPos * dataPoint.id)
			pointPos++
		}
	}
	return checksum
}

//---------functions declaration---------

func interpretBytesAsNumerals(input []byte) []byte {
	var output []byte = make([]byte, len(input))
	for i, byt := range input {
		// assumes only [0-9]
		output[i] = byte(int(rune(byt) - '0'))
	}
	return output
}

func readInput(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
