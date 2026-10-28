# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT license.

[CmdletBinding()]
param(
    [switch]$VerifyOnly,
    [ValidatePattern('^[0-9a-fA-F]{16}$')]
    [string]$ExpectedSequenceNumber = '01dd5f6ca45c9d86',
    [string]$VerificationLogPath
)

$ErrorActionPreference = 'Stop'

function Invoke-EmudhraCertutil {
    param([string[]]$ArgumentList)

    Write-Host "certutil.exe $($ArgumentList -join ' ')"
    $output = & certutil.exe @ArgumentList
    $exitCode = $LASTEXITCODE
    $output | ForEach-Object { Write-Host $_ }
    if ($exitCode -ne 0) {
        throw "eMudhra CTL operation failed: certutil.exe exited with $exitCode."
    }
}

function Get-EmudhraEncodedCtl {
    $path = 'HKLM:\SOFTWARE\Microsoft\SystemCertificates\AuthRoot\AutoUpdate'
    $encoded = Get-ItemPropertyValue -LiteralPath $path -Name DisallowedCertEncodedCtl -ErrorAction Stop
    if ($encoded -isnot [byte[]] -or $encoded.Length -eq 0) {
        throw 'The persisted Windows Disallowed CTL is missing or is not binary data.'
    }
    return ,$encoded
}

function Get-EmudhraCtlSequence {
    param([byte[]]$EncodedCtl)

    if (-not ('AKSe.EmudhraCtl' -as [type])) {
        Add-Type -TypeDefinition @'
using System;
using System.ComponentModel;
using System.Runtime.InteropServices;

namespace AKSe {
    public static class EmudhraCtl {
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
    return [AKSe.EmudhraCtl]::GetSequence($EncodedCtl)
}

function Update-EmudhraCtl {
    $directory = Join-Path $env:TEMP ('akse-emudhra-' + [guid]::NewGuid().ToString('N'))
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
        Invoke-EmudhraCertutil -ArgumentList @('-f', '-verifyCTL', 'Disallowed', $directory, $probePath)
    }
    finally {
        if ($certificate) { $certificate.Dispose() }
        if ($rsa) { $rsa.Dispose() }
        Remove-Item -LiteralPath $directory -Recurse -Force -ErrorAction Stop
    }
}

function Assert-EmudhraCtl {
    param([string]$ExpectedSequence)

    Invoke-EmudhraCertutil -ArgumentList @('-verifyCTL', 'Disallowed')
    $encoded = Get-EmudhraEncodedCtl
    $sequence = Get-EmudhraCtlSequence -EncodedCtl $encoded
    if ($sequence -ne $ExpectedSequence) {
        throw "Disallowed CTL sequence '$sequence' does not match the reviewed eMudhra CTL '$ExpectedSequence'. Review Microsoft's mitigation coverage before accepting another sequence."
    }
    return "Verified persisted Windows Disallowed CTL: SequenceNumber=$sequence; Microsoft eMudhra eight-key mitigation."
}

function Invoke-EmudhraRemediation {
    param(
        [switch]$VerifyOnly,
        [string]$ExpectedSequence
    )

    if (-not $VerifyOnly) {
        Update-EmudhraCtl
    }
    return Assert-EmudhraCtl -ExpectedSequence $ExpectedSequence
}

if ($MyInvocation.InvocationName -ne '.') {
    $result = Invoke-EmudhraRemediation -VerifyOnly:$VerifyOnly -ExpectedSequence $ExpectedSequenceNumber
    if ($VerificationLogPath) {
        $result | Set-Content -LiteralPath $VerificationLogPath -Encoding UTF8 -ErrorAction Stop
    }
    Write-Output $result
}
