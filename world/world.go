package world

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/midnightfreddie/goleveldb/leveldb"
)

// World holds the LevelDB instance, wraps its functions and provides functions for any other World needs
type World struct {
	db       *leveldb.DB
	filePath string
}

// OpenWorld opens a Minecraft Pocket Edition world folder
func OpenWorld(path string) (World, error) {
	world := World{nil, path}
	var err error
	dbPath := path + "/db"

	// For now, abort if path/db doesn't exist or if it's not a directory . Later may add an option to create if not exist or otherwise validate world
	fileInfo, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		return world, errors.New(dbPath + " does not exist. This must be run against a valid world folder.")
	}
	if !fileInfo.IsDir() {
		return world, errors.New(dbPath + " is not a directory. This must be run against a valid world folder.")
	}

	world.db, err = leveldb.OpenFile(dbPath, nil)
	if err != nil {
		_ = world.db.Close()
		return world, err
	}
	return world, nil
}

// Close needs to be called before exiting to relase the LevelDB locks
func (world *World) Close() error {
	err := world.db.Close()
	return err
}

// FilePath returns the path used to open the world
func (world *World) FilePath() string {
	return world.filePath
}

// GetKeys returns all keys in the LevelDB database
func (world *World) GetKeys() ([][]byte, error) {
	keylist := [][]byte{}
	iter := world.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		tmp := make([]byte, len(key))
		copy(tmp, key)
		keylist = append(keylist, tmp)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return keylist, err
	}
	return keylist, nil
}

// // GetPlayerKeys returns player keys
// func (world *World) GetPlayerKeys() ([][]byte, error) {
// 	keylist := [][]byte{}
// 	iter := world.db.NewIterator(nil, nil)
// 	for iter.Next() {
// 		key := iter.Key()
// 		tmp := make([]byte, len(key))
// 		copy(tmp, key)
// 		if isPlayer(key) {
// 			keylist = append(keylist, tmp)
// 		}
// 	}
// 	iter.Release()
// 	err := iter.Error()
// 	if err != nil {
// 		return keylist, err
// 	}
// 	return keylist, nil
// }

// Get gets a value from the world database
func (world *World) Get(key []byte) ([]byte, error) {
	tmp, err := world.db.Get(key, nil)
	// goleveldb docs say not to modify the returned slice. Unsure if that would create a problem or if they mean you can't update the DB that way
	// In any case, returning a copy of the slice to be safe
	value := make([]byte, len(tmp))
	copy(value, tmp)
	return value, err
}

// Put puts a key/value pair into the world database, superceding/replacing/deleting the existing value for that key, if any
func (world *World) Put(key []byte, value []byte) error {
	err := world.db.Put(key, value, nil)
	return err
}

// Delete deletes a key and its value from the world database
func (world *World) Delete(key []byte) error {
	err := world.db.Delete(key, nil)
	return err
}

// GetLevelDat gets a the contents of level.dat
func (world *World) GetLevelDat() ([]byte, error) {
	// levelDatPath := world.filePath + "/level.dat"
	value, err := ioutil.ReadFile(world.filePath + "/level.dat")
	return value, err
}
