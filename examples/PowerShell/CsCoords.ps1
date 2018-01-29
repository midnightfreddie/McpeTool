$source = @"
using System;

public class McBedrockTool {
    public static string SubchunkKey(int x, int z, int y) {
        int ChunkX = x / 16;
        // return (BitConverter.GetBytes(ChunkX));
        // return String.Concat(Array.ConvertAll(BitConverter.GetBytes(ChunkX), e => e.ToString("X2")));
        // return ChunkX;
        return HexKey(ChunkX);
    }
    public static string HexKey(int i) {
        return String.Concat(Array.ConvertAll(BitConverter.GetBytes(i), x => x.ToString("X2")));
    }
}
"@

Add-Type -TypeDefinition $source

[McBedrockTool]::SubchunkKey(413,54,105)

pause
