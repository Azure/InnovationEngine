# Script to install documentation from our upstream repository.

set -e

# TODO: make parameters mandatory
LANG="$1"
RELEASE="$2"

# If no release is specified, download the latest release
if [ "$RELEASE" == "" ]; then
	RELEASE="latest"
fi

# Set a default scenarios file
SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/scenarios.zip"

# If the LANG parameter was set, download appropriate script
if [ "$LANG" != "" ]; then
	# Map the language parameter to the corresponding scenarios file
	# If no parameter, download the scenarios from IE
	MAIN_LANG_PREFIX="$(echo "$LANG" | head -c2 | tr '[:upper:]' '[:lower:]')"
	LANG_ARRAY=("de" "es" "fr" "it" "nl" "pt" "zh" "cs" "hu" "id" "ja" "ko" "pl" "pt" "ru" "sv" "tr")

	if [[ "${LANG_ARRAY[*]}" =~ "$MAIN_LANG_PREFIX" ]]; then
		SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/$MAIN_LANG_PREFIX-$MAIN_LANG_PREFIX-scenarios.zip"
		if [ "$MAIN_LANG_PREFIX" = "pt" ]; then
			if [ "$LANG" = "pt-pt" ]; then
				SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/pt-pt-scenarios.zip"
			elif [ "$LANG" = "pt-br" ]; then
				SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/pt-br-scenarios.zip"
			else
				SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/pt-pt-scenarios.zip"
			fi
		fi
		if [ "$MAIN_LANG_PREFIX" = "zh" ]; then
			if [ "$LANG" = "zh-cn" ]; then
				SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/zh-cn-scenarios.zip"
			elif [ "$LANG" = "zh-tw" ]; then
				SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/zh-tw-scenarios.zip"
			else
				SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/zh-cn-scenarios.zip"
			fi
		fi
		if [ "$MAIN_LANG_PREFIX" = "cs" ]; then
			SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/$MAIN_LANG_PREFIX-cz-scenarios.zip"
		fi
		if [ "$MAIN_LANG_PREFIX" = "ja" ]; then
			SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/$MAIN_LANG_PREFIX-jp-scenarios.zip"
		fi
		if [ "$MAIN_LANG_PREFIX" = "ko" ]; then
			SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/$MAIN_LANG_PREFIX-kr-scenarios.zip"
		fi
		if [ "$MAIN_LANG_PREFIX" = "sv" ]; then
			SCENARIOS="https://github.com/MicrosoftDocs/executable-docs/releases/download/$RELEASE/$MAIN_LANG_PREFIX-se-scenarios.zip"
		fi
	fi
fi

# Download the scenarios.
echo "Installing scenarios from the $RELEASE release..."
wget -q -O scenarios.zip "$SCENARIOS" >/dev/null

# Unzip the scenarios, overwrite if they already exist.
unzip -o scenarios.zip -d ~ >/dev/null
rm scenarios.zip >/dev/null

echo "Done."
