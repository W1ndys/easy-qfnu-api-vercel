<#
.SYNOPSIS
    Deploy the embedded frontend + Go binary to a remote Linux host.

.DESCRIPTION
    1. Build the frontend
    2. Build a Linux amd64 Go binary
    3. Detect a local SSH private key if not provided
    4. Ensure the remote directory exists
    5. Upload the single binary via tar + ssh
    6. Run the remote restart command
#>

param(
    [string]$Server = "mylinux",
    [int]$Port = 22,
    [string]$User = "root",
    [string]$RemotePath = "/root/easy-qfnu-api-lite",
    [string]$RestartCmd = "supervisorctl restart easy-qfnu-api-lite:easy-qfnu-api-lite_00",
    [string]$LocalPath = ".",
    [string]$IdentityFile = ""
)

$ErrorActionPreference = "Stop"

function Write-Step {
    param([string]$Message)
    Write-Host ("[-] " + $Message) -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host ("[+] " + $Message) -ForegroundColor Green
}

function Write-Detail {
    param([string]$Message)
    Write-Host ("[-] " + $Message) -ForegroundColor DarkGray
}

function Invoke-Checked {
    param(
        [Parameter(Mandatory = $true)]
        [string]$FilePath,

        [Parameter()]
        [string[]]$ArgumentList = @(),

        [Parameter(Mandatory = $true)]
        [string]$FailureMessage
    )

    & $FilePath @ArgumentList
    if ($LASTEXITCODE -ne 0) {
        Write-Error $FailureMessage
    }
}

function Get-UnixPath {
    param([Parameter(Mandatory = $true)][string]$WindowsPath)

    $resolvedPath = [System.IO.Path]::GetFullPath($WindowsPath)
    $unixPath = $resolvedPath -replace "\\", "/"

    if ($unixPath -match "^([A-Za-z]):/(.+)$") {
        return "/$($matches[1].ToLower())/$($matches[2])"
    }

    return $unixPath
}

Push-Location $LocalPath

try {
    if ([string]::IsNullOrWhiteSpace($IdentityFile)) {
        $sshDir = Join-Path $env:USERPROFILE ".ssh"
        $possibleKeys = @("id_rsa", "id_ed25519")

        foreach ($keyName in $possibleKeys) {
            $candidate = Join-Path $sshDir $keyName
            if (Test-Path $candidate) {
                $IdentityFile = $candidate
                Write-Step "Detected SSH key: $IdentityFile"
                break
            }
        }
    }

    if ([string]::IsNullOrWhiteSpace($IdentityFile) -or -not (Test-Path $IdentityFile)) {
        Write-Error "SSH private key not found. Use -IdentityFile to specify one."
    }

    $IdentityFile = (Resolve-Path $IdentityFile).Path

    $projectName = "easy-qfnu-api-lite"
    $targetOS = "linux"
    $targetArch = "amd64"
    $binaryName = "$projectName-$targetOS-$targetArch"

    $sshBaseArgs = @(
        "-i", $IdentityFile,
        "-p", $Port.ToString(),
        "-o", "StrictHostKeyChecking=no",
        "$User@$Server"
    )

    Write-Step "Building frontend"
    Push-Location "frontend"
    try {
        if (-not (Test-Path "node_modules")) {
            Write-Detail "Installing frontend dependencies"
            Invoke-Checked -FilePath "npm" -ArgumentList @("ci") -FailureMessage "Frontend dependency install failed."
        }

        Invoke-Checked -FilePath "npm" -ArgumentList @("run", "build") -FailureMessage "Frontend build failed."
        Write-Success "Frontend build completed"
    }
    finally {
        Pop-Location
    }

    $originalCGOEnabled = $env:CGO_ENABLED
    $originalGOOS = $env:GOOS
    $originalGOARCH = $env:GOARCH

    Write-Step "Building Go binary: $binaryName"
    try {
        $env:CGO_ENABLED = "0"
        $env:GOOS = $targetOS
        $env:GOARCH = $targetArch

        Invoke-Checked -FilePath "go" -ArgumentList @("build", "-ldflags", "-s -w", "-o", $binaryName, ".") -FailureMessage "Go build failed."
        Write-Success "Binary build completed"
    }
    finally {
        $env:CGO_ENABLED = $originalCGOEnabled
        $env:GOOS = $originalGOOS
        $env:GOARCH = $originalGOARCH
    }

    Write-Step "Ensuring remote directory exists: $RemotePath"
    Invoke-Checked -FilePath "ssh" -ArgumentList ($sshBaseArgs + @("mkdir -p $RemotePath")) -FailureMessage "Failed to create remote directory."

    $gitBashPaths = @(
        "$env:ProgramFiles\Git\bin\bash.exe",
        "${env:ProgramFiles(x86)}\Git\bin\bash.exe",
        "$env:LOCALAPPDATA\Programs\Git\bin\bash.exe",
        "$env:ProgramFiles\Git\usr\bin\bash.exe"
    )

    $bashExe = $null
    foreach ($candidate in $gitBashPaths) {
        if (Test-Path $candidate) {
            $bashExe = $candidate
            break
        }
    }

    if (-not $bashExe) {
        $gitCmd = Get-Command "git" -ErrorAction SilentlyContinue
        if ($gitCmd) {
            $gitDir = Split-Path (Split-Path $gitCmd.Source -Parent) -Parent
            $gitBash = Join-Path $gitDir "bin\bash.exe"
            if (Test-Path $gitBash) {
                $bashExe = $gitBash
            }
        }
    }

    if (-not $bashExe) {
        Write-Error "Git Bash was not found. Install Git for Windows first."
    }

    Write-Detail "Using Git Bash: $bashExe"

    $unixIdentityFile = Get-UnixPath -WindowsPath $IdentityFile
    $bashCmd = "ssh -i '$unixIdentityFile' -p $Port -o StrictHostKeyChecking=no ${User}@${Server} 'rm -f ${RemotePath}/${binaryName}' && tar cf - $binaryName | ssh -i '$unixIdentityFile' -p $Port -o StrictHostKeyChecking=no ${User}@${Server} 'tar xf - -C ${RemotePath}' && ssh -i '$unixIdentityFile' -p $Port -o StrictHostKeyChecking=no ${User}@${Server} 'chmod +x ${RemotePath}/${binaryName}'"

    Write-Step "Uploading binary"
    Invoke-Checked -FilePath $bashExe -ArgumentList @("-c", $bashCmd) -FailureMessage "Binary upload failed."
    Write-Success "Binary upload completed"

    if (Test-Path $binaryName) {
        Remove-Item $binaryName -Force
        Write-Detail "Removed local artifact: $binaryName"
    }

    if (-not [string]::IsNullOrWhiteSpace($RestartCmd)) {
        Write-Step "Running remote command: $RestartCmd"
        & "ssh" @sshBaseArgs $RestartCmd
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Remote command completed"
        }
        else {
            Write-Warning "Remote command returned a non-zero exit code."
        }
    }

    Write-Host "`nDeployment finished." -ForegroundColor Cyan
}
finally {
    Pop-Location
}
