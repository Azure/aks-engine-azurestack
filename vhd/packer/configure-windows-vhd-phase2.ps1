<#
    .SYNOPSIS
        Used to produce Windows AKS images.

    .DESCRIPTION
        This script is used by packer to produce Windows AKS images.
#>

param()

$ErrorActionPreference = "Stop"

filter Timestamp { "$(Get-Date -Format o): $_" }

$global:containerdPackageUrl = "https://mobyartifacts.azureedge.net/moby/moby-containerd/1.6.36+azure/windows/windows_amd64/moby-containerd-1.6.36+azure-u1.amd64.zip"

function Write-Log($Message) {
    $msg = $message | Timestamp
    Write-Output $msg
}
function Get-ContainerImages {
    $containerdImagePullNotesFilePath = "c:\containerd-image-pull-notes.txt"
    $imagesToPull = @(
        "mcr.microsoft.com/windows/servercore:ltsc2019",
        "mcr.microsoft.com/windows/nanoserver:1809",
        "mcr.microsoft.com/oss/kubernetes/pause:3.8",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.31.6",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.32.5",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.33.0",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.31.12",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.32.11",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.33.5",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.31.12-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.32.11-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.33.5-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/csi-node-driver-registrar:v2.15.0",
        "mcr.microsoft.com/oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.17.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",
        "mcr.microsoft.com/oss/kubernetes/windows-host-process-containers-base-image:v1.0.0")

    # start containerd to pre-pull the images to disk on VHD
    # CSE will configure and register containerd as a service at deployment time
    Start-Job -Name containerd -ScriptBlock { containerd.exe }
    foreach ($image in $imagesToPull) {
        & ctr.exe -n k8s.io images pull $image > $containerdImagePullNotesFilePath
    }
    Write-Log "Begin listing containerd images"
    $imagesList = & ctr.exe -n k8s.io images list
    foreach ($line in $imagesList) {
        Write-Output $line
    }
    Write-Log "End listing containerd images"
    Stop-Job  -Name containerd
    Remove-Job -Name containerd
}

function Get-FilesToCacheOnVHD {
    Write-Log "Caching misc files on VHD"

    $map = @{
        "c:\akse-cache\"              = @(
            "https://github.com/Azure/aks-engine-azurestack/raw/master/scripts/collect-windows-logs.ps1",
            "https://github.com/Microsoft/SDN/raw/master/Kubernetes/flannel/l2bridge/cni/win-bridge.exe",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/debug/collectlogs.ps1",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/debug/dumpVfpPolicies.ps1",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/debug/portReservationTest.ps1",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/debug/starthnstrace.cmd",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/debug/startpacketcapture.cmd",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/debug/stoppacketcapture.cmd",
            "https://github.com/Microsoft/SDN/raw/master/Kubernetes/windows/debug/VFP.psm1",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/helper.psm1",
            "https://github.com/microsoft/SDN/raw/master/Kubernetes/windows/hns.v2.psm1"
        );
        "c:\akse-cache\containerd\"   = @(
            $global:containerdPackageUrl
        );
        "c:\akse-cache\csi-proxy\"    = @(
            "https://packages.aks.azure.com/csi-proxy/v1.1.3/binaries/csi-proxy-v1.1.3.tar.gz"
        );
        "c:\akse-cache\win-k8s\"      = @(
            "https://packages.aks.azure.com/kubernetes/v1.33.5/windowszip/v1.33.5-1int.zip",
            "https://packages.aks.azure.com/kubernetes/v1.32.9/windowszip/v1.32.9-1int.zip",
            "https://packages.aks.azure.com/kubernetes/v1.31.13/windowszip/v1.31.13-1int.zip"
        );
        "c:\akse-cache\win-vnet-cni\" = @(
            "https://packages.aks.azure.com/azure-cni/v1.4.59/binaries/azure-vnet-cni-windows-amd64-v1.4.59.zip"
        )
    }

    foreach ($dir in $map.Keys) {
        New-Item -ItemType Directory $dir -Force | Out-Null

        foreach ($URL in $map[$dir]) {
            $fileName = [IO.Path]::GetFileName($URL)
            $dest = [IO.Path]::Combine($dir, $fileName)

            Write-Log "Downloading $URL to $dest"
            curl.exe -f --retry 5 --retry-delay 0 -L $URL -o $dest
            if ($LASTEXITCODE) {
                throw "Curl exited with '$LASTEXITCODE' while attemping to downlaod '$URL'"
            }
        }
    }

    $acrCredentialProviderUrls = @(
        @{ Url = "https://github.com/kubernetes-sigs/cloud-provider-azure/releases/download/v1.31.7/azure-acr-credential-provider-windows-amd64.exe"; K8sVersion = "v1.31" },
        @{ Url = "https://github.com/kubernetes-sigs/cloud-provider-azure/releases/download/v1.30.13/azure-acr-credential-provider-windows-amd64.exe"; K8sVersion = "v1.30" }
    )

    $credentialProviderDir = "c:\k\credential-provider\"
    New-Item -ItemType Directory $credentialProviderDir -Force | Out-Null
    
    foreach ($providerInfo in $acrCredentialProviderUrls) {
        $versionedFileName = "azure-acr-credential-provider-windows-amd64-$($providerInfo.K8sVersion).exe"
        $dest = [IO.Path]::Combine($credentialProviderDir, $versionedFileName)

        curl.exe -f --retry 5 --retry-delay 0 -L $providerInfo.Url -o $dest
        if ($LASTEXITCODE) {
            throw "Curl exited with '$LASTEXITCODE' while attempting to download '$($providerInfo.Url)'"
        }
    }
}

function Install-ContainerD {
    Write-Log "Getting containerD binaries from $global:containerdPackageUrl"

    $installDir = "c:\program files\containerd"
    Write-Log "Installing containerd to $installDir"
    New-Item -ItemType Directory $installDir -Force | Out-Null

    if ($global:containerdPackageUrl.endswith(".zip")) {
        $zipPath = [IO.Path]::Combine($installDir, "containerd.zip")
        Invoke-WebRequest -UseBasicParsing -Uri $global:containerdPackageUrl -OutFile $zipPath
        Expand-Archive -path $zipPath -DestinationPath $installDir -Force
        Remove-Item -Path $zipPath | Out-Null
    } else {
        $tarPath = [IO.Path]::Combine($installDir, "containerd.tar.gz")
        Invoke-WebRequest -UseBasicParsing -Uri $global:containerdPackageUrl -OutFile $tarPath
        tar -xzf $tarPath --strip=1 -C $installDir
        Remove-Item -Path $tarPath | Out-Null
    }

    $newPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine) + ";$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, [EnvironmentVariableTarget]::Machine)
    $env:Path += ";$installDir"
}

function Set-WinRmServiceAutoStart {
    Write-Log "Setting WinRM service start to auto"
    sc.exe config winrm start=auto
}

function Update-Registry {
    # Enable HNS fixed gated behind reg keys for Windows Server 2019
    Write-Log "Enable a HNS fix (0x40) in 2022-11B and another HNS fix (0x10)"
    $hnsControlFlag=0x50
    $currentValue=(Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\hns\State" -Name HNSControlFlag -ErrorAction Ignore)
    if (![string]::IsNullOrEmpty($currentValue)) {
        Write-Log "The current value of HNSControlFlag is $currentValue"
        $hnsControlFlag=([int]$currentValue.HNSControlFlag -bor $hnsControlFlag)
    }
    Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\hns\State" -Name HNSControlFlag -Value $hnsControlFlag -Type DWORD

    # --- WCIFS Fix ---
    Write-Log "Enable a WCIFS fix in 2022-10B"
    $currentValue=(Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\wcifs" -Name WcifsSOPCountDisabled -ErrorAction Ignore)
    if (![string]::IsNullOrEmpty($currentValue)) {
        Write-Log "The current value of WcifsSOPCountDisabled is $currentValue"
    }
    Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\wcifs" -Name WcifsSOPCountDisabled -Value 0 -Type DWORD

    # --- TLS 1.0 and 1.1 Disablement ---
    Write-Host "Disabling TLS 1.0 and 1.1 for Client and Server roles"
    $basePath = "HKLM:\SYSTEM\CurrentControlSet\Control\SecurityProviders\SCHANNEL\Protocols"
    $protocols = @(
        @{ Name = "TLS 1.0"; Roles = @("Client", "Server") },
        @{ Name = "TLS 1.1"; Roles = @("Client", "Server") }
    )
    foreach ($protocol in $protocols) {
        foreach ($role in $protocol.Roles) {
            $regPath = Join-Path -Path $basePath -ChildPath "$($protocol.Name)\$role"
            if (-not (Test-Path $regPath)) {
                New-Item -Path $regPath -Force | Out-Null
            }
            New-ItemProperty -Path $regPath -Name "Enabled" -Value 0 -PropertyType "DWORD" -Force
            New-ItemProperty -Path $regPath -Name "DisabledByDefault" -Value 1 -PropertyType "DWORD" -Force
        }
    }
    Write-Host "TLS 1.0 and 1.1 have been disabled successfully."
}

# Disable progress writers for this session to greatly speed up operations such as Invoke-WebRequest
$ProgressPreference = 'SilentlyContinue'

Write-Log "Performing actions for provisioning phase 2 for container runtime 'containerd'"
Set-WinRmServiceAutoStart
Install-ContainerD
Update-Registry
Get-ContainerImages
Get-FilesToCacheOnVHD
(New-Guid).Guid | Out-File -FilePath 'c:\vhd-id.txt'