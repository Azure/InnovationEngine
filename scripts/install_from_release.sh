set -e

# Download the binary from the latest
wget -O ie https://github.com/Azure/InnovationEngine/releases/download/latest/ie

# Setup permissions & move to the local bin
chmod +x ie
mkdir -p ~/.local/bin
mv ie ~/.local/bin

# Export the path to IE if it's not already available
if ![[ "$PATH" =~ "~/.local/bin" || "$PATH" =~ "$HOME/.local/bin" ]]; then
  export PATH="$PATH:~/.local/bin"
fi
