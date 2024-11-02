$ErrorActionPreference = "Continue"

# Stop and remove Azure Agents to enable use in Azure Stack
# If deploying an Azure VM the agents will be re-added to the VMs at deployment time
Stop-Service WindowsAzureGuestAgent
# Stop-Service WindowsAzureNetAgentSvc
Stop-Service RdAgent
& sc.exe delete WindowsAzureGuestAgent
# & sc.exe delete WindowsAzureNetAgentSvc
& sc.exe delete RdAgent
Write-Output '>>> Deleted agents complete ...'

# Remove the WindowsAzureGuestAgent registry key for sysprep 
# This removes AzureGuestAgent from participating in sysprep 
# There was an update that is missing VMAgentDisabler.dll
$path = "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Setup\SysPrepExternal\Generalize"
$generalizeKey = Get-Item -Path $path
$generalizeProperties = $generalizeKey | Select-Object -ExpandProperty property
$values = $generalizeProperties | ForEach-Object {
    New-Object psobject -Property @{"Name"=$_;
    "Value" = (Get-ItemProperty -Path $path -Name $_).$_}
}

$values | ForEach-Object {
    $item = $_;
    if( $item.Value.Contains("VMAgentDisabler.dll")) {
            Write-HOST "Removing " $item.Name - $item.Value;
            Remove-ItemProperty -Path $path -Name $item.Name;
    }
}

# Get-ChildItem c:\\WindowsAzure -Force | Sort-Object -Property FullName -Descending | ForEach-Object { try { Remove-Item -Path $_.FullName -Force -Recurse -ErrorAction SilentlyContinue; } catch { } }
# Remove-Item -Path WSMan:\\Localhost\\listener\\listener* -Recurse -ErrorAction SilentlyContinue

Write-Output '>>> Remove agent script complete ...'