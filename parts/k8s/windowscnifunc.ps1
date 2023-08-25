function Get-HnsPsm1
{
    Param(
        [string]
        $HnsUrl = "https://github.com/Microsoft/SDN/raw/master/Kubernetes/windows/hns.v2.psm1",
        [Parameter(Mandatory=$true)][string]
        $HNSModule
    )
    DownloadFileOverHttp -Url $HnsUrl -DestinationPath "$HNSModule"
}

function Install-SdnBridge
{
    Param(
        [Parameter(Mandatory=$true)][string]
        $Url,
        [Parameter(Mandatory=$true)][string]
        $CNIPath
    )

    $cnizip = [Io.path]::Combine($CNIPath, "cni.zip")
    DownloadFileOverHttp -Url $Url -DestinationPath $cnizip
    Expand-Archive -path $cnizip -DestinationPath $CNIPath
    del $cnizip
}

# TODO: Move the code that creates the wincni configuration file out of windowskubeletfunc.ps1 and put it here