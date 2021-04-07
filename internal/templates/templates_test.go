package templates

import (
	"testing"
)

func TestGenerateTemplates(t *testing.T) {
	t.Run("returns error if generating templates fails", func(t *testing.T) {
		_, err := GenerateTemplates("")

		if err == nil {
			t.Error("expected for error to be returned from template generator")
		}
	})
}
