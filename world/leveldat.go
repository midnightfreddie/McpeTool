package world

import "io/ioutil"

// GetLevelDat gets a the contents of level.dat
func (world *World) GetLevelDat() ([]byte, error) {
	value, err := ioutil.ReadFile(world.filePath + "/level.dat")
	return value, err
}

// PutLevelDat replaces the contents of level.dat
func (world *World) PutLevelDat(levelDatData []byte) error {
	err := ioutil.WriteFile(world.filePath+"/level.dat", levelDatData, 0644)
	return err
}
