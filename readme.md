# MCPE Tool

A command line tool to read world data from exported MCPE worlds.
Currently it can read some raw data.
Later it will interpret the data and perhaps even be able to write information into the world.
Code will be split into modules making a non-command-line utility possible in the future.

I am using the Windows 10 Beta version of Minecraft to export the files, but presumably this should also work for Android- and iPhone-exported .mcworld files.

The author has no affiliation with Minecraft, Mojang or Microsoft.

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

- `McpeTool poc` - Runs the original proof-of-concept code I did. This will be removed soon. Expects the db folder to be in your current working directory. It reads and prints the `~local_player` key value and the first 10 keys and values.
- `McpeTool` shows the help screen:

		NAME:
		   MCPE Tool - A utility to access Minecraft Portable Edition .mcworld exported world files.

		USAGE:
		   McpeTool.exe [global options] command [command options] [arguments...]

		VERSION:
		   0.0.0

		COMMANDS:
			 keys, k              Lists all keys in the database. Be sure to include the path to the db, e.g. 'McpeTool keys db'
			 proofofconcept, poc  Run the original POC code which assumes a folder "db" is present with the *.ldb and other level files
			 help, h              Shows a list of commands or help for one command

		GLOBAL OPTIONS:
		   --help, -h     show help
		   --version, -v  print the version

## Goals

I don't know. My original impulse was to create simple block structures--mob spawners, perhaps--in an existing world.
And to create a flat survival world.
Step one was to see if I could even access the world data, and surprisingly I can, so now I'll continue tinkering and seeing what I can read and later what I can write into the world.

Whatever is done, I expect it to mostly be command-line based or at least batch-oriented and not like in-game creative mode or MCEdit.

Some possible near-term goals:

- Print simple statistics on the world (numbers of block types; player/spawn location)
- Place one or more terrain blocks into the map from the utility

Some potential long-term goals:

- Place terrain blocks or possibly complete structures in-world
- Convert Anvil worlds to MCPE
    - This wasn't originally a goal of mine, but it seems to have some interest from others
    - The [jteeuwen/mctools](https://github.com/jteeuwen/mctools) repo is unstarred, but his other projects are well-starred and the project well-documented. I'll look at its Anvil module for this purpose.
    - Create a new playable world
        - And by playable I mean the game will load it and place the player
        - The world itself will probably be very simple and flat
- Visualizations, likely via JSON output files to be read by a web browser page to generate with d3 or similar library
    - Overhead map
    - Simple 3d representation

## References

- http://minecraft.gamepedia.com/Pocket_Edition_level_format