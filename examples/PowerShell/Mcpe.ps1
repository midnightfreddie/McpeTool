[cmdletbinding()]
param (
    [string]$EndPoint = "http://localhost:8080",
    [string]$ApiRoot = "/api/v1/db"
)

function ConvertTo-Base64 {
    param ( [byte[]]$ByteArray )
    [System.Convert]::ToBase64String($ByteArray)
}

function ConvertFrom-Base64 {
    param ( [string]$Base64String )
    [System.Text.Encoding]::UTF8.GetString(
        [System.Convert]::FromBase64String($Base64String)
    )
}

$Data = [System.Text.Encoding]::UTF8.GetBytes("Hello")
$HexKey = "010203"
$Uri = $EndPoint + $ApiRoot + "/" + $HexKey
$Body = New-Object psobject -Property @{
    base64Data = ConvertTo-Base64 $Data
} | ConvertTo-Json

$DeleteResult = Invoke-WebRequest -Uri $Uri -Method Delete

$PutResult = Invoke-WebRequest -Uri $Uri -Method Put -Body $Body

Remove-Variable GetResult
$GetResult = Invoke-WebRequest -Uri $Uri

$GetResult
