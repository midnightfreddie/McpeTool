# MCPE Tool

An offline tool to read and write world data from Minecraft Pocket Edition and Win10 Edition (Bedrock). It can also read and write the raw LevelDB keys for Minecraft 9.0.

The author has no affiliation with Minecraft, Mojang or Microsoft.

## Use

- Since LevelDB only allows one process to use it at a time, this tool can not be used while the world is running in Minecraft
- Input and most output defaults to stdin/stdout, but `--in` or `--out` parameters are available where appropriate. Some key lists only print to stdout.
- This is a work in progress. I intend to release only working versions, but always be sure to back up your worlds before modifying them.
- This tool does not currently validate that data input is valid for the game, but the game has proven pretty resilient to my messing around with saves. It seems to generate new terrain or respawn the player as appropriate if it encounters missing or invalid data.
- `mcpetool [command [subcommand]] -h` - Shows help and options for commands and subcommands.
- You can put the path to your world in the environment variable `MCPETOOL_WORLD` to save typing on repeated commands
- If `%MCPETOOL_WORLD%`/`$MCPETOOL_WORLD`/`$env:MCPETOOL_WORLD` isn't set and `--path` isn't specified, `mcpetool` will look in the current directory for a world
- NBT to JSON and YAML formatting is using [my nbt2json module](https://github.com/midnightfreddie/nbt2json). When using JSON or YAML as input, it needs to be compatible with this. That basically means the top-level document key `"nbt":` should be an array of objects like `{ "name": "<name>", "tagType": <number>, "value": <varies based on type> }`, and of course sub-tags of compound tags and list tags are arrays of nbt2json objects in the value field.

### level.dat

- `mcpetool leveldat get [--path <path/to/world>] [--out <file>] [--dump] [--yaml] [--base64] [--binary]` - Retrieves the contents of level.dat in JSON format, or in hex dump, YAML, base64, or binary format if one of those options is given. For JSON and YAML output, the level.dat version is put in the comment field, and the length is discarded. For hex dump, base64 and binary output, you get the unmodified contents of the file including the 8-byte header.
- `mcpetool leveldat put [--path <path/to/world>] [--in <file>] [--ver <version>] [--yaml] [--base64] [--binary]` - Replaces the contents of level.dat with your JSON input, or YAML, base64, or binary input per the parameters. For JSON and YAML input, the header is generated using version 6 or the version number provided by `--ver` (not from the JSON/YAML), and the length is calculated after conversion. For base64 or binary input, the level.dat file is written with those contents exactly, so they should include the proper 8-byte header.

### LevelDB

Most world data is stored in a modified LevelDB key/value store. The leveldb commands provide raw access to this key/value store. ref: https://minecraft.gamepedia.com/Bedrock_Edition_level_format

- `mcpetool db list [--path <path/to/world>]` - This will list the keys in the LevelDB world store in hex string format
- `mcpetool db get [--path <path/to/world>] [--json] [--dump] [--yaml] [--base64] [--binary] <hexkey>` - Returns the data for the given key in base64 or specified format
	- Example: `mcpetool db get --path path/to/world --yaml 7e6c6f63616c5f706c61796572` returns local player data in YAML format
	- Example: `mcpetool db get --path path/to/world --json 00000000000000002f00` returns chunk X=0, Z=0, Y=0 in JSON format if it exists
- `mcpetool db put [--path <path/to/world>] [--json] [--yaml] [--base64] [--binary] <hexkey>` - Puts a key/value pair in the database, replacing the previous value if present or creating the key if not. They key and value are not checked for game validity; it will place any data in any key you specify.
- `mcpetool db delete [--path <path/to/world>] <hexkey>` - Deletes the key/value pair for that key if present
	- Example: `mcpetool db delete [--path <path/to/world>] 7e6c6f63616c5f706c61796572` deletes the local player data including inventory and equipped items. If you do this and play the world, you will spawn at the world spawn point.

### API

Starts a local HTTP server allowing REST API access to the MCPE world database. When done accessing the world, stop the server with control-C (or your OS'es equivalent BREAK).

The PUT request bodies should be formatted like the GET requests with base64-encoded data in `"base64Data":`.

I intend to make the next API version with Swagger, and add more features and data input/output options.

- `mcpetool api [--path <path/to/world>]` - Starts the local REST API server
	- GET http://localhost:8080/api/v1/db/ will return the DB keys
	- GET http://localhost:8080/api/v1/db/7e6c6f63616c5f706c61796572 - Returns `7e6c6f63616c5f706c61796572` in a JSON object with the base64-encoded value.
	- PUT http://localhost:8080/api/v1/db/7e6c6f63616c5f706c61796572 - Creates/overwrites `7e6c6f63616c5f706c61796572` with request body formatted as seen with the GET requests
	- DELETE http://localhost:8080/api/v1/db/7e6c6f63616c5f706c61796572 - Deletes `7e6c6f63616c5f706c61796572` from LevelDB
