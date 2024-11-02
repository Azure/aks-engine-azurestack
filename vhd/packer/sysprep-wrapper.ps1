# Set ErrorActionPreference to Continue to suppress errors
$ErrorActionPreference = "Continue"

# Try to run the sysprep script
try {
    Write-Output "Running sysprep.ps1..."
    & "c:\akse-cache\sysprep.ps1"
    Write-Output "sysprep.ps1 completed successfully."
} catch {
    Write-Output "An error occurred while running sysprep.ps1: $_"
}

Write-Output "Continuing with the steps..."