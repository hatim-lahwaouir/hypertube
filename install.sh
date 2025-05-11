#!/bin/zsh


# Set Go version
GO_VERSION="1.23.0"
GO_TARBALL="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TARBALL}"
INSTALL_DIR="$HOME/go"


# Download Go
echo "📥 Downloading Go $GO_VERSION..."

wget -q --show-progress "$GO_URL" -O "$HOME/$GO_TARBALL"

# Extract
echo "📦 Extracting Go to $INSTALL_DIR..."
tar -xzf "$HOME/$GO_TARBALL" -C "$HOME"

rm "$HOME/$GO_TARBALL"

# Set up environment variables
echo "⚙️ Configuring environment in ~/.bashrc..."
{
    echo ""
    echo "# Go environment setup"
    echo "export GOROOT=\$HOME/go"
    echo "export GOPATH=$HOME/gopkg"
    echo "export PATH=\$GOROOT/bin:\$PATH"
} >> "$HOME/.zshrc"


# Apply changes
echo "🔁 Reloading ~/.bashrc..."
source "$HOME/.zshrc"

# Test installation
echo "✅ Installation complete. Checking Go version:"
"$HOME/go/bin/go" version
