param(
    [string]$InstallDir,
    [switch]$PrintUrl,
    [switch]$NoModifyPath
)

$ErrorActionPreference = "Stop"

$Repo = "Boursyt/img2ascii"
$BaseUrl = "https://github.com/$Repo/releases/latest/download"

function Stop-Install {
    param([string]$Message)

    Write-Error "img2ascii install error: $Message"
    exit 1
}

function Get-InstallArch {
    $Arch = $env:PROCESSOR_ARCHITEW6432

    if ([string]::IsNullOrWhiteSpace($Arch)) {
        $Arch = $env:PROCESSOR_ARCHITECTURE
    }

    if ([string]::IsNullOrWhiteSpace($Arch)) {
        Stop-Install "could not detect Windows architecture"
    }

    switch ($Arch.ToUpperInvariant()) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Stop-Install "unsupported architecture: $Arch" }
    }
}

function Download-File {
    param(
        [string]$Url,
        [string]$OutputPath
    )

    $Command = Get-Command Invoke-WebRequest -ErrorAction SilentlyContinue
    if ($null -eq $Command) {
        Stop-Install "Invoke-WebRequest is required"
    }

    $Params = @{
        Uri = $Url
        OutFile = $OutputPath
    }

    if ($Command.Parameters.ContainsKey("UseBasicParsing")) {
        $Params.UseBasicParsing = $true
    }

    Invoke-WebRequest @Params
}

function Normalize-PathEntry {
    param([string]$PathEntry)

    if ([string]::IsNullOrWhiteSpace($PathEntry)) {
        return ""
    }

    $ExpandedPath = [Environment]::ExpandEnvironmentVariables($PathEntry.Trim())

    try {
        return ([System.IO.Path]::GetFullPath($ExpandedPath).TrimEnd([char[]]@("\", "/")))
    } catch {
        return $ExpandedPath.TrimEnd([char[]]@("\", "/"))
    }
}

function Add-InstallDirToUserPath {
    param([string]$PathToAdd)

    $ResolvedPath = Normalize-PathEntry $PathToAdd
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ([string]::IsNullOrWhiteSpace($UserPath)) {
        $Entries = @()
    } else {
        $Entries = $UserPath -split ";"
    }

    foreach ($Entry in $Entries) {
        if ((Normalize-PathEntry $Entry).Equals($ResolvedPath, [StringComparison]::OrdinalIgnoreCase)) {
            return $false
        }
    }

    if ([string]::IsNullOrWhiteSpace($UserPath)) {
        $NewPath = $ResolvedPath
    } else {
        $NewPath = "$UserPath;$ResolvedPath"
    }

    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    $env:Path = "$env:Path;$ResolvedPath"

    return $true
}

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    if ([string]::IsNullOrWhiteSpace($env:LOCALAPPDATA)) {
        Stop-Install "LOCALAPPDATA is not set. Pass -InstallDir explicitly"
    }

    $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\img2ascii\bin"
}

$Arch = Get-InstallArch
$Archive = "img2ascii-windows-$Arch.zip"
$Url = "$BaseUrl/$Archive"

if ($PrintUrl) {
    Write-Output $Url
    exit 0
}

$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("img2ascii-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir $Archive
$ExtractDir = Join-Path $TempDir "extract"
$TargetPath = Join-Path $InstallDir "img2ascii.exe"

New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

try {
    New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null

    Download-File $Url $ArchivePath
    Expand-Archive -LiteralPath $ArchivePath -DestinationPath $ExtractDir -Force

    $BinaryPath = Join-Path $ExtractDir "img2ascii.exe"
    if (-not (Test-Path -LiteralPath $BinaryPath -PathType Leaf)) {
        Stop-Install "archive did not contain img2ascii.exe"
    }

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Copy-Item -LiteralPath $BinaryPath -Destination $TargetPath -Force

    $PathUpdated = $false
    if (-not $NoModifyPath) {
        $PathUpdated = Add-InstallDirToUserPath $InstallDir
    }

    Write-Output "img2ascii installed to $TargetPath"

    if ($PathUpdated) {
        Write-Output "Added $InstallDir to the user PATH. Restart your terminal if img2ascii is not found."
    } elseif ($NoModifyPath) {
        Write-Output "PATH was not modified."
    }
} finally {
    Remove-Item -LiteralPath $TempDir -Recurse -Force -ErrorAction SilentlyContinue
}
