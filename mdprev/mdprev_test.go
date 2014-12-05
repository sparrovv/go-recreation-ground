package mdprev

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewMdPrev(t *testing.T) {
	expectedContnet := "Content"
	mdFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(mdFile.Name())

	err = ioutil.WriteFile(mdFile.Name(), []byte(expectedContnet), 0644)
	if err != nil {
		panic(err)
	}

	mdp := NewMdPrev(mdFile.Name())
	if mdp.MdContent != expectedContnet {
		t.Fatalf("The file's contents: %s is not eql to the expected %s", mdp.MdContent, expectedContnet)
	}

	if len(mdp.WSConns) != 0 {
		t.Fatalf("Expecting no WS connections, but got %s", len(mdp.WSConns))
	}
}
