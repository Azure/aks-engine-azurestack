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