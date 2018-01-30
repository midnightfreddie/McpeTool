package world

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
)

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

// GetLevelDatNbtAndVersion returns the version and nbt portions of level.dat
func (world *World) GetLevelDatNbtAndVersion() ([]byte, int32, error) {
	var version int32
	value, err := world.GetLevelDat()
	if err != nil {
		return nil, 0, err
	}
	buf := bytes.NewBuffer(value[:4])
	err = binary.Read(buf, binary.LittleEndian, &version)
	return value[8:], version, err
}

// PutLevelDatNbtAndVersion builds the header info for level.dat given the nbt and version, and then replaces level.dat
func (world *World) PutLevelDatNbtAndVersion(levelDatData []byte, version int32) error {
	var header []byte
	var err error
	nbtLen := int32(len(levelDatData))
	buf := bytes.NewBuffer(header)
	err = binary.Write(buf, binary.LittleEndian, version)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.LittleEndian, nbtLen)
	if err != nil {
		return err
	}
	err = world.PutLevelDat(append(buf.Bytes(), levelDatData...))
	return err
}
