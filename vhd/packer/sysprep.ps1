Write-Output '>>> Waiting for GA Service (RdAgent) to start ...'
while ((Get-Service RdAgent).Status -ne 'Running') { Start-Sleep -s 5 }
Write-Output '>>> Waiting for GA Service (WindowsAzureTelemetryService) to start ...'
while ((Get-Service WindowsAzureTelemetryService) -and ((Get-Service WindowsAzureTelemetryService).Status -ne 'Running')) { Start-Sleep -s 5 }
Write-Output '>>> Waiting for GA Service (WindowsAzureGuestAgent) to start ...'
while ((Get-Service WindowsAzureGuestAgent).Status -ne 'Running') { Start-Sleep -s 5 }
if( Test-Path $Env:SystemRoot\system32\Sysprep\unattend.xml ) {
  Write-Output '>>> Removing Sysprep\unattend.xml ...'
  Remove-Item $Env:SystemRoot\system32\Sysprep\unattend.xml -Force
}
if (Test-Path $Env:SystemRoot\Panther\unattend.xml) {
  Write-Output '>>> Removing Panther\unattend.xml ...'
  Remove-Item $Env:SystemRoot\Panther\unattend.xml -Force
}
Write-Output '>>> Sysprepping VM ...'
& $Env:SystemRoot\System32\Sysprep\Sysprep.exe /oobe /generalize /quiet /quit
while($true) {
  $imageState = (Get-ItemProperty HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Setup\State).ImageState
  Write-Output $imageState
  if ($imageState -eq 'IMAGE_STATE_GENERALIZE_RESEAL_TO_OOBE') { break }
  Start-Sleep -s 5
}
Write-Output '>>> Sysprep complete ...'

# Stop and remove Azure Agents to enable use in Azure Stack
# If deploying an Azure VM the agents will be re-added to the VMs at deployment time
Stop-Service WindowsAzureGuestAgent
# Stop-Service WindowsAzureNetAgentSvc
Stop-Service RdAgent
& sc.exe delete WindowsAzureGuestAgent
# & sc.exe delete WindowsAzureNetAgentSvc
& sc.exe delete RdAgent

# when done clean up
Get-ChildItem c:\\WindowsAzure -Force | Sort-Object -Property FullName -Descending | ForEach-Object { try { Remove-Item -Path $_.FullName -Force -Recurse -ErrorAction SilentlyContinue; } catch { } }
Remove-Item -Path WSMan:\\Localhost\\listener\\listener* -Recurse -ErrorAction SilentlyContinue

Write-Output '>>> Sysprep script complete ...'