
# Return codes:
#  0 - success
#  1 - install failure
#  2 - download failure
#  3 - unrecognized patch extension
#  4 - patch validation failure

param(
    [Parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [string[]] $URIs,

    [Parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [ValidatePattern('^[A-Fa-f0-9]{64}$')]
    [string[]] $SHA256Hashes,

    [switch] $EnableTestSigning
)

function DownloadAndVerifyFile([Uri] $URI, [string] $FullName, [string] $ExpectedSHA256)
{
    try {
        Write-Host "Downloading $($URI.GetLeftPart([UriPartial]::Path))"
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -UseBasicParsing -Uri $URI.AbsoluteUri -OutFile $FullName
    } catch {
        Remove-Item -LiteralPath $FullName -Force -ErrorAction SilentlyContinue
        Write-Error $_
        exit 2
    }

    $actualSHA256 = (Get-FileHash -LiteralPath $FullName -Algorithm SHA256).Hash
    if (-not [string]::Equals($actualSHA256, $ExpectedSHA256, [StringComparison]::OrdinalIgnoreCase)) {
        Remove-Item -LiteralPath $FullName -Force -ErrorAction SilentlyContinue
        Write-Error "SHA-256 verification failed for $($URI.GetLeftPart([UriPartial]::Path))"
        exit 4
    }

    Write-Host "Verified SHA-256 for $([IO.Path]::GetFileName($FullName))"
}

if ($URIs.Count -ne $SHA256Hashes.Count) {
    Write-Error "Each patch URI must have a corresponding SHA-256 hash"
    exit 4
}

$patches = @()
for ($index = 0; $index -lt $URIs.Count; $index++) {
    $parsedURI = $null
    if (-not [Uri]::TryCreate($URIs[$index], [UriKind]::Absolute, [ref] $parsedURI)) {
        Write-Error "Patch URI must be an absolute URI"
        exit 4
    }

    if ($parsedURI.Scheme -ine [Uri]::UriSchemeHttps) {
        Write-Error "Patch URI must use HTTPS: $($parsedURI.GetLeftPart([UriPartial]::Path))"
        exit 4
    }

    $fileName = [IO.Path]::GetFileName($parsedURI.LocalPath)
    if ([string]::IsNullOrWhiteSpace($fileName)) {
        Write-Error "Patch URI must identify a file"
        exit 4
    }

    $extension = [IO.Path]::GetExtension($fileName).ToLowerInvariant()
    if ($extension -ne ".exe" -and $extension -ne ".msu") {
        Write-Error "This script extension doesn't know how to install $extension files"
        exit 3
    }

    $patches += [PSCustomObject]@{
        URI = $parsedURI
        ExpectedSHA256 = $SHA256Hashes[$index]
        FileName = $fileName
        Extension = $extension
        FullName = [IO.Path]::Combine($env:TEMP, $fileName)
    }
}

$testSigningEnabled = $false
$patches | ForEach-Object {
    $patch = $_
    Write-Host "Processing $($patch.URI.GetLeftPart([UriPartial]::Path))"
    DownloadAndVerifyFile -URI $patch.URI -FullName $patch.FullName -ExpectedSHA256 $patch.ExpectedSHA256

    switch ($patch.Extension) {
        ".exe" {
            if ($EnableTestSigning -and -not $testSigningEnabled) {
                Write-Warning "Enabling Windows test-signing mode for executable patches"
                $bcdedit = Start-Process -Passthru -Wait -FilePath bcdedit.exe -ArgumentList "/set {current} testsigning on"
                if ($bcdedit.ExitCode -ne 0) {
                    Write-Error "Failed to enable Windows test-signing mode, exitcode $($bcdedit.ExitCode)"
                    exit 1
                }
                $testSigningEnabled = $true
            }

            Write-Host "Starting $($patch.FullName)"
            $proc = Start-Process -Passthru -FilePath $patch.FullName -ArgumentList "/q /norestart"
            Wait-Process -InputObject $proc
            switch ($proc.ExitCode)
            {
                0 {
                    Write-Host "Finished running $($patch.FullName)"
                }
                3010 {
                    Write-Host "Finished running $($patch.FullName). Reboot required to finish patching."
                }
                Default {
                    Write-Error "Error running $($patch.FullName), exitcode $($proc.ExitCode)"
                    exit 1
                }
            }
        }
        ".msu" {
            Write-Host "Installing $($patch.FullName)"
            $proc = Start-Process -Passthru -FilePath wusa.exe -ArgumentList "`"$($patch.FullName)`" /quiet /norestart"
            Wait-Process -InputObject $proc
            switch ($proc.ExitCode)
            {
                0 {
                    Write-Host "Finished running $($patch.FullName)"
                }
                3010 {
                    Write-Host "Finished running $($patch.FullName). Reboot required to finish patching."
                }
                Default {
                    Write-Error "Error running $($patch.FullName), exitcode $($proc.ExitCode)"
                    exit 1
                }
            }
        }
    }
}

# No failures, schedule reboot now

schtasks /create /TN RebootAfterPatch /RU SYSTEM /TR "shutdown.exe /r /t 0 /d 2:17" /SC ONCE /ST $(([System.DateTime]::Now + [timespan]::FromMinutes(5)).ToString("HH:mm")) /V1 /Z
exit 0