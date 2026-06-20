#!/usr/bin/env bash

set -euo pipefail

REPO="Boursyt/img2ascii"
BASE_URL="https://github.com/${REPO}/releases/latest/download"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

abort() {
  echo "img2ascii install error: $*" >&2
  exit 1
}

detect_os() {
  case "$(uname -s)" in
    Linux) echo "linux" ;;
    Darwin) echo "darwin" ;;
    *) abort "unsupported OS: $(uname -s). Use install.ps1 on Windows" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    riscv64) echo "riscv64" ;;
    *) abort "unsupported architecture: $(uname -m)" ;;
  esac
}

download() {
  url="$1"
  output="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$url" -O "$output"
  else
    abort "curl or wget is required"
  fi
}

install_binary() {
  source="$1"
  target="${INSTALL_DIR}/img2ascii"

  mkdir -p "$INSTALL_DIR" 2>/dev/null || sudo mkdir -p "$INSTALL_DIR"

  if [ -w "$INSTALL_DIR" ]; then
    install -m 0755 "$source" "$target"
  else
    sudo install -m 0755 "$source" "$target"
  fi
}

os="$(detect_os)"
arch="$(detect_arch)"

case "${os}-${arch}" in
  darwin-amd64|darwin-arm64|linux-amd64|linux-arm64|linux-riscv64)
    archive="img2ascii-${os}-${arch}.tar.gz"
    ;;
  *)
    abort "no release asset for ${os}-${arch}"
    ;;
esac

url="${BASE_URL}/${archive}"

if [ "${1:-}" = "--print-url" ]; then
  echo "$url"
  exit 0
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

download "$url" "${tmp_dir}/${archive}"
tar -xzf "${tmp_dir}/${archive}" -C "$tmp_dir"

[ -f "${tmp_dir}/img2ascii" ] || abort "archive did not contain img2ascii"

install_binary "${tmp_dir}/img2ascii"

echo "img2ascii installed to ${INSTALL_DIR}/img2ascii"
