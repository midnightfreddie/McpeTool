[cmdletbinding()]
param (
    [string]$EndPoint = "http://localhost:8080",
    [string]$ApiRoot = "/api/v1/db"
)

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

# Returns byte offset of block in a terrain chunk. Does not validate input; X and Z should be from 0-15 each and Y from 0-127
function Get-TerrainChunkOffset {
    [cmdletbinding()]
    param (
        [int]$RelativeX,
        [int]$RelativeZ,
        [int]$Y
    )
    $Value = $RelativeX * 2048 + $RelativeZ * 128 + $Y
    Write-Output $Value
}

# Loads chunk 0,0 terrain and writes a pillar in the middle and staircase around the perimeter
function Invoke-McpeSpiralStaircase {
    [cmdletbinding()]
    param (
        $OddYBlockID = 42,
        $EvenYBlockID = 57,
        $HexKey = "000000000000000030"
    )
    $Uri = $EndPoint + $ApiRoot + "/" + $HexKey
    $Result = Invoke-WebRequest -Uri $Uri -ErrorAction Stop
    $ChunkData = ConvertFrom-Base64 (($Result.Content | ConvertFrom-Json).base64Data)
    $RelativeX = 0
    $RelativeZ = 0
    for ($Y = 0; $Y -lt 128; $Y++) {
        if ($Y % 2 -eq 0) {
            $BlockID = $EvenYBlockID
        } else {
            $BlockID = $OddYBlockID
        }
        # Place a pillar block
        $ChunkData[(Get-TerrainChunkOffset -RelativeX 7 -Relativez 7 -Y $Y )] = $BlockID
        # Place a "spiral staircase" block
        $ChunkData[(Get-TerrainChunkOffset -RelativeX $RelativeX  -RelativeZ $RelativeZ -Y $Y )] = $BlockID
        # Move staircase block X or Z
        switch ([Math]::Floor($Y / 15) %4) {
            0 { $RelativeX++; break }
            1 { $RelativeZ++; break }
            2 { $RelativeX--; break }
            3 { $RelativeZ--; break }
        }
    }
    $Body = New-Object psobject -Property @{
        base64Data = ConvertTo-Base64 $ChunkData
    } | ConvertTo-Json
    $Result = Invoke-WebRequest -Uri $Uri -Method Put -Body $Body
}

