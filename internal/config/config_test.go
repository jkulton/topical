package config

import (
	"os"
	"testing"
)

var initialArgs []string

func testSetup() {
	initialArgs = os.Args
}

func testTeardown() {
	os.Args = initialArgs
}

func TestParseAppConfig(t *testing.T) {
	t.Run("parses known flags and returns config object", func(t *testing.T) {
		want := AppConfig{Port: 1234, DBFile: "./test.db", SessionKey: "big_session_key"}
		testSetup()

		mockArgs := []string{"_", "-p=1234", "-db-file=./test.db", "-session-key=big_session_key"}
		os.Args = mockArgs
		got := ParseAppConfig()

		if got != want {
			t.Error("parsed config does not match expected config")
		}

		testTeardown()
	})
}
