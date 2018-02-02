# LevelDB and key info

Minecraft uses LevelDB but doesn't use its standard Snappy compression, instead adding compression types 2 and 4 for zlib deflate—which they term ZlibCompressor in their modified LevelDB code—and flate deflate (ZlibCompressorRaw).

LevelDB keys (and data) are byte arrays, but hex strings are more human-friendly, so we use those to represent the keys.

## Chunk-related keys

Most of the game data is tied to chunks (16x16x256-block-sized bits of the world), so most keys consist of an X and Z chunk coordinate, then optionally a dimension coordinate (for Nether, End, etc.), then a tag that identifies the type of data, and then for block data a Y subchunk index. [See below](#how-to-convert-world-coordinates-to-leveldb-keys) for more details.

Decimal Tag | Hex Tag | Data Type | Data
---|---|---|---
45 | 2d | int16[] ? | Data2D
46 | 2e |  byte[] | Data2DLegacy
47 | 2f | nbt | SubChunkPrefix
48 | 30 |  byte[] | LegacyTerrain (pre-1.0 128 height worlds)
49 | 31 | nbt | BlockEntity
50 | 32 | nbt | Entity
51 | 33 |  nbt | PendingTicks
52 | 34 |  int32, int32, int16 ? | BlockExtraData
53 | 35 |  byte[3] ? | BiomeState
54 | 36 | int32 | FinalizedState
118 | 76 | byte | Version

## Non-chunk keys

Some game data is not tied to chunks, and those keys are simple text descriptions. Here are some hex string keys for known simple keys:

Simple Name | Hex Key | Format | Contents
---|---|---|---
~local_player | 7e6c6f63616c5f706c61796572 | nbt | Local player data including inventory
AutonomousEntities | 4175746f6e6f6d6f7573456e746974696573 | nbt |
BiomeData | 42696f6d6544617461 | nbt |
Overworld | 4f766572776f726c64 | nbt |
mVillages | 6d56696c6c61676573 | nbt |
portals | 706f7274616c73 | nbt |
player_ | 706c617965725f | nbt | Keys beginning with this represent multiplayer players

## How to convert world coordinates to leveldb keys

X, Y, Z is the typical coordinate order, but when dealing with the data we find X, Z, Y is the order of greatest-to-least significance, which is why this explanation tends to express X, Z, Y ordering.

All division below is of course integer division. The remainder/modulus will be used to find the byte offset within the subchunk data. X, Z, and dimension are 32-bit signed integers in little endian byte order. In the examples below, I've bolded the chunk Z coordinate for clarity.

Each chunk is 16x16x256 (X,Z,Y), and the subchunk block data keys are 16 high. So for x, z, y coordinates of 413, 54, 105:

- chunk X = 413 / 16 = 25 or 0x19 signed 32-bit integer in little endian byte order ([0x19,0, 0, 0] == 19000000)
- chunk Z = 54 / 16 = 3 ([0x3, 0, 0, 0] == **03000000**) 

So all keys beginning with 19000000**03000000** are about this coordinate's chunk. (In the overworld; other dimensions add a 32-bit dimension ID, so the same coordinates in the Nether I think have keys that start with 19000000**03000000**FFFFFFFF and 19000000**03000000**01000000 for the End.)

The tags and subchunk indexes are 8-bit values. (Unsigned? Not sure it matters as there are no negative Y chunk coordinates and no tags <0 or > 127.)

47 ([0x2F]) is the subchunk prefix tag, so all keys beginning with 19000000**03000000**2f are the Y subchunks for this coordinate.

- subchunk Y = 105 / 16 = 6 ([0x*06*])

So, the subchunk key for X=413, Z=54, Y=105 is 19000000**03000000**2f*06*

ref: <https://minecraft.gamepedia.com/Bedrock_Edition_level_format>

## How to convert world coordinates to block subchunk byte offsets

`%` is modulo division, the remainder of integer division. Note that subchunks may not exist if they are empty.

- `(X % 16) * 256 + (Z % 16) * 16 + (Y % 16) + 1` - 1-byte Block ID that determines if the block is e.g. air, water, sand, dirt, tree, lava, etc.
- The following three are odd in that each byte holds data for two blocks, in the high and low nybbles.
- `(X % 16) * 128 + (Z % 16) * 8 + (Y % 16) / 2 + 4096 + 1` - Block data e.g. rotation
- `(X % 16) * 128 + (Z % 16) * 8 + (Y % 16) / 2 + 6144 + 1` - Sky light
- `(X % 16) * 128 + (Z % 16) * 8 + (Y % 16) / 2 + 8192 + 1` - Block light
