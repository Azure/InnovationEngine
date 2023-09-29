package engine

import (
	"testing"
)

func TestRegex(t *testing.T) {

	t.Run("Test ssh command regex", func(t *testing.T) {
		testCases := []string{
			"Run ssh -i key.pem username@host to connect",
			"ssh -p 22 -L 8080:localhost:8080 username@host",
			"ssh -Y username@host",
			"Use ssh to connect",
			"sshusername@host is not correct",
			" ssh username@domain.com",
			"Invalid ssh username@@domain.com",
			"ssh -o StrictHostKeyChecking=no $MY_USERNAME@$IP_ADDRESS",
		}

		testResults := []bool{
			true,
			true,
			true,
			false,
			false,
			false,
			false,
			true,
		}

		for index, testCase := range testCases {
			match := sshCommand.FindString(testCase)
			if match == "" && testResults[index] {
				t.Errorf("Expected match not found: %s\n", testCase)
			} else if match != "" && !testResults[index] {
				t.Errorf("Unexpected match found: %s\n", testCase)
			}
		}
	})

}
