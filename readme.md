# MCPE Tool

A tool to read and write world data from Minecraft Pocket Edition.
Currently it can read and write raw key/value data from command line or via HTTP REST API.
For example, each chunk's terrain is stored in a single key/value pair, and the player is 

Game-tested against Minecraft Windows 10 Edition and Android Minecraft Pocket Edition.

The author has no affiliation with Minecraft, Mojang or Microsoft.

## Use

- **Back up your worlds**
- Copy an MCPE world to a working folder
	- Windows 10 Edition: Unzip an exported .mcworld file--rename it to .mcworld.zip if needed--so that the contents including the "db" folder are accessible. (In the future this utility may be able to read from .mcworld files directly.)
	- Also, in Windows 10 the save folders can be found in `%LOCALAPPDATA%\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds` (if you paste that into an explorer window it will show you the right location)
	- For Android or iOS, backup and copy the worlds via USB or file manager app (transfer via cloud drive or Bluetooth)
	- The world folder will include a "level.dat" file and a "db" directory
	- Set the path to the world folder in the MCPETOOL_WORLD environment variable or pass it to mcpetool with the --path or -p parameter. If no path is provided, it will attempt to use the current folder as the world path.

- `McpeTool keys` - This will list the keys in the LevelDB world store in hex format
- `McpeTool get [--dump] <hexkey>` - Returns the data for the given key in base64 format
	- `--dump` or `-d` flag outputs as hexdump instead
	- `--json` or `-j` flag attempts nbt2json conversion and outputs json. It will not produce an error if the value is not nbt; it will output `{ "tagType": 0, "name": "" }` which doesn't represent the value. Not all keys are NBT!
	- `--yaml` or `y` - exactly like `--json` but then converted to YAML
	- Example: `McpeTool.exe get 7e6c6f63616c5f706c61796572` returns player data
- `McpeTool put <hexkey>` - Puts a key/value pair in the database, replacing the previous value if present or creating the key if not.
	- `--json` or `-j` flag attempts json2nbt conversion on input data. The best use case for this is to get a value in json format, edit it, and then put the updated value back. This does not validate that the resulting NBT is valid game data.
	- `--yaml` or `y` - exactly like `--json` but converts input from YAML to JSON internally first
- `McpeTool delete <hexkey>` - Deletes the key/value pair for that key if present
	- Example: `McpeTool.exe delete 7e6c6f63616c5f706c61796572` deletes the player data. If you do this and play the world, you will spawn at the original location with no inventory.
- `McpeTool api` - Starts a local HTTP server allowing REST API access to the MCPE world database. When done accessing the world, stop the server with control-C (or your OS'es equivalent BREAK)
	- http://localhost:8080/api/v1/db/ will return the DB keys. Sample parital output:

			{
			  "apiVersion": "1.0",
			  "keys": [
				{
				  "stringKey": "AutonomousEntities",
				  "hexKey": "4175746f6e6f6d6f7573456e746974696573"
				},
				{
				  "stringKey": "BiomeData",
				  "hexKey": "42696f6d6544617461"
				},
				{
				  "stringKey": "Overworld",
				  "hexKey": "4f766572776f726c64"
				},
				{
				  "stringKey": "mVillages",
				  "hexKey": "6d56696c6c61676573"
				},
				{
				  "stringKey": "portals",
				  "hexKey": "706f7274616c73"
				},
				{
				  "hexKey": "7e000000000000002d"
				},
				{
				  "hexKey": "7e000000000000002f00"
				},
				// <truncated for example>
				{
				  "hexKey": "8a000000f8ffffff36"
				},
				{
				  "hexKey": "8a000000f8ffffff76"
				}
			  ]
			}

	- Can retrieve, e.g. http://localhost:8080/api/v1/db/83000000ffffffff2f04 for terrain subchunk X=2098 Z=-16 Y=65, with HTTP GET requests, or a simple web browser

				{
					"apiVersion": "1.0",
					"hexKey": "83000000ffffffff2f04",
					"base64Data": "AAAAAAAAAAAAAAAAAA <truncated for example> AAAAAAAAAAAAA=="
				}

	- Can also PUT and DELETE keys via web API using tools and language of your choice   

## How to convert world coordinates to leveldb keys

All division below is of course integer division. The remainder/modulus will be used to find the byte offset within the subchunk data. Also, X, Z, and dimension are 32-bit signed integers in little endian byte order. The true key is a byte array, but we represent it with a hex string for convenience. In the examples below, I've bolded the chunk Z coordinate for clarity.

Each chunk is 16x16x256 (X,Z,Y), and the subchunk block data keys are 16 high. So for x, z, y coordinates of 413, 54, 105:

- chunk X = 413 / 16 = 25 or 0x19 signed 32-bit integer in little endian byte order ([0x19,0, 0, 0] == 19000000)
- chunk Z = 54 / 16 = 3 ([0x3, 0, 0, 0] == **03000000**) 

So all keys beginning with 19000000**03000000** are about this coordinate's chunk. (In the overworld; other dimensions add a 32-bit dimension ID, so the same coordinates in the Nether I think have keys that start with 19000000**03000000**FFFFFFFF and 19000000**03000000**01000000 for the End.)

The tags and subchunk indexes are 8-bit values. (Unsigned? Not sure it matters as there are no negative Y chunk coordinates and no tags <0 or > 127.)

47 ([0x2F]) is the subchunk prefix tag, so all keys beginning with 19000000**03000000**2f are the Y subchunks for this coordinate.

- subchunk Y = 105 / 16 = 6 ([0x*06*])

So, the subchunk key for X=413, Z=54, Y=105 is 19000000**03000000**2f*06*

ref: https://minecraft.gamepedia.com/Bedrock_Edition_level_format

## How to convert world coordinates to block subchunk byte offsets

`%` is modulo division, the remainder of integer division. Note that subchunks may not exist if they are empty.

- `(X % 16) * 256 + (Z % 16) * 16 + (Y % 16) + 1` - 1-byte Block ID that determines if the block is e.g. air, water, sand, dirt, tree, lava, etc.
- The following three are odd in that each byte holds data for two blocks, in the high and low nybbles.
- `(X % 16) * 32 + (Z % 16) * 8 + (Y % 16) / 2 + 4096 + 1` - Block data e.g. rotation
- `(X % 16) * 32 + (Z % 16) * 8 + (Y % 16) / 2 + 6144 + 1` - Sky light
- `(X % 16) * 32 + (Z % 16) * 8 + (Y % 16) / 2 + 8192 + 1` - Block light

## Goals

My original impulse was to create simple block structures--mob spawners, perhaps--in an existing world.
And to create a flat survival world.
My early attempts to read and later write the world were successful, and then I realized that a simple API would make this program versatile and not break every time MCPE updates.
My main focus is to allow low-level access via simple APIs so more complex logic can be handled by any program in any language.

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
	- Challenges:
		- Anvil worlds are 256 blocks high while PE worlds are 128 blocks high
		- At least some block IDs are different and/or do not exist in both versions
		- Data is organized differently
		- Hopefully NBT data is mostly tagged the same way and portable at the Compund Tag level
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
- Deleted terrain chunks
	- Somehow the game preserved or recovered the deleted chunks, complete with in-game player modifications
- Deleted player
	- Player spanwed at original spawn point with no inventory
- Used a script against the web API to replace blocks in chunk 0,0
- Used the command line to get the contents of chunk 0,0 and put them in chunk 1,0 to duplicate a chunk (**Note**: this was pre-MCPE-1.0 when chunks were 128-high and in keys ending in 0x30 for terrain. In 1.0 and later, terrain is stored in subchunks with keys ending in 0x2fnn where nn is the subchunk number. Also, the tool now uses a flag or env var to specify world path.)
	- `.\McpeTool.exe get "<path/to/world>" 000000000000000030 | .\McpeTool.exe put "<path/to/world>" 000000000100000030`

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
- Keys of interest
	- 9-byte-log keys ending in 0x30 are terrain chunks with terrain and lighting data. Fixed size, no NBT coding.
	- 13-byte-log keys ending in 0x30 are Nether terrain chunks, same structure as overworld chunks
	- `42696f6d6544617461` - "BiomeData"
	- `6d56696c6c61676573` - "mVillages"
	- `706f7274616c73` - "portals"
	- `7e6c6f63616c5f706c61796572` - "~local_player" - NBT-coded player data. The other string keys are also likely NBT-coded. 

### World file locations

#### Windows 10

- Windows 10 can export and import .mcworld files (which are zip files containing the world data) in the game interface to any location you choose.
- The actual locations of the saved worlds is under `%LOCALAPPDATA%\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds` (if you paste that into an explorer window it will show you the right location)

#### Android

- In one of the storage volumes as `\games\com.mojang\minecraftWorlds`

#### iOS

Sorry, I don't have an iOS device.