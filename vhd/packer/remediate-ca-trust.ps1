# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT license.

[CmdletBinding()]
param(
    [switch]$VerifyOnly,
    [string]$ProfilePath = (Join-Path $PSScriptRoot 'ca-remediation.json'),
    [string]$CtlDirectory = (Join-Path $PSScriptRoot 'ca-trust-ctl'),
    [string]$VerificationLogPath
)

$ErrorActionPreference = 'Stop'

function Read-CaRemediationProfile {
    param([string]$Path)

    $profile = Get-Content -LiteralPath $Path -Raw -ErrorAction Stop | ConvertFrom-Json -ErrorAction Stop
    if (($profile.schemaVersion -isnot [int] -and $profile.schemaVersion -isnot [long]) -or $profile.schemaVersion -ne 1) {
        throw 'Unsupported CA remediation profile schemaVersion.'
    }
    if ($profile.incident -isnot [string] -or $profile.incident -cnotmatch '\A[A-Za-z0-9][A-Za-z0-9 ._-]{0,99}\z') {
        throw 'Invalid CA remediation incident label.'
    }
    if ($profile.windows.expectedDisallowedCtlSequence -isnot [string] -or
        $profile.windows.expectedDisallowedCtlSequence -cnotmatch '\A[0-9a-fA-F]{16}\z') {
        throw 'The CA remediation profile requires a reviewed 16-digit Windows Disallowed CTL sequence.'
    }
    return $profile
}

function Invoke-CaCertutil {
    param(
        [string[]]$ArgumentList,
        [switch]$ExpectUnlistedProbe
    )

    Write-Host "certutil.exe $($ArgumentList -join ' ')"
    $output = & certutil.exe @ArgumentList
    $exitCode = $LASTEXITCODE
    $output | ForEach-Object { Write-Host $_ }
    if ($ExpectUnlistedProbe -and $exitCode -eq -2146877435) {
        Write-Host 'The generated probe is not in the Disallowed CTL (0x80094005 CERTSRV_E_INVALID_CA_CERTIFICATE), as expected. Persisted CTL verification is still required.'
        return
    }
    if ($exitCode -ne 0) {
        throw "CA CTL operation failed: certutil.exe exited with $exitCode."
    }
}

function Get-CaEncodedCtl {
    $path = 'HKLM:\SOFTWARE\Microsoft\SystemCertificates\AuthRoot\AutoUpdate'
    $encoded = Get-ItemPropertyValue -LiteralPath $path -Name DisallowedCertEncodedCtl -ErrorAction Stop
    if ($encoded -isnot [byte[]] -or $encoded.Length -eq 0) {
        throw 'The persisted Windows Disallowed CTL is missing or is not binary data.'
    }
    return ,$encoded
}

function Get-CaCtlSequence {
    param([byte[]]$EncodedCtl)

    if (-not ('AKSe.CaTrustCtl' -as [type])) {
        Add-Type -TypeDefinition @'
using System;
using System.ComponentModel;
using System.Runtime.InteropServices;

namespace AKSe {
    public static class CaTrustCtl {
        [StructLayout(LayoutKind.Sequential)]
        private struct Blob {
            public uint Length;
            public IntPtr Data;
        }

        [StructLayout(LayoutKind.Sequential)]
        private struct Usage {
            public uint Count;
            public IntPtr Identifiers;
        }

        [StructLayout(LayoutKind.Sequential)]
        private struct CtlInfoHeader {
            public uint Version;
            public Usage SubjectUsage;
            public Blob ListIdentifier;
            public Blob SequenceNumber;
        }

        [StructLayout(LayoutKind.Sequential)]
        private struct CtlContextHeader {
            public uint Encoding;
            public IntPtr Encoded;
            public uint EncodedLength;
            public IntPtr Info;
        }

        [DllImport("crypt32.dll", SetLastError = true)]
        private static extern IntPtr CertCreateCTLContext(uint encoding, byte[] encoded, uint length);

        [DllImport("crypt32.dll")]
        [return: MarshalAs(UnmanagedType.Bool)]
        private static extern bool CertFreeCTLContext(IntPtr context);

        public static string GetSequence(byte[] encoded) {
            IntPtr context = CertCreateCTLContext(0x00010001, encoded, (uint)encoded.Length);
            if (context == IntPtr.Zero) {
                throw new Win32Exception(Marshal.GetLastWin32Error(), "Cannot decode the persisted Disallowed CTL.");
            }
            try {
                CtlContextHeader header = (CtlContextHeader)Marshal.PtrToStructure(context, typeof(CtlContextHeader));
                CtlInfoHeader info = (CtlInfoHeader)Marshal.PtrToStructure(header.Info, typeof(CtlInfoHeader));
                if (info.SequenceNumber.Length != 8 || info.SequenceNumber.Data == IntPtr.Zero) {
                    throw new InvalidOperationException("The Disallowed CTL has no supported 8-byte sequence number.");
                }
                byte[] sequence = new byte[8];
                Marshal.Copy(info.SequenceNumber.Data, sequence, 0, sequence.Length);
                Array.Reverse(sequence);
                return BitConverter.ToString(sequence).Replace("-", "").ToLowerInvariant();
            }
            finally {
                CertFreeCTLContext(context);
            }
        }
    }
}
'@
    }
    return [AKSe.CaTrustCtl]::GetSequence($EncodedCtl)
}

function Update-CaCtl {
    param([string]$CtlDirectory)

    New-Item -ItemType Directory -Path $CtlDirectory -Force -ErrorAction Stop | Out-Null
    Invoke-CaCertutil -ArgumentList @('-f', '-syncWithWU', $CtlDirectory)
    if (-not (Test-Path -LiteralPath (Join-Path $CtlDirectory 'disallowedcert.sst') -PathType Leaf)) {
        throw 'Windows Update synchronization did not produce disallowedcert.sst for CTL certificate matching.'
    }

    $directory = Join-Path $env:TEMP ('akse-ca-trust-' + [guid]::NewGuid().ToString('N'))
    $probePath = Join-Path $directory 'untrusted-probe.cer'
    $rsa = $null
    $certificate = $null
    New-Item -ItemType Directory -Path $directory -ErrorAction Stop | Out-Null
    try {
        # Microsoft documents an untrusted CertFile as the forced CTL refresh trigger.
        # The probe is generated in memory and is never added to a certificate store.
        $rsa = [System.Security.Cryptography.RSACng]::new(2048)
        $request = [System.Security.Cryptography.X509Certificates.CertificateRequest]::new(
            'CN=AKSe CTL refresh probe',
            $rsa,
            [System.Security.Cryptography.HashAlgorithmName]::SHA256,
            [System.Security.Cryptography.RSASignaturePadding]::Pkcs1
        )
        $certificate = $request.CreateSelfSigned(
            [DateTimeOffset]::UtcNow.AddMinutes(-5),
            [DateTimeOffset]::UtcNow.AddDays(1)
        )
        [IO.File]::WriteAllBytes($probePath, $certificate.Export(
            [System.Security.Cryptography.X509Certificates.X509ContentType]::Cert
        ))
        Invoke-CaCertutil -ArgumentList @('-f', '-verifyCTL', 'Disallowed', $CtlDirectory, $probePath) -ExpectUnlistedProbe
    }
    finally {
        if ($certificate) { $certificate.Dispose() }
        if ($rsa) { $rsa.Dispose() }
        Remove-Item -LiteralPath $directory -Recurse -Force -ErrorAction Stop
    }
}

function Assert-CaCtl {
    param(
        [string]$ExpectedSequence,
        [string]$Incident,
        [string]$CtlDirectory
    )

    if (-not (Test-Path -LiteralPath (Join-Path $CtlDirectory 'disallowedcert.sst') -PathType Leaf)) {
        throw 'The synchronized disallowedcert.sst is missing. Run image remediation before verification-only mode.'
    }
    Invoke-CaCertutil -ArgumentList @('-verifyCTL', 'Disallowed', $CtlDirectory)
    $encoded = Get-CaEncodedCtl
    $sequence = Get-CaCtlSequence -EncodedCtl $encoded
    if ($sequence -ne $ExpectedSequence) {
        throw "Disallowed CTL sequence '$sequence' does not match the reviewed CTL '$ExpectedSequence' for '$Incident'. Review Microsoft's mitigation coverage before accepting another sequence."
    }
    return "Verified persisted Windows Disallowed CTL: SequenceNumber=$sequence; reviewed incident=$Incident."
}

function Invoke-CaRemediation {
    param(
        [switch]$VerifyOnly,
        [string]$ProfilePath,
        [string]$CtlDirectory
    )

    $profile = Read-CaRemediationProfile -Path $ProfilePath
    if (-not $VerifyOnly) {
        Update-CaCtl -CtlDirectory $CtlDirectory
    }
    return Assert-CaCtl -ExpectedSequence $profile.windows.expectedDisallowedCtlSequence -Incident $profile.incident -CtlDirectory $CtlDirectory
}

if ($MyInvocation.InvocationName -ne '.') {
    $result = Invoke-CaRemediation -VerifyOnly:$VerifyOnly -ProfilePath $ProfilePath -CtlDirectory $CtlDirectory
    if ($VerificationLogPath) {
        $result | Set-Content -LiteralPath $VerificationLogPath -Encoding UTF8 -ErrorAction Stop
    }
    Write-Output $result
}
