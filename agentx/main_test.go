package agentx

import (
	"os"
	"testing"
)

type tlognone struct{}

func (tlognone) Write(msg string) {
}
func (tlognone) Writef(format string, args ...interface{}) {
}

var lognone = tlognone{}

func TestMain(m *testing.M) {
	setupTesting()
	code := m.Run()
	teardownTesting()
	os.Exit(code)
}

func setupTesting() {
}

func teardownTesting() {
}
