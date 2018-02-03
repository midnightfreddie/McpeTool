### 2018-02-02

Way back in 0.9--128-high world, blocks stored in tag 0x30--I noticed that if I deleted chunks and placed chunks
with my "spiral staircase", the game would generate surrounding terrain and try to match the edges of my
manual chunks.

Today I'm trying to see if the behavior is similar.

Lessons learned:

- The terrain kept appearing to restore itself, and I finally figured out it wasn't restoring but regenerating with the same seed
- When deleting the world and placing new subchunks, the version tag key (0x76, e.g. 010000000100000076) must be present or the game discards my subchunk and generates a new one
- The game does **not** generate new subchunks underneath a placed chunk. e.g if I place 01000000010000002f04 and don't have 2f03, 2f02, 2f01, and 2f00, then the world is just empty air, and if I dig through the bottom of my placed chunk I fall all the way out of the world
- When placing a large blob of 2f04 blocks with only one layer of grass at 0 (64), the generated terrain does not seem to try to match the edges, although I blocks 0-63 are all air, so that might be contributing.
