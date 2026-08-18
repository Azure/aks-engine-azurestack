
# Return codes:
#  0 - success
#  1 - install failure
#  2 - download failure
#  3 - unrecognized patch extension
#  4 - patch validation failure

param(
    [string[]] $URIs,

    [string[]] $SHA256Hashes,

    [switch] $EnableTestSigning
)

function New-PatchException([string] $Message, [int] $ExitCode)
{
    $exception = [InvalidOperationException]::new($Message)
    $exception.Data["ExitCode"] = $ExitCode
    return $exception
}

function DownloadAndVerifyFile([Uri] $URI, [string] $FullName, [string] $ExpectedSHA256)
{
    try {
        Write-Host "Downloading $($URI.GetLeftPart([UriPartial]::Path))"
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -UseBasicParsing -Uri $URI.AbsoluteUri -OutFile $FullName
    } catch {
        Remove-Item -LiteralPath $FullName -Force -ErrorAction SilentlyContinue
        throw (New-PatchException -Message $_.Exception.Message -ExitCode 2)
    }

    $actualSHA256 = (Get-FileHash -LiteralPath $FullName -Algorithm SHA256).Hash
    if (-not [string]::Equals($actualSHA256, $ExpectedSHA256, [StringComparison]::OrdinalIgnoreCase)) {
        Remove-Item -LiteralPath $FullName -Force -ErrorAction SilentlyContinue
        throw (New-PatchException -Message "SHA-256 verification failed for $($URI.GetLeftPart([UriPartial]::Path))" -ExitCode 4)
    }

    Write-Host "Verified SHA-256 for $([IO.Path]::GetFileName($FullName))"
}

function Get-PatchDefinitions([string[]] $URIs, [string[]] $SHA256Hashes)
{
    if ($null -eq $URIs -or $URIs.Count -eq 0 -or $null -eq $SHA256Hashes -or $SHA256Hashes.Count -eq 0) {
        throw (New-PatchException -Message "Patch URIs and SHA-256 hashes are required" -ExitCode 4)
    }

    if ($URIs.Count -ne $SHA256Hashes.Count) {
        throw (New-PatchException -Message "Each patch URI must have a corresponding SHA-256 hash" -ExitCode 4)
    }

    $patches = @()
    for ($index = 0; $index -lt $URIs.Count; $index++) {
        if ($SHA256Hashes[$index] -notmatch '^[A-Fa-f0-9]{64}$') {
            throw (New-PatchException -Message "Each SHA-256 hash must contain 64 hexadecimal characters" -ExitCode 4)
        }

        $parsedURI = $null
        if (-not [Uri]::TryCreate($URIs[$index], [UriKind]::Absolute, [ref] $parsedURI)) {
            throw (New-PatchException -Message "Patch URI must be an absolute URI" -ExitCode 4)
        }

        if ($parsedURI.Scheme -ine [Uri]::UriSchemeHttps) {
            throw (New-PatchException -Message "Patch URI must use HTTPS: $($parsedURI.GetLeftPart([UriPartial]::Path))" -ExitCode 4)
        }

        $fileName = [IO.Path]::GetFileName($parsedURI.LocalPath)
        if ([string]::IsNullOrWhiteSpace($fileName)) {
            throw (New-PatchException -Message "Patch URI must identify a file" -ExitCode 4)
        }

        $extension = [IO.Path]::GetExtension($fileName).ToLowerInvariant()
        if ($extension -ne ".exe" -and $extension -ne ".msu") {
            throw (New-PatchException -Message "This script extension doesn't know how to install $extension files" -ExitCode 3)
        }

        $patches += [PSCustomObject]@{
            URI = $parsedURI
            ExpectedSHA256 = $SHA256Hashes[$index]
            FileName = $fileName
            Extension = $extension
            FullName = [IO.Path]::Combine($env:TEMP, $fileName)
        }
    }

    return $patches
}

function Wait-PatchProcess($Process)
{
    Wait-Process -InputObject $Process
    return $Process.ExitCode
}

function Install-Patches([string[]] $URIs, [string[]] $SHA256Hashes, [switch] $EnableTestSigning)
{
    $patches = @(Get-PatchDefinitions -URIs $URIs -SHA256Hashes $SHA256Hashes)
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
                        throw (New-PatchException -Message "Failed to enable Windows test-signing mode, exitcode $($bcdedit.ExitCode)" -ExitCode 1)
                    }
                    $testSigningEnabled = $true
                }

                Write-Host "Starting $($patch.FullName)"
                $proc = Start-Process -Passthru -FilePath $patch.FullName -ArgumentList "/q /norestart"
                $exitCode = Wait-PatchProcess -Process $proc
                switch ($exitCode)
                {
                    0 {
                        Write-Host "Finished running $($patch.FullName)"
                    }
                    3010 {
                        Write-Host "Finished running $($patch.FullName). Reboot required to finish patching."
                    }
                    Default {
                        throw (New-PatchException -Message "Error running $($patch.FullName), exitcode $exitCode" -ExitCode 1)
                    }
                }
            }
            ".msu" {
                Write-Host "Installing $($patch.FullName)"
                $proc = Start-Process -Passthru -FilePath wusa.exe -ArgumentList "`"$($patch.FullName)`" /quiet /norestart"
                $exitCode = Wait-PatchProcess -Process $proc
                switch ($exitCode)
                {
                    0 {
                        Write-Host "Finished running $($patch.FullName)"
                    }
                    3010 {
                        Write-Host "Finished running $($patch.FullName). Reboot required to finish patching."
                    }
                    Default {
                        throw (New-PatchException -Message "Error running $($patch.FullName), exitcode $exitCode" -ExitCode 1)
                    }
                }
            }
        }
    }

    # No failures, schedule reboot now
    schtasks /create /TN RebootAfterPatch /RU SYSTEM /TR "shutdown.exe /r /t 0 /d 2:17" /SC ONCE /ST $(([System.DateTime]::Now + [timespan]::FromMinutes(5)).ToString("HH:mm")) /V1 /Z
}

if ($MyInvocation.InvocationName -ne '.') {
    try {
        Install-Patches -URIs $URIs -SHA256Hashes $SHA256Hashes -EnableTestSigning:$EnableTestSigning
        exit 0
    } catch {
        Write-Error $_.Exception.Message
        $exitCode = $_.Exception.Data["ExitCode"]
        if ($null -eq $exitCode) {
            $exitCode = 1
        }

        exit $exitCode
    }
}