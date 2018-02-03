<#
  An attempt to try to make a randomly-generated flat world

  I've noticed that when deleting world chunks and creating standalone isolated chunks,
  the in-game terrain generator tries to match the edges of the isolated chunk. But
  this was in 0.9 with 127-high worlds.

  So, in Bedrock, if I create a subchunks 05 that just have one layer of dirt at y=65,
  will the generator fill in below and around making a flat survival world?
#>

# $yStack = New-Object byte[] 16
# $yStack[0] = [byte]2

# Make a zero-filled byte array the size of a subchunk
# $subChunk = New-Object byte[] (1 + 4096 + 2048 + 2048 + 2048)
# It seems to be okay to leave off sky light and block light; at least generated chunks don't store them
$subChunk = New-Object byte[] (1 + 4096 + 2048)
# Fill subchunk with stone
for ($x = 1; $x -lt 4097; $x++) { 
    $subChunk[$x] = [byte]1
}
$TopSubChunk = $subChunk | ForEach-Object { $PSItem }
$BottomSubChunk = $subChunk | ForEach-Object { $PSItem }
# Put grass/dirt blocks in the top 3 layers
for ($x = 0; $x -lt 16; $x++) { 
    for ($z = 0; $z -lt 16; $z++) { 
        $TopSubChunk[1 + ($x * 16 * 16) + ($z * 16) + 12] = [byte]3
        $TopSubChunk[1 + ($x * 16 * 16) + ($z * 16) + 13] = [byte]3
        $TopSubChunk[1 + ($x * 16 * 16) + ($z * 16) + 14] = [byte]2
        $TopSubChunk[1 + ($x * 16 * 16) + ($z * 16) + 15] = [byte]0
    }    
}
# Put bedrock blocks in bottom layer
for ($x = 0; $x -lt 16; $x++) { 
    for ($z = 0; $z -lt 16; $z++) { 
        $BottomSubChunk[1 + ($x * 16 * 16) + ($z * 16)] = [byte]7
    }    
}

# DELETE ALL THE THINGS
$KeyList = Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/"
$KeyList.keys | ForEach-Object {
    Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/$($PSItem.hexKey)" -Method Delete
}

$SubchunkBody = New-Object psobject -Property @{
    base64Data = [System.Convert]::ToBase64String($subChunk)
} | ConvertTo-Json

$TopSubchunkBody = New-Object psobject -Property @{
    base64Data = [System.Convert]::ToBase64String($TopSubChunk)
} | ConvertTo-Json

$BottomSubchunkBody = New-Object psobject -Property @{
    base64Data = [System.Convert]::ToBase64String($BottomSubChunk)
} | ConvertTo-Json

$ChunkVersionBody = New-Object psobject -Property @{
    base64Data = [System.Convert]::ToBase64String([byte]7)
} | ConvertTo-Json

# Manual spawn point until I get the leveldat API implemented
[Int32]$SpawnX = 92
[Int32]$SpawnZ = 20
[Int32]$SpawnY = 32767

[int32]$ChunkX = [Math]::Floor($SpawnX / 16)
[int32]$ChunkZ = [Math]::Floor($SpawnZ / 16)

for ($x = $ChunkX - 6; $x -lt $ChunkX + 6; $x+=2) { 
    for ($z = $ChunkZ - 6; $z -lt $ChunkZ + 6; $z+=2) { 
        $HexPrefix = (
            @([System.BitConverter]::GetBytes([int32]$x)) +
            [System.BitConverter]::GetBytes([int32]$z) |
            ForEach-Object {
                '{0:x2}' -f $PSItem
            }
            ) -join ''
        # $HexPrefix
        Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/${HexPrefix}2f00" -Method Put -Body $BottomSubchunkBody
        Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/${HexPrefix}2f01" -Method Put -Body $SubchunkBody
        Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/${HexPrefix}2f02" -Method Put -Body $SubchunkBody
        Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/${HexPrefix}2f03" -Method Put -Body $TopSubchunkBody
        Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/${HexPrefix}76" -Method Put -Body $ChunkVersionBody
        # Meh, for starters just manually put one subchunk:
        # Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/05000000010000002f05" -Method Put -Body $Body
        # Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/05000000010000002f04" -Method Put -Body $Body
        # Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/05000000010000002f03" -Method Put -Body $Body
        # Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/05000000010000002f02" -Method Put -Body $Body
        # Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/05000000010000002f01" -Method Put -Body $Body
        # Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/v1/db/05000000010000002f00" -Method Put -Body $Body
    }
}