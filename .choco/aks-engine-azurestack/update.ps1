Import-Module au -Force

function global:au_SearchReplace {
  Write-Host "au_SearchReplace"
  @{
    'tools\chocolateyInstall.ps1' = @{
      "(^\s*url64bit\s*=\s*)('.*')"   = "`$1'$($Latest.URL64)'"
      "(^\s*checksum64\s*=\s*)('.*')" = "`$1'$($Latest.Checksum64)'"
    }
  }
}

function global:au_GetLatest {
  $ProgressPreference = 'SilentlyContinue'
  $hash_check_file_path = "$pwd/aksengine.hashcheck"
  [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12;

  $url = "https://github.com/Azure/aks-engine-azurestack/releases/download/v${ReleaseVersion}/aks-engine-azurestack-v${ReleaseVersion}-windows-amd64.zip"
  $wc = New-Object net.webclient
  $wc.Downloadfile($url, $hash_check_file_path)
  
  $checksum = Get-FileHash -Algorithm SHA256 -Path $hash_check_file_path
  $checksum64 = $checksum.Hash.ToLower()

  $Latest = @{ URL64 = "$url"; Version = $ReleaseVersion; Checksum64 = $checksum64 }
  return $Latest
}

Update-Package -ChecksumFor none -Verbose
