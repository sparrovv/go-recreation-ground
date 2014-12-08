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

	mdPrev := NewMdPrev(mdFile.Name())
	if mdPrev.MdContent != expectedContent {
		t.Fatalf("The file's contents: %s is not eql to the expected %s", mdPrev.MdContent, expectedContent)
	}

	if len(mdPrev.WSConns) != 0 {
		t.Fatalf("Expecting no WS connections, but got %s", len(mdPrev.WSConns))
	}
}

func TestWatcher(t *testing.T) {
	expectedContent := "Content"
	mdPrev := testMdprevObj("")
	defer os.Remove(mdPrev.MdFile)

	mdPrev.Watch()

	var expectedChange bool
	expectedChange = false

	go func() {
		expectedChange = <-mdPrev.MdChanges
	}()

	_ = ioutil.WriteFile(mdPrev.MdFile, []byte(expectedContent), 0644)
	// it sucks, but watcher needs time to notice changes
	time.Sleep(10 * time.Millisecond)

	if expectedChange != true {
		t.Errorf("Expected that watcher will notify about the change, but it didn't")
	}
}

func TestUpdateWSConnections(t *testing.T) {
	content := "A lot of new content"
	mdPrev := testMdprevObj(content)
	defer os.Remove(mdPrev.MdFile)
	var b *bytes.Buffer = new(bytes.Buffer)

	mdPrev.KeepWSConn(b)
	mdPrev.UpdateWSConnections()
	mdPrev.MdChanges <- true

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
