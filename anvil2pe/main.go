// Using modules from https://github.com/jteeuwen/mctools/blob/master/anvil for this file
// It is BSD "1-clause" licensed. Pretty sure it and MIT are compatible, but will need to ensure I comply with all notifications

// in early coding, assuming Anvil worlds are always 256 blocks high and PE always 128 blocks high
// assuming all chunks are 16x16

package main

import (
	"encoding/base64"
	"fmt"

	"github.com/jteeuwen/mctools/anvil"
)

// does not validate for sanity
func anvilOffset(x, y, z int) (section, offset int) {
	tmp := 256*y + 16*z + x
	section = tmp / 4096
	offset = tmp % 4096
	return
}

// does not validate for sanity
func peOffset(x, y, z int) (offset int) {
	offset = 2048*x + 128*z + y
	return
}

func main() {
	region, err := anvil.LoadRegion("region/r.-1.0.mca")
	if err != nil {
		panic(err.Error())
	}
	// fmt.Println(region.X, region.Z)
	// fmt.Printf("%v\n", region.ChunkLen())
	chunks := region.Chunks()
	anvilChunk := anvil.Chunk{}
	// for i := range chunks {
	for i := 1; i < 2; i++ {
		success := region.ReadChunk(chunks[i][0], chunks[i][1], &anvilChunk)
		if !success {
			continue
		}
	}
	peChunk := make([]byte, 83200)
	// Full brighness for blocks instead of full dark
	for i := 0xc000; i < 0x14000; i++ {
		peChunk[i] = 0xff
	}
	// Grass color
	grassColor := [...]byte{0x7, 0x7d, 0xac, 0x6c}
	for i := 0x14100; i < 0x14500; i += 4 {
		copy(peChunk[i:], grassColor[:])
	}
	for sIdx := range anvilChunk.Sections {
		section := anvilChunk.Sections[sIdx]
		if section.Y > 7 {
			continue
		}
		yBase := 16 * int(section.Y)
		for y := yBase; y < yBase+16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {
					_, aIdx := anvilOffset(x, y, z)
					peChunk[peOffset(x, y, z)] = section.Blocks[aIdx]
				}
			}
		}
	}
	fmt.Println(base64.StdEncoding.EncodeToString(peChunk))
}
