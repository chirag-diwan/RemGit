#!/usr/bin/env bash

set -e

APP_NAME="remgit"
INSTALL_DIR="$HOME/.local/bin"

echo "Detecting shell..."

detect_shell_rc() {
  if [[ -n "$ZSH_VERSION" ]]; then
    echo "$HOME/.zshrc"
  elif [[ -n "$BASH_VERSION" ]]; then
    if [[ -f "$HOME/.bashrc" ]]; then
      echo "$HOME/.bashrc"
    else
      echo "$HOME/.bash_profile"
    fi
  elif [[ -n "$FISH_VERSION" ]]; then
    echo "$HOME/.config/fish/config.fish"
  else
    echo "$HOME/.profile"
  fi
}

RC_FILE="$(detect_shell_rc)"

echo "Using shell config: $RC_FILE"


echo "Building binary..."
go build -o "$APP_NAME"


echo "Installing binary to $INSTALL_DIR"
mkdir -p "$INSTALL_DIR"
mv "$APP_NAME" "$INSTALL_DIR/"


if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  echo "Adding $INSTALL_DIR to PATH"

  if [[ "$RC_FILE" == *fish* ]]; then
    echo "set -Ux PATH $INSTALL_DIR \$PATH" >> "$RC_FILE"
  else
    echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$RC_FILE"
  fi
else
  echo "PATH already contains $INSTALL_DIR"
fi

echo "Reloading shell config"
source "$RC_FILE" 2>/dev/null || true

echo "Installation complete!"
echo "Try running: remgit"
