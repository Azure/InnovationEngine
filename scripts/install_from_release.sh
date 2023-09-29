set -e

# Download the binary from the latest
echo "Installing IE & scenarios from the latest release..."
wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/latest/ie > /dev/null
wget -q -O scenarios.zip https://github.com/Azure/InnovationEngine/releases/download/latest/scenarios.zip > /dev/null

# Setup permissions & move to the local bin
chmod +x ie > /dev/null
mkdir -p ~/.local/bin > /dev/null
mv ie ~/.local/bin > /dev/null

# Unzip the scenarios, overwrite if they already exist.
unzip -o scenarios.zip -d ~ > /dev/null
rm scenarios.zip > /dev/null

# Export the path to IE if it's not already available
if [[ !"$PATH" =~ "~/.local/bin" || !"$PATH" =~ "$HOME/.local/bin" ]]; then
  export PATH="$PATH:~/.local/bin"
fi

echo "Done."
