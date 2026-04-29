# RTDB API publish script (PowerShell)
param(
    [switch]$d,
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$VersionArgs
)

$version = $VersionArgs[0]

if (-not $version) {
    Write-Host "Usage:"
    Write-Host "  Publish: .\publish.ps1 <version>"
    Write-Host "  Delete:  .\publish.ps1 -d <version>"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\publish.ps1 v4.0.15_0.2.0"
    Write-Host "  .\publish.ps1 -d v4.0.15_0.2.0"
    exit 1
}

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location $scriptDir

if ($d) {
    Write-Host "========== Delete version: $version =========="
    go run tools\publish_version.go -d "$version"
} else {
    Write-Host "========== Publish version: $version =========="
    go run tools\publish_version.go "$version"
}
