package az

import (
	"testing"
)

func TestSetCorrelationId(t *testing.T) {
	t.Run("Test setting a custom correlation ID", func(t *testing.T) {
		correlationId := "test-correlation-id"
		env := map[string]string{}
		SetCorrelationId(correlationId, env)
		if env["AZURE_HTTP_USER_AGENT"] != "innovation-engine-test-correlation-id" {
			t.Errorf("Expected AZURE_HTTP_USER_AGENT to be set to innovation-engine-test-correlation-id, got %s", env["AZURE_HTTP_USER_AGENT"])
		}
	})
}
