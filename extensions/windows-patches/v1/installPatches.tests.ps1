BeforeAll {
    . $PSScriptRoot\installPatches.ps1
    $validHash = '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef'
}

Describe 'Get-PatchDefinitions' {
    It 'creates a patch definition for a valid HTTPS URI and SHA-256 hash' {
        $patches = @(Get-PatchDefinitions -URIs 'https://example.test/patch.msu?token=secret' -SHA256Hashes $validHash)

        $patches.Count | Should -Be 1
        $patches[0].FileName | Should -Be 'patch.msu'
        $patches[0].Extension | Should -Be '.msu'
        $patches[0].ExpectedSHA256 | Should -Be $validHash
    }

    It 'rejects an HTTP URI' {
        { Get-PatchDefinitions -URIs 'http://example.test/patch.msu' -SHA256Hashes $validHash } |
            Should -Throw '*must use HTTPS*'
    }

    It 'rejects mismatched URI and hash counts' {
        { Get-PatchDefinitions -URIs @('https://example.test/one.msu', 'https://example.test/two.msu') -SHA256Hashes $validHash } |
            Should -Throw '*corresponding SHA-256 hash*'
    }

    It 'rejects a malformed SHA-256 hash' {
        { Get-PatchDefinitions -URIs 'https://example.test/patch.msu' -SHA256Hashes 'not-a-hash' } |
            Should -Throw '*64 hexadecimal characters*'
    }

    It 'rejects an unsupported patch extension' {
        { Get-PatchDefinitions -URIs 'https://example.test/patch.zip' -SHA256Hashes $validHash } |
            Should -Throw '*doesn''t know how to install*'
    }
}

Describe 'DownloadAndVerifyFile' {
    BeforeEach {
        Mock Invoke-WebRequest
        Mock Remove-Item
    }

    It 'accepts a matching SHA-256 hash without deleting the file' {
        Mock Get-FileHash { [PSCustomObject]@{ Hash = $validHash.ToUpperInvariant() } }

        { DownloadAndVerifyFile -URI 'https://example.test/patch.msu' -FullName 'TestDrive:\patch.msu' -ExpectedSHA256 $validHash } |
            Should -Not -Throw

        Should -Invoke Invoke-WebRequest -Times 1 -Exactly
        Should -Invoke Remove-Item -Times 0 -Exactly
    }

    It 'deletes the file and rejects a mismatched SHA-256 hash' {
        Mock Get-FileHash { [PSCustomObject]@{ Hash = ('f' * 64) } }

        { DownloadAndVerifyFile -URI 'https://example.test/patch.msu' -FullName 'TestDrive:\patch.msu' -ExpectedSHA256 $validHash } |
            Should -Throw '*SHA-256 verification failed*'

        Should -Invoke Remove-Item -Times 1 -Exactly -ParameterFilter { $LiteralPath -eq 'TestDrive:\patch.msu' }
    }

    It 'deletes a partial file when the download fails' {
        Mock Invoke-WebRequest { throw 'download failed' }

        { DownloadAndVerifyFile -URI 'https://example.test/patch.msu' -FullName 'TestDrive:\patch.msu' -ExpectedSHA256 $validHash } |
            Should -Throw '*download failed*'

        Should -Invoke Remove-Item -Times 1 -Exactly -ParameterFilter { $LiteralPath -eq 'TestDrive:\patch.msu' }
    }
}

Describe 'Install-Patches' {
    BeforeEach {
        Mock DownloadAndVerifyFile
        Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
        Mock Wait-PatchProcess { 0 }
        Mock schtasks
    }

    It 'does not enable test signing by default' {
        Install-Patches -URIs 'https://example.test/patch.exe' -SHA256Hashes $validHash

        Should -Invoke Start-Process -Times 0 -Exactly -ParameterFilter { $FilePath -eq 'bcdedit.exe' }
        Should -Invoke Start-Process -Times 1 -Exactly -ParameterFilter { $FilePath -like '*patch.exe' }
    }

    It 'enables test signing when explicitly requested for an executable patch' {
        Install-Patches -URIs 'https://example.test/patch.exe' -SHA256Hashes $validHash -EnableTestSigning

        Should -Invoke Start-Process -Times 1 -Exactly -ParameterFilter { $FilePath -eq 'bcdedit.exe' }
        Should -Invoke Start-Process -Times 1 -Exactly -ParameterFilter { $FilePath -like '*patch.exe' }
    }
}