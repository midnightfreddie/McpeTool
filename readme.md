# MCPE Tool

A command line tool to read world data from exported MCPE worlds.
Currently it can read and write raw, hard-coded data.
Later it will interpret the data and perhaps even be able to write information into the world.
Code will be split into modules making a non-command-line utility possible in the future.

I am using the Windows 10 Beta version of Minecraft to export the files, but presumably this should also work for Android- and iPhone-exported .mcworld files.

The author has no affiliation with Minecraft, Mojang or Microsoft.

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

## Use

- Unzip an .mcworld file--rename it to .mcworld.zip if needed--so that the "db" folder is accessible. (In the future this utility may be able to read from .mcworld files directly.)

- `McpeTool keys path/to/db` - This will list the keys in the LevelDB world store. Sample partial output:

		[11 0 0 0 255 255 255 255 48]
		[11 0 0 0 255 255 255 255 118]
		[12 0 0 0 254 255 255 255 48]
		[12 0 0 0 254 255 255 255 118]
		[12 0 0 0 255 255 255 255 48]
		[12 0 0 0 255 255 255 255 118]
		BiomeData
		mVillages
		portals
		~local_player
		[244 255 255 255 0 0 0 0 48]
		[244 255 255 255 0 0 0 0 49]
		[244 255 255 255 0 0 0 0 118]
		[244 255 255 255 252 255 255 255 48]
		[244 255 255 255 252 255 255 255 118]
		[245 255 255 255 2 0 0 0 48]

- `McpeTool api path/to/db` - Starts httpd daemon on localhost:8080. Any web requests to it will return a JSON-endoded list of keys in the database (each base64-encoded) 
- `McpeTool` shows the help screen:

	NAME:
	   MCPE Tool - A utility to access Minecraft Portable Edition .mcworld exported world files.

	USAGE:
	   McpeTool.exe [global options] command [command options] [arguments...]

	VERSION:
	   0.0.0

	COMMANDS:
		 keys, k       Lists all keys in the database. Be sure to include the path to the db, e.g. 'McpeTool keys db'
		 develop, dev  Random thing the dev is working on
		 api, www      Open world and start http API. Hit control-c to exit.
		 help, h       Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	   --help, -h     show help
	   --version, -v  print the version

## Goals

I don't know. My original impulse was to create simple block structures--mob spawners, perhaps--in an existing world.
And to create a flat survival world.
Step one was to see if I could even access the world data, and surprisingly I can, so now I'll continue tinkering and seeing what I can read and later what I can write into the world.

Whatever is done, I expect it to mostly be command-line based or at least batch-oriented and not like in-game creative mode or MCEdit.

Some possible near-term goals:

- Modularize world access
	- Allow get/put of individual terrain blocks (no entity/expanded data yet)
	- Wrap chunk read/write calls
- Allow reading/writing .mcworld files without manual unzip/zip
- Auto-detect if provided path is a world directory, LevelDB directory or zipped exported world
- Print simple statistics on the world (numbers of block types; player/spawn location)

Some potential long-term goals:

- Replace terrain blocks with other types
- Place complete structures
- Convert Anvil worlds to MCPE
    - This wasn't originally a goal of mine, but it seems to have some interest from others
    - The [jteeuwen/mctools](https://github.com/jteeuwen/mctools) repo is unstarred, but his other projects are well-starred and the project well-documented. I'll look at its Anvil module for this purpose.
    - Create a new playable world
        - And by playable I mean the game will load it and place the player
        - The world itself will probably be very simple and flat
	- First problem: Anvil worlds are 256 blocks high while PE worlds are 128 blocks high.
- Visualizations, likely via JSON output files to be read by a web browser page to generate with d3 or similar library
    - Overhead map
    - Simple 3d representation

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