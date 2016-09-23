# MCPE Tool

It does very little at the moment. It's not really even a tool yet; I just had to pick a project name.
But it can, if set up properly, read a few keys/values from the db folder extracted from a .mcworld file exported from Minecraft Portable Edition.
I am using the Windows 10 Beta version of Minecraft to export the files.

The author has no affiliation with Minecraft, Mojang or Microsoft.

## Goals

I don't know. My original impulse was to create simple block structures--mob spawners, perhaps--in an existing world.
And to create a flat survival world.
Step one was to see if I could even access the world data, and surprisingly I can, so now I'll continue tinkering and seeing what I can read and later what I can write into the world.

Whatever is done, I expect it to mostly be command-line based or at least batch-oriented and not like in-game creative mode or MCEdit.

Some possible near-term goals:

- Create a command-line utility to read information from .mcworld files and/or the extracted db folder
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

## How to make it do anything at all

I'm still in proof-of-concept phase, so it's a bit manual:

- Place https://github.com/midnightfreddie/goleveldb/tree/addzlib (addzlib branch code) into $GOPATH/src/github.com/syndtr/goleveldb/leveldb , or apply [the changes](https://github.com/midnightfreddie/goleveldb/commit/7e93013f9e155f7d70a4bae670630566c6bfc61f) to your local copy of the original repo (had to do this because of fully-qualified import statments) (the changes add zlib decompression as type 2 compression for reading MCPE-modified LevelDB files)
- Unzip an .mcworld file--rename it to .mcworld.zip if needed--then copy the db folder to the folder where you'll be running this program
- Build and execute this program