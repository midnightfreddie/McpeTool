$EndPoint = "http://localhost:8080"
$ApiRoot = "/api/v1/db"

function Get-KeyByCoords {
    [cmdletbinding()]
    param (
        [Int32]$X,
        [Int32]$Z,
        [Int32]$Y,
        [Int32]$Dimension = 0,
        [byte]$Tag = 0x2f
    )
    [Int32]$ChunkX = $X / 16
    [Int32]$ChunkZ = $Z / 16
    [byte]$SubChunkY = $Y / 16
    return (
        [System.BitConverter]::GetBytes($ChunkX) +
        [System.BitConverter]::GetBytes($ChunkZ) +
        $(if ($Dimension) { [System.BitConverter]::GetBytes($Dimension) }) +
        $Tag +
        $(if ($Tag -eq 0x2f) { $SubChunkY }) |
        ForEach-Object {
            '{0:x2}' -f $PSItem
        }
    ) -join ''
}

function ConvertTo-Base64 {
    [cmdletbinding()]
    param ( [byte[]]$ByteArray )
    [System.Convert]::ToBase64String($ByteArray)
}

function ConvertFrom-Base64 {
    [cmdletbinding()]
    param ( [string]$Base64String )
    [System.Convert]::FromBase64String($Base64String)
}

# Simple API exercise. Put, get, then delete a key/value pair that has nothing to do with the game
function Invoke-McpeApiTest {
    [cmdletbinding()]
    param(
        $HexKey = "010203",
        $Data = [System.Text.Encoding]::ASCII.GetBytes("Hello")
    )
    $Uri = $EndPoint + $ApiRoot + "/" + $HexKey
    $Body = New-Object psobject -Property @{
        base64Data = ConvertTo-Base64 $Data
    } | ConvertTo-Json

    # Put the key/value combo in MCPE LevelDB
    $PutResult = Invoke-WebRequest -Uri $Uri -Method Put -Body $Body

    # Get that value
    $GetResult = Invoke-WebRequest -Uri $Uri

    # Convert the JSON/base64 value and output it. Most data isn't a printable string, but since this is a test I know it should be "Hello"
    Write-Output [System.Text.Encoding]::UTF8.GetString((ConvertFrom-Base64 (($GetResult.Content | ConvertFrom-Json).base64Data)))

    # Delete the key/value pair from MCPE LevelDB
    $DeleteResult = Invoke-WebRequest -Uri $Uri -Method Delete
}

# Returns byte offset of block in a terrain chunk.
# Can pass in chunk-relative or world coordinates because X % 16 is the same as (X % 16) % 16
function Get-TerrainChunkOffset {
    [cmdletbinding()]
    param (
        [int]$X,
        [int]$Z,
        [int]$Y
    )
    $Value = ($X % 16) * 256 + ($Z % 16) * 16 + ($Y % 16) + 1
    Write-Output $Value
}

# Loads chunk 0,0 terrain and writes a pillar in the middle and staircase around the perimeter
# X and Z are world coordinates, not chunk coordinates
# Y is irrelevant for this function
function Invoke-McpeSpiralStaircase {
    [cmdletbinding()]
    param (
        $X,
        $Z,
        $Y,
        $OddYBlockID = 42,
        $EvenYBlockID = 57
    )
    # Set up a loop to add the pillar and spiral
    $RelativeX = 0
    $RelativeZ = 0
    for ($Y = 0; $Y -lt 256; $Y++) {
        if ($Y % 16 -eq 0) {
            # Read subchunk
            $Uri = $EndPoint + $ApiRoot + "/" + (Get-KeyByCoords -X $X -Z $Z -Y $Y)
            $Result = Invoke-WebRequest -Uri $Uri -ErrorAction Stop
            # The response body is a JSON object with a "base64Data" key holding the base64-encoded data. Decode it to [byte[]]
            $ChunkData = ConvertFrom-Base64 (($Result.Content | ConvertFrom-Json).base64Data)
        }
        if ($Y % 2 -eq 0) {
            $BlockID = $EvenYBlockID
        } else {
            $BlockID = $OddYBlockID
        }
        # Place a pillar block
        $ChunkData[(Get-TerrainChunkOffset -X 7 -Z 7 -Y ($Y % 16) )] = $BlockID
        # Place a "spiral staircase" block
        $ChunkData[(Get-TerrainChunkOffset -X $RelativeX  -Z $RelativeZ -Y ($Y % 16) )] = $BlockID
        # Move staircase block X or Z
        switch ([Math]::Floor($Y / 15) %4) {
            0 { $RelativeX++; break }
            1 { $RelativeZ++; break }
            2 { $RelativeX--; break }
            3 { $RelativeZ--; break }
        }
        if ($Y % 16 -eq 15) {
            # Write subchunk
            # Create a JSON-encoded request body with the "base64Data" key holding the chunk data in base64 format
            $Body = New-Object psobject -Property @{
                base64Data = ConvertTo-Base64 $ChunkData
            } | ConvertTo-Json
            # Write the chunk back to the database
            $Result = Invoke-WebRequest -Uri $Uri -Method Put -Body $Body
        }        
    }
    # Remember to stop the API server before copying/playing the level again
}
