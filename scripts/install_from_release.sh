# Script to install scenarios file.  Pass in language code parameter for a particular language, such as it-it for Italian.
set -e

# Define the language parameter
LANG="${1:-''}"
SCENARIOS=""

# Map the language parameter to the corresponding scenarios file
# If no parameter, download the scenarios from IE
if [ "$lang" = "" ]; then
  SCENARIOS='https://github.com/Azure/InnovationEngine/releases/download/latest/scenarios.zip'
# Otherwise, download the scenarios from Microsoft Docs in the appropriate langauge
elif [ "$lang" = "en-us" ]; then
  SCENARIOS='https://github.com/MicrosoftDocs/executable-docs/releases/download/v1.0.1/scenarios.zip'
else
  SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/v1.0.1/$lang-scenarios.zip"
fi

# Download the binary from the latest
echo "Installing IE & scenarios from the latest release..."
wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/latest/ie > /dev/null
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
