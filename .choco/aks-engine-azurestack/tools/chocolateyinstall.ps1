﻿$ErrorActionPreference = 'Stop'

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
  packageName   = 'aks-engine-azurestack'
  url64bit      = '_placeholder_'
  checksum64    = '_placeholder_'
  checksumType64= 'sha256'
  unzipLocation = $toolsDir
}

Install-ChocolateyZipPackage @packageArgs