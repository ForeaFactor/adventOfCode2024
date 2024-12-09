package day_09

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func Main() {
	input := readInput("./day_09/sample.txt")
	diskMap := interpretBytesAsNumerals(input)
	disk := newStorageFromDiskMap(diskMap)
	disk.defragment()

	fmt.Printf("\n====== DAY 8 ======\n")
	fmt.Printf("%d = \n", 0)
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
	var iterFront int = 0
	var iterBack int = len(s.data) - 1

	for iterFront < iterBack {
		s.swapDatapoints(iterFront, iterBack)

	}

	iterFront, err = s.nextDatapoint(iterFront, 0)
	iterBack, err = s.prevDatapoint(iterBack, 1)
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

	for {
		*iter = *iter + 1
		if iter >= len(s.data) || iter < 0 {
			return -1, io.EOF
		}
		if s.data[iter].typ == typ {
			return iter, nil
		}
	}
}

func (s *storage) prevDatapoint(iter int, typ byte) (int, error) {
	//scan storage for the next dataPoint of type typ
	for {
		iter--
		if iter >= len(s.data) || iter < 0 {
			return -1, io.EOF
		}
		if s.data[iter].typ == typ {
			return iter, nil
		}
	}
}

func (s *storage) swapDatapoints(front int, back int) {
	tmp := s.data[front]
	s.data[front] = s.data[back]
	s.data[back] = tmp
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
