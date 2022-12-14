<#
.DESCRIPTION
    This script rotates a windows node certificates.
    It assumes that client.key, client.crt and ca.crt will be dropped in $env:temp.
#>

. c:\AzureData\k8s\windowskubeletfunc.ps1
. c:\AzureData\k8s\kuberneteswindowsfunctions.ps1

$global:KubeDir = "c:\k"

$global:AgentKeyPath = [io.path]::Combine($env:temp, "client.key")
$global:AgentCertificatePath = [io.path]::Combine($env:temp, "client.crt")
$global:CACertificatePath = [io.path]::Combine($env:temp, "ca.crt")

function Prereqs {
    Assert-FileExists $global:AgentKeyPath
    Assert-FileExists $global:AgentCertificatePath
    Assert-FileExists $global:CACertificatePath
}

function Backup {
    Copy-Item "c:\k\config" "c:\k\config.bak"
    Copy-Item "c:\k\ca.crt" "c:\k\ca.crt.bak"
}

function Update-CACertificate {
    Write-Log "Write ca root"
    Write-CACert -CACertificate $global:CACertificate -KubeDir $global:KubeDir
}

function Update-KubeConfig {
    Write-Log "Write kube config"
    $ClusterConfiguration = ConvertFrom-Json ((Get-Content "c:\k\kubeclusterconfig.json" -ErrorAction Stop) | out-string) 
    $MasterIP = $ClusterConfiguration.Kubernetes.ControlPlane.IpAddress

    $CloudProviderConfig = ConvertFrom-Json ((Get-Content "c:\k\azure.json" -ErrorAction Stop) | out-string) 
    $MasterFQDNPrefix = $CloudProviderConfig.ResourceGroup

    $AgentKey = [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((Get-Content -Raw $AgentKeyPath)))
    $AgentCertificate = [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((Get-Content -Raw $AgentCertificatePath)))

    Write-KubeConfig -CACertificate $global:CACertificate `
        -KubeDir $global:KubeDir `
        -MasterFQDNPrefix $MasterFQDNPrefix `
        -MasterIP $MasterIP `
        -AgentKey $AgentKey `
        -AgentCertificate $AgentCertificate
}

function Force-Kubelet-CertRotation {
    Remove-Item "/var/lib/kubelet/pki/kubelet-client-current.pem" -Force -ErrorAction Ignore
    Remove-Item "/var/lib/kubelet/pki/kubelet.crt" -Force -ErrorAction Ignore
    Remove-Item "/var/lib/kubelet/pki/kubelet.key" -Force -ErrorAction Ignore

    try {
        $err = Retry-Command -Command "c:\k\windowsnodereset.ps1" -Args @{Foo="Bar"} -Retries 3 -RetryDelaySeconds 10
    } catch {
        Write-Error "Error reseting Windows node. Error: $_"
        throw $_
    }
}

function Start-CertRotation {
    try
    {
        $global:CACertificate = [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((Get-Content -Raw $CACertificatePath)))

        Prereqs
        Update-CACertificate
        Update-KubeConfig
        Force-Kubelet-CertRotation
    }
    catch
    {
        Write-Error $_
        throw $_
    }
}

function Clean {
    Remove-Item "c:\k\config.bak" -Force -ErrorAction Ignore
    Remove-Item "c:\k\ca.crt.bak" -Force -ErrorAction Ignore
    Remove-Item $global:AgentKeyPath -Force -ErrorAction Ignore
    Remove-Item $global:AgentCertificatePath -Force -ErrorAction Ignore
    Remove-Item $global:CACertificatePath -Force -ErrorAction Ignore
}
