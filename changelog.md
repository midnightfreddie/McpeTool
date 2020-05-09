## v0.3.2

The API can now read and write nbt2json-formatted nbt data if the db key is all nbt.
Also updates for newer dependencies and a fix for nbt long values getting messed up.

- Fixed CsCoords.ps1; it was miscalculating negative coordinates
- Added go.mod
- Moved mcpetool/ to cmd/mcpetool/
- Merged readme pull request
- Updated to urfave/cli/v2
- Started this changelog (previous version notes are from memory or release page)
- **Breaking change**: upgraded to nbt2json v0.4.0. NBT long values (int64) are now represented as strings instead of numbers because many JSON libraries won't preserve a 64-bit integer value.
- Http address and port are now configurable
- Aded GitHub link to help page
- Added CORS header to API to allow all origins
- Added `?json` parameter to api to allow GETs to provide nbt in json format if the entire chunk is nbt, and PUT will use the json instead of base64 when PUTting with `?json` (I can't believe I hadn't done this before!)

## v0.3.0

- level.dat - Can now get and put level.dat as binary or JSON/YAML
- Commands refactored - Since more features are coming, moved `list`, `get`, `put` and `delete` under `db` command. Added `leveldat` command for level.dat `get`s and `put`s.
- Updated documentation for the new commands and include how to convert world coordinates to a db key

## v0.2.1

- Now can use YAML or JSON. YAML is translated to/from JSON for NBT conversion
- JSON is reformatted
  - There is now a top-level document object with version and other info
  - The NBT tag list is always an array now
- Placing tool, world path, and key info in the new upstream JSON comment field

The JSON/YAML produced and consumed by this version is not compatible with pre-0.2.0 versions.

## v0.1.3

[/u/Flaming5asquatch](https://www.reddit.com/user/Flaming5asquatch) discovered that some new save files don't work with the older version. This release now works with the newer save files.

The path to the world is no longer provided as an argument but can be provided one of the following ways:

- With the `--path` parameter
- With the `MCPETOOL_WORD` environment variable
- Defaults to current directory (".")

Also, the command-line `get` and `put` actions can experimentally output and input in JSON as converted by my [nbt2json](https://github.com/midnightfreddie/nbt2json). Use the `--json` option with `get` or `put`. This is not yet heavily tested, but I have been able to, for example, change item counts in my inventory, and add inventory items to chests or player by outputting the player key or world block, editing the JSON, then putting the updated value back. As always, **back up your worlds before making changes**.

## v0.1.1 First release

I think this has enough features to be useful to others, so I am releasing as v0.1.1 . (v0.1.0 was mis-tagged and removed)
- Get, put, and delete access to key/value pairs in MCPE LevelDB world via command line or http API
- Option to hex dump values
- Others should be able to use this as a basic access tool for more intelligent and/or user-friendly world editing
