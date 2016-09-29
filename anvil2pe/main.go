// Taking some code snippets from https://github.com/jteeuwen/mctools/blob/master/anvil/anvil-dump/region.go for this file
// Also using its modules
// It is BSD "1-clause" licensed. Pretty sure it and MIT are compatible, but will need to ensure I comply with all notifications

package main

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// sectorSize defines the byte size of a single sector.
const sectorSize = 4096

func doChunk(r io.ReadSeeker, offset int64) error {
	address := offset*sectorSize + 4
	_, err := r.Seek(address, 0)
	if err != nil {
		return err
	}
	var scheme [1]byte
	_, err = io.ReadFull(r, scheme[:])
	if err != nil {
		return err
	}
	var rr io.ReadCloser
	switch scheme[0] {
	case 1:
		rr, err = gzip.NewReader(r)
	case 2:
		rr, err = zlib.NewReader(r)
	default:
		return fmt.Errorf("chunk(%d); invalid compression scheme: %d", offset, scheme[0])
	}

	if err != nil {
		return err
	}

	// err = dump(w, rr)
	rr.Close()
	return nil
}

func main() {
	fmt.Println("Get hype")
	var locations [sectorSize]byte
	anvilFile, err := os.Open("region/r.-1.0.mca")
	if err != nil {
		panic(err.Error())
	}
	defer anvilFile.Close()
	_, err = io.ReadFull(anvilFile, locations[:])
	if err != nil {
		panic(err.Error())
	}
	// doing one chunk for now
	d := locations[32:]
	offset := int64(d[0])<<16 | int64(d[1])<<8 | int64(d[2])
	sectors := d[3]
	if offset == 0 && sectors == 0 {
		fmt.Println("offset and sectors are 0")
		return
	}
	err = doChunk(anvilFile, offset)
	if err != nil {
		panic(err.Error())
	}
}
