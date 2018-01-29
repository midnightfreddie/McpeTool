$source = @"
using System;

public class McBedrockTool {
    public static string GetKeyByCoords(int x, int z, int y = 0, int Dimension = 0, byte Tag = 0x2f) {
        int ChunkX = x / 16;
        int ChunkZ = z / 16;
        byte SubChunkY = (byte) (y / 16);
        string MyKey = HexKey(ChunkX) +
            HexKey(ChunkZ) +
            (Dimension == 0 ? "" : HexKey(Dimension)) +
            HexKey(Tag) +
            (Tag == 0x2f ? HexKey(SubChunkY) : "");
        return MyKey;
    }
    public static string HexKey(int i) {
        byte[] ByteArray = BitConverter.GetBytes(i);
        // Force output order to little endian no matter the local platform
        if (! BitConverter.IsLittleEndian)
            Array.Reverse(ByteArray);
        return String.Concat(Array.ConvertAll(ByteArray, x => x.ToString("X2")));
    }
    public static string HexKey(byte i) {
        return i.ToString("X2");
    }
}
"@

Add-Type -TypeDefinition $source

[McBedrockTool]::GetKeyByCoords(413,54,90)

pause
