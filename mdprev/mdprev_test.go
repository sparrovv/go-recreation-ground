package mdprev

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNewMdPrev(t *testing.T) {
	expectedContent := "Content"
	mdFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(mdFile.Name())

	err = ioutil.WriteFile(mdFile.Name(), []byte(expectedContent), 0644)
	if err != nil {
		panic(err)
	}

	mdp := NewMdPrev(mdFile.Name())
	if mdp.MdContent != expectedContent {
		t.Fatalf("The file's contents: %s is not eql to the expected %s", mdp.MdContent, expectedContent)
	}

	if len(mdp.WSConns) != 0 {
		t.Fatalf("Expecting no WS connections, but got %s", len(mdp.WSConns))
	}
}

func TestWatcher(t *testing.T) {
	expectedContent := "Content"
	mdp := testMdprevObj("")
	defer os.Remove(mdp.MdFile)

	mdp.Watch()

	var expectedChange bool
	expectedChange = false

	go func() {
		expectedChange = <-mdp.MdChanges
	}()

	_ = ioutil.WriteFile(mdp.MdFile, []byte(expectedContent), 0644)
	// it sucks, but watcher needs time to notice changes
	time.Sleep(10 * time.Millisecond)

	if expectedChange != true {
		t.Errorf("Expected that watcher will notify about the change, but it didn't")
	}
}

func TestUpdateWSConnections(t *testing.T) {
	content := "A lot of new content"
	mdp := testMdprevObj(content)
	defer os.Remove(mdp.MdFile)
	var b *bytes.Buffer = new(bytes.Buffer)

	mdp.KeepWSConn(b)
	mdp.UpdateWSConnections()
	mdp.MdChanges <- true

	if b.String() != content {
		t.Errorf("Expecting that io.Writer gets %s, but got %s", "foo", b.String())
	}
}

func testMdprevObj(content string) *MdPrev {
	mdFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile(mdFile.Name(), []byte(content), 0644)
	return NewMdPrev(mdFile.Name())
}
