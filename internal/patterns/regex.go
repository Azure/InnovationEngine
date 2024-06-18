package patterns

import "regexp"

var (
	// An SSH command regex where there must be a username@host somewhere present in the command.
	SshCommand = regexp.MustCompile(
		`(^|\s)\bssh\b\s+([^\s]+(\s+|$))+((?P<username>[a-zA-Z0-9_-]+|\$[A-Z_0-9]+)@(?P<host>[a-zA-Z0-9.-]+|\$[A-Z_0-9]+))`,
	)

	// Multiline quoted string
	MultilineQuotedStringCommand = regexp.MustCompile(`\"(.*\\\n.*)+\"`)

	// Az cli command regex
	AzCommand     = regexp.MustCompile(`az\s+([a-z]+)\s+([a-z]+)`)
	AzGroupDelete = regexp.MustCompile(`az group delete`)

	// ARM regex
	AzResourceURI       = regexp.MustCompile(`\"id\": \"(/subscriptions/[^\"]+)\"`)
	AzResourceGroupName = regexp.MustCompile(`resourceGroups/([^\"\\/\ ]+)`)
)
