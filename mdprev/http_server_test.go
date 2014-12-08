package mdprev

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	mdPrev := testMdprevObj("#content")
	iHandler := indexHandler(mdPrev)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	iHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

func TestWsHandler(t *testing.T) {
	// assert that connections is added to the slice?
	//mdPrev := testMdprevObj("#content")
	//var blockChan chan bool
	//wshandler := wsHandler(mdPrev, blockChan)

	//wshandler.ServeHTTP(w, req)

	//if len(mdPrev.WSConns) != 1 {
	//t.Errorf("There are no saved websocket connections")
	//}
}
