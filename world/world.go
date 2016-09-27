package world

import "github.com/midnightfreddie/goleveldb/leveldb"

// World holds the LevelDB instance, wraps its functions and provides functions for any other World needs
type World struct {
	db *leveldb.DB
}

// OpenWorld opens a Minecraft Pocket Edition world folder
func OpenWorld(path string) (World, error) {
	world := World{}
	var err error
	world.db, err = leveldb.OpenFile(path+"/db", nil)
	if err != nil {
		panic("error")
	}
	return world, nil
}

// Close needs to be called before exiting to relase the LevelDB locks
func (world *World) Close() error {
	err := world.db.Close()
	return err
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
		panic(err.Error())
	}
	return keylist, nil
}
