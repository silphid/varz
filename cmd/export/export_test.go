package export

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	file, err := ioutil.TempFile("", "TestExport_d.varz")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	yaml :=
`base:
  BASE_VAR1: "base1"
  BASE_VAR2: "base2"
section1:
  ENV_VAR1: "abc"
  ENV_VAR2: 123
  subSection:
    ENV_VAR3: "ghi"
    ENV_VAR4: 456
    subSubSection:
      ENV_VAR5: "jkl"
  subSection2:
    ENV_VAR6: "mno"
`
	if _, err := file.WriteString(yaml); err != nil {
		t.Error(err)
	}
	if err := file.Close(); err != nil {
		t.Error(err)
	}

	gotStdout, gotStderr, err := do(file.Name(), "section1/subSection")
	if err != nil {
		t.Error(err)
	}
	wantStdout :=
`export BASE_VAR1=base1
export BASE_VAR2=base2
export ENV_VAR3=ghi
export ENV_VAR4=456
`
    wantStderr := ""

	assert.Equal(t, wantStdout, gotStdout, "stdout")
	assert.Equal(t, wantStderr, gotStderr, "stderr")
}
