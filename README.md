# img2ascii

Convert images to ASCII art from the command line.

## Install

The install scripts download the matching asset from the latest GitHub release.

### Linux

```sh
curl -fsSL https://raw.githubusercontent.com/Boursyt/img2ascii/main/install/install.sh | bash
```

### macOS

```sh
curl -fsSL https://raw.githubusercontent.com/Boursyt/img2ascii/main/install/install.sh | bash
```

### Windows

Run this in PowerShell:

```powershell
$script = Join-Path $env:TEMP "install-img2ascii.ps1"; curl.exe -fsSL https://raw.githubusercontent.com/Boursyt/img2ascii/main/install/install.ps1 -o $script; powershell -ExecutionPolicy Bypass -NoProfile -File $script
```

## Install directory

Linux and macOS install to `/usr/local/bin` by default. You can override it:

```sh
curl -fsSL https://raw.githubusercontent.com/Boursyt/img2ascii/main/install/install.sh | INSTALL_DIR="$HOME/.local/bin" bash
```

Windows installs to `%LOCALAPPDATA%\Programs\img2ascii\bin` by default and adds it to the user `PATH`.

## Usage

```sh
img2ascii path/to/image.png
```

```sh
img2ascii -width 80 -height 40 -ramp blocks path/to/image.png
```

Available ramps:

```text
short      [ .:-=+*#%@]
short-alt  [ .,:;i1tfLCG08@]
bourke70   [ .'`^",:;Il!i><~+_-?][}{1)(|\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$]
blocks     [ ░▒▓█]
inverted   [@%#*+=-:. ]
binary     [ #]
```

You can also print this list with:

```sh
img2ascii -display-ramp
```
