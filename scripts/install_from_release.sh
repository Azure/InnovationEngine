# Script to install scenarios file.  Pass in language code parameter for a particular language, such as it-it for Italian.
set -e

RELEASE="$1"

# If no release is specified, download the latest release
if [ "$RELEASE" == "" ]; then
	RELEASE="latest"
fi

# Download the binary
echo "Installing IE from the $RELEASE release..."
wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/"$RELEASE"/ie >/dev/null

# Setup permissions & move to the local bin
chmod +x ie >/dev/null
mkdir -p ~/.local/bin >/dev/null
mv ie ~/.local/bin >/dev/null

# Export the path to IE if it's not already available
if [[ !"$PATH" =~ "~/.local/bin" || !"$PATH" =~ "$HOME/.local/bin" ]]; then
	export PATH="$PATH:~/.local/bin"
fi

echo "Done."
