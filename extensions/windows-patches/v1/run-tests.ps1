[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'

if ($PSVersionTable.PSEdition -eq 'Desktop') {
    $pwsh = Get-Command pwsh.exe -ErrorAction SilentlyContinue
    if ($null -eq $pwsh) {
        throw "PowerShell 7 is required to bootstrap the local test dependencies. Install it and run this script again."
    }

    & $pwsh.Source -NoProfile -File $PSCommandPath
    exit $LASTEXITCODE
}

$requiredPesterVersion = [Version]'5.7.1'
$pester = Get-Module -ListAvailable Pester |
    Where-Object { $_.Version -eq $requiredPesterVersion } |
    Select-Object -First 1

if ($null -eq $pester) {
    Write-Host "Pester $requiredPesterVersion is not installed. Installing it for the current user..."
    Install-Module Pester -Scope CurrentUser -RequiredVersion $requiredPesterVersion -Force -SkipPublisherCheck
    $pester = Get-Module -ListAvailable Pester |
        Where-Object { $_.Version -eq $requiredPesterVersion } |
        Select-Object -First 1
}

Remove-Module Pester -Force -ErrorAction SilentlyContinue
Import-Module $pester.Path -Force

$configuration = New-PesterConfiguration
$configuration.Run.Path = Join-Path $PSScriptRoot 'installPatches.tests.ps1'
$configuration.Run.Exit = $true
$configuration.Output.Verbosity = 'Detailed'

Invoke-Pester -Configuration $configuration