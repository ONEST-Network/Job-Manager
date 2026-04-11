# run.ps1
# Load environment variables from .env file
if (Test-Path ".env") {
    Get-Content .env | Where-Object { $_ -match "^[^#]" -and $_ -match "=" } | ForEach-Object {
        $name, $value = $_ -split '=', 2
        [Environment]::SetEnvironmentVariable($name.Trim(), $value.Trim(), "Process")
    }
    Write-Host "Loaded environment variables from .env"
} else {
    Write-Host ".env file not found, running with existing environment variables"
}

# Run the Go application
go run main.go
