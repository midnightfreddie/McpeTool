# MCPE Tool

A command line tool to read and write world data from Minecraft Pocket Edition worlds.
Currently it can read raw data from command line or web API.
It's capable of writing, and game-tested, but I'm currently working on the API to write data.
Later it will interpret the data and perhaps even be able to write information into the world.
Code is now in modules making a non-command-line utility possible.

I am using the Windows 10 Beta version of Minecraft to export and import the .mcworld files, but this should also work for Android and iPhone world directories if you manually copy them from/to the device.

The author has no affiliation with Minecraft, Mojang or Microsoft.

## Use

- **Back up your worlds**
- Copy an MCPE world to a working folder
	- Windows 10 Edition: Unzip an .mcworld file--rename it to .mcworld.zip if needed--so that the contents including the "db" folder are accessible. (In the future this utility may be able to read from .mcworld files directly.)
	- For Android or iOS, backup and copy the worlds via USB or file manager app (transfer via cloud drive or Bluetooth)
	- The world folder will include a "level.dat" file and a "db" directory

- `McpeTool api path/to/world` - Starts a web server on port 8080 allowing REST API access to the world. http://localhost:8080/api/v1/db/ will return the DB keys. Keys and data are base64 encoded.
- `McpeTool keys path/to/world` - This will list the keys in the LevelDB world store in base64 format
- `McpeTool get path/to/world base64key` - Returns the data for the given key in base64 format. Example: `McpeTool.exe get path/to/world fmxvY2FsX3BsYXllcg==` returns player data in base64 format.
- `put` is not yet implemented via command line because the data would be too big for parameters
- `McpeTool delete path/to/world base64key` - Deletes the key/value pair for that key if present. Example: `McpeTool.exe delete path/to/world fmxvY2FsX3BsYXllcg==` deletes the player data. If you do this and play the world, you will spawn at the original location with no inventory.
- `McpeTool` shows the help screen:

		NAME:
		   MCPE Tool - A utility to access Minecraft Pocket Edition .mcworld exported world files.

		USAGE:
		   McpeTool.exe [global options] command [command options] [arguments...]

		VERSION:
		   0.0.0

		COMMANDS:
			 api, www      Open world and start http API. Hit control-c to exit.
			 keys, k       Lists all keys in the database in base64 format. Be sure to include the path to the world folder, e.g. 'McpeTool keys path/to/world'
			 get           Retruns the value of a key. Both key and value are in base64 format. e.g. 'McpeTool get path/to/world AAAAAAAAAAAw' for terrain chunk 0,0 or 'McpeTool get path/to/world fmxvY2FsX3BsYXllcg==' for ~local_player player data
			 delete        Deletes a key and its value. The key is in base64 format. e.g. 'McpeTool delete path/to/world AAAAAAAAAAAw' to delete terrain chunk 0,0 or 'McpeTool delete path/to/world fmxvY2FsX3BsYXllcg==' to delete ~local_player player data
			 develop, dev  Random thing the dev is working on
			 help, h       Shows a list of commands or help for one command

		GLOBAL OPTIONS:
		   --help, -h     show help
		   --version, -v  print the version

## Goals

My original impulse was to create simple block structures--mob spawners, perhaps--in an existing world.
And to create a flat survival world.
My early attempts to read and later write the world were successful, and then I realized that a simple API would make this program versatile and not break every time MCPE updates.
My main focus is to allow low-level access via simple APIs so more complex logic can be handled by any program in any language.

- Read/Write command-line access to worlds
- Read/Write local REST API access to worlds
	- Allow raw key get/put/delete
	- Allow raw get/put/delete to level.dat and level.txt
	- Provide basic in-browser web app to do simple world edits using API
- Friendlier API allowing access to blocks, entities, players and other game settings (villages? portals?). Possibly also:
	- Find/replace blocks
	- Place predefined structures
- For friendly APIs, allow user to provide config files so the program should be usable on future world version updates
- Allow reading/writing .mcworld files without manual unzip/zip
- Auto-detect if provided path is a world directory, LevelDB directory or zipped exported world
- Print simple statistics on the world (numbers of block types; player/spawn location)
- (Maybe) Convert Anvil worlds to LevelDB
    - This wasn't originally a goal of mine, but it seems to have some interest from others
    - The [jteeuwen/mctools](https://github.com/jteeuwen/mctools) repo is unstarred, but his other projects are well-starred and the project well-documented. I'll look at its Anvil module for this purpose.
	- First problem: Anvil worlds are 256 blocks high while PE worlds are 128 blocks high.
- Visualizations, likely via JSON output files to be read by a web browser page to generate with d3 or similar library
    - Overhead map
    - Simple 3d representation

## Accomplishments

- Read raw keys and data from db/ folder LevelDB database (manually unzipped from .mcworld file)
- Wrote terrain blocks—a column of diamond blocks from y=0 to y=127—into an existing survival mode world and was able to continue playing that level in Win10 Edition and Android. (Manually added modified db folder into .mcworld for Win10, and manually copied world folders to Android.)
- Wrote 40 chunks from 1-40 on the x axis beyond where terrain was pre-generated
	- Written chunks had glass bottom, a water layer and jack o'lantern pillar and "stairs" around the chunk perimeter
	- It was playable
	- The chunks surrounding the new chunks generated automatically when playing
	- Tree leaves would generate protruding into my placed chunks' air space
	- The generated terrain tried to match my chunks, but since my chunks had a "spiral staircase" around the perimeter the results were interesting (uneven on each side of the chunk; generally very mountainous)
	- Jack o'lanterns and water are supposed to have additional data in the db, but MCPE seems tolerant of this missing data and assigns sane defaults, apparently
- Blanked out db, only put in very simple row of chunks
	- level.dat still in place
	- Game spawned the player on the chunk ground and generated surrounding terrain

## Notes

### World format

An .mcworld file is the zipped contents of a world folder. .mcworld files can be exported from and imported to Minecraft Windows 10 edition.
The Android version (and presumably iOS version) do not seem to have export/import commands, although the world format is reportedly the same.

- World folder
	- level.dat - NBT-coded information about the world
	- level.dat_old - a backup version of level.dat?
	- levelname.txt - Self-explanatory
	- db folder - This is a LevelDB database. **The files should not be directly altered.** The db as a whole should be accessed from a LevelDB library. Unlike databases most people have heard of, there is no query language or general client program for this database.
		- CURRENT - not for humans
		- LOCK - not for humans
		- LOG - Human-readable LevelDB event log.
		- MANIFEST - not for humans
		- *.ldb - LevelDB data files
		- *.log - *Not* a readable file. It's a binary DB log, not anything a human would be interested in reading.

### Minecraft data in LevelDB

- LevelDB is a simple key/value store.
- Structure of MCPE keys and values: http://minecraft.gamepedia.com/Pocket_Edition_level_format

### World file locations

#### Windows 10

- Windows 10 can export and import .mcworld files (which are zip files containing the world data) in the game interface to any location you choose.
- The actual locations of the saved worlds is under `%LOCALAPPDATA%\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds` (if you paste that into an explorer window it will show you the right location)

#### Android

- In one of the storage volumes as `\games\com.mojang\minecraftWorlds`

#### iOS

Sorry, I don't have an iOS device.