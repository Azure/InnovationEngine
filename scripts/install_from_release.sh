# Script to install scenarios file.  Pass in language code parameter for a particular language, such as it-it for Italian.
set -e

# TODO: make parameters mandatory
LANG="$1"
RELEASE="$2"

# If no release is specified, download the latest release
if [ "$RELEASE" == "" ]; then
  RELEASE="latest"
fi

# Map the language parameter to the corresponding scenarios file
# If no parameter, download the scenarios from IE
if [ "$LANG" = "" ]; then
  SCENARIOS="https://github.com/Azure/InnovationEngine/releases/download/$RELEASE/scenarios.zip"
# Otherwise, download the scenarios from Microsoft Docs in the appropriate langauge
elif [ "$LANG" = "en-us" ]; then
  SCENARIOS='https://github.com/MicrosoftDocs/executable-docs/releases/download/v1.0.1/scenarios.zip'
else
  SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/v1.0.1/$LANG-scenarios.zip"
fi

# Download the binary
echo "Installing IE & scenarios from the $RELEASE release..."
wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/$RELEASE/ie > /dev/null
wget -q -O scenarios.zip "$SCENARIOS" > /dev/null

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
