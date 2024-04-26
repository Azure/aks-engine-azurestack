<#
    .SYNOPSIS
        Used to produce Windows AKS images.

    .DESCRIPTION
        This script is used by packer to produce Windows AKS images.
#>

param()

$ErrorActionPreference = "Stop"

filter Timestamp { "$(Get-Date -Format o): $_" }

$global:containerdPackageUrl = "https://mobyartifacts.azureedge.net/moby/moby-containerd/1.6.28+azure/windows/windows_amd64/moby-containerd-1.6.28+azure-u1.amd64.zip"

function Write-Log($Message) {
    $msg = $message | Timestamp
    Write-Output $msg
}
function Get-ContainerImages {
    $imagesToPull = @(
        "mcr.microsoft.com/windows/servercore:ltsc2019",
        "mcr.microsoft.com/windows/nanoserver:1809",
        "mcr.microsoft.com/oss/kubernetes/pause:3.4.1",
        "mcr.microsoft.com/oss/kubernetes/pause:3.8",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.27.13",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.28.5",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.28.3",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.28.3-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.29.1",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.29.1-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/csi-node-driver-registrar:v2.6.2",
        "mcr.microsoft.com/oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.8.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.10.0",
        "mcr.microsoft.com/oss/kubernetes/windows-host-process-containers-base-image:v1.0.0")

    # start containerd to pre-pull the images to disk on VHD
    # CSE will configure and register containerd as a service at deployment time
    Start-Job -Name containerd -ScriptBlock { containerd.exe }
    foreach ($image in $imagesToPull) {
        & ctr.exe -n k8s.io images pull $image
    }
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
            "https://globalcdn.nuget.org/packages/microsoft.applicationinsights.2.11.0.nupkg",
            "https://akseashartifacts.blob.core.windows.net/windows/provisioning/signedscripts-v0.0.18.zip"
        );
        "c:\akse-cache\containerd\"   = @(
            $global:containerdPackageUrl
        );
        "c:\akse-cache\csi-proxy\"    = @(
            "https://kubernetesartifacts.azureedge.net/csi-proxy/v1.1.3/binaries/csi-proxy-v1.1.3.tar.gz"
        );
        "c:\akse-cache\win-k8s\"      = @(
            "https://kubernetesartifacts.azureedge.net/kubernetes/v1.28.6/windowszip/v1.28.6-1int.zip",
            "https://kubernetesartifacts.azureedge.net/kubernetes/v1.27.10/windowszip/v1.27.10-1int.zip"
        );
        "c:\akse-cache\win-vnet-cni\" = @(
            "https://kubernetesartifacts.azureedge.net/azure-cni/v1.4.32/binaries/azure-vnet-cni-singletenancy-windows-amd64-v1.4.32.zip"
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

    Write-Log "Enable a WCIFS fix in 2022-10B"
    $currentValue=(Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\wcifs" -Name WcifsSOPCountDisabled -ErrorAction Ignore)
    if (![string]::IsNullOrEmpty($currentValue)) {
        Write-Log "The current value of WcifsSOPCountDisabled is $currentValue"
    }
    Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\wcifs" -Name WcifsSOPCountDisabled -Value 0 -Type DWORD
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