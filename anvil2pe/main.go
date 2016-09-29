// Taking some code snippets from https://github.com/jteeuwen/mctools/blob/master/anvil for this file
// Also using its modules
// It is BSD "1-clause" licensed. Pretty sure it and MIT are compatible, but will need to ensure I comply with all notifications

package main

import (
	"fmt"

	"github.com/jteeuwen/mctools/anvil"
)

// sectorSize defines the byte size of a single sector.
const sectorSize = 4096

func main() {
	fmt.Println("Get hype")
	region, err := anvil.LoadRegion("region/r.-1.0.mca")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(region.X, region.Z)
}
