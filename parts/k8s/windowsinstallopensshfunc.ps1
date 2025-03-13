function
Install-OpenSSH {
    Param(
        [Parameter(Mandatory = $true)][string[]] 
        $SSHKeys
    )

    $adminpath = "c:\ProgramData\ssh"
    $adminfile = "administrators_authorized_keys"

    $sshdService = Get-Service | ? Name -like 'sshd'
    if ($sshdService.Count -eq 0)
    {
        Write-Log "Installing OpenSSH"
        $isAvailable = Get-WindowsCapability -Online | ? Name -like 'OpenSSH*'

        if (!$isAvailable) {
            throw "OpenSSH is not available on this machine"
        }

        # Somehow openssh client got added to Windows 2019 base image. 
        # Remove openssh client in order to install the server.
        Remove-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0
        Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
    }
    else
    {
        if ($sshdService.Status -ne 'Running') {
            Write-Log "OpenSSH Server service detected but not running. Reinstalling OpenSSH..."
            # Somehow openssh client got added to Windows 2019 base image. 
            # Remove openssh client in order to install the server.
            Remove-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
            Remove-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0
            Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
        }
        else {
            Write-Log "OpenSSH Server service detected and running - skipping online install..."
        }
    }

    # It’s by design that files within the C:\Windows\System32\ folder are not modifiable. 
    # When the OpenSSH Server starts, it copies C:\windows\system32\openssh\sshd_config_default to C:\programdata\ssh\sshd_config, if the file does not already exist.
    $OriginalConfigPath = "C:\windows\system32\OpenSSH\sshd_config_default"
    $ConfigDirectory = "C:\programdata\ssh"
    New-Item -ItemType Directory -Force -Path $ConfigDirectory
    $ConfigPath = $ConfigDirectory + "\sshd_config"
    Write-Log "Updating $ConfigPath for CVE-2023-48795"
    $ModifiedConfigContents = Get-Content $OriginalConfigPath `
        | %{ $_ -replace "#RekeyLimit default none", "$&`r`n# Disable cipher to mitigate CVE-2023-48795`r`nCiphers -chacha20-poly1305@openssh.com`r`nMacs -*-etm@openssh.com`r`n" }
    Write-Log "Updating $ConfigPath for CVE-2006-5051"
    $ModifiedConfigContents = $ModifiedConfigContents.Replace("#LoginGraceTime 2m", "LoginGraceTime 0")
    Stop-Service sshd
    Out-File -FilePath $ConfigPath -InputObject $ModifiedConfigContents -Encoding UTF8
    Start-Service sshd

    if (!(Test-Path "$adminpath")) {
        Write-Log "Created new file and text content added"
        New-Item -path $adminpath -name $adminfile -type "file" -value ""
    }

    Write-Log "$adminpath found."
    Write-Log "Adding keys to: $adminpath\$adminfile ..."
    $SSHKeys | foreach-object {
        Add-Content $adminpath\$adminfile $_
    }

    Write-Log "Setting required permissions..."
    icacls $adminpath\$adminfile /remove "NT AUTHORITY\Authenticated Users"
    icacls $adminpath\$adminfile /inheritance:r
    icacls $adminpath\$adminfile /grant SYSTEM:`(F`)
    icacls $adminpath\$adminfile /grant BUILTIN\Administrators:`(F`)

    Write-Log "Restarting sshd service..."
    Restart-Service sshd
    # OPTIONAL but recommended:
    Set-Service -Name sshd -StartupType 'Automatic'

    # Confirm the Firewall rule is configured. It should be created automatically by setup. 
    $firewall = Get-NetFirewallRule -Name *ssh*

    if (!$firewall) {
        throw "OpenSSH is firewall is not configured properly"
    }
    Write-Log "OpenSSH installed and configured successfully"
}
