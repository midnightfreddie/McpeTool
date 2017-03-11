// Using modules from https://github.com/jteeuwen/mctools/blob/master/anvil for this file
// It is BSD "1-clause" licensed. Pretty sure it and MIT are compatible, but will need to ensure I comply with all notifications

// in early coding, assuming Anvil worlds are always 256 blocks high and PE always 128 blocks high
// assuming all chunks are 16x16

package main

import (
	"os"

	"github.com/jteeuwen/mctools/anvil"
	"github.com/midnightfreddie/McpeTool/world"
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

func int2bytes(i int) [4]byte {
	return [4]byte{byte(i % 256), byte(i >> 8 % 256), byte(i >> 16 % 256), byte(i >> 24 % 256)}
}

func main() {
	// region, err := anvil.LoadRegion("region/r.0.-1.mca")
	region, err := anvil.LoadRegion(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
	world, err := world.OpenWorld("C:\\Users\\Jim.AD\\AppData\\Local\\Packages\\Microsoft.MinecraftUWP_8wekyb3d8bbwe\\LocalState\\games\\com.mojang\\minecraftWorlds\\SuEBAD-5BAA=")
	if err != nil {
		panic(err.Error())
	}
	defer world.Close()
	// fmt.Println(region.X, region.Z)
	// fmt.Printf("%v\n", region.ChunkLen())
	chunks := region.Chunks()
	anvilChunk := anvil.Chunk{}
	for chunkIdx := range chunks {
		// for i := 1; i < 2; i++ {
		success := region.ReadChunk(chunks[chunkIdx][0], chunks[chunkIdx][1], &anvilChunk)
		if !success {
			continue
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
		// fmt.Println(base64.StdEncoding.EncodeToString(peChunk))
		// fmt.Println(int2bytes(chunks[chunkIdx][0]), int2bytes(chunks[chunkIdx][1]))
		// glass ceiling so I can see which chunks were inserted
		for i := 0x70; i < 0x8000; i += 0x80 {
			peChunk[i] = 20
		}
		peKey := make([]byte, 9)
		tmp := int2bytes(chunks[chunkIdx][0])
		copy(peKey, tmp[:])
		tmp = int2bytes(chunks[chunkIdx][1])
		copy(peKey[4:], tmp[:])
		peKey[8] = 0x30
		// fmt.Println(peKey)
		err = world.Put(peKey, peChunk)
		if err != nil {
			panic(err.Error())
		}
	}
}
