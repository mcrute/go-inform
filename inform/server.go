package inform

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type StateTree struct {
	states map[string]map[int]int
}

func NewStateTree() *StateTree {
	return &StateTree{make(map[string]map[int]int)}
}

func (t *StateTree) ensureNode(device string, port int) {
	_, ok := t.states[device]
	if !ok {
		t.states[device] = make(map[int]int)
	}

	_, ok = t.states[device][port]
	if !ok {
		t.states[device][port] = 0
	}
}

func (t *StateTree) GetState(device string, port int) int {
	t.ensureNode(device, port)
	return t.states[device][port]
}

func (t *StateTree) SetState(device string, port, value int) {
	t.ensureNode(device, port)
	t.states[device][port] = value
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		t := time.Now().Format("02/Jan/2006 15:04:05")
		log.Printf("%s - - [%s] \"%s %s %s\"", r.RemoteAddr, t, r.Method, r.URL, r.Proto)
		// Addr - - [D/M/Y H:M:S] "Method RequestURI Proto" Code Size
		// 127.0.0.1 - - [24/Sep/2016 14:30:35] "GET / HTTP/1.1" 200 -
	})
}

type InformHandler struct {
	Codec     *Codec
	StateTree *StateTree
}

func (h *InformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL != nil && r.URL.Path != "/inform" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	msg, err := h.Codec.Unmarshal(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	pl := NewInformWrapper()
	copy(pl.MacAddr, msg.MacAddr)
	pl.SetEncrypted(true)

	// TODO: compare current state to tree and update

	res, err := h.Codec.Marshal(pl)
	if err != nil {
		http.Error(w, "Server Error", 500)
		return
	}

	fmt.Fprintf(w, "%s", res)
}

func NewServer(handler *InformHandler) *http.Server {
	return &http.Server{
		Addr:    ":6080",
		Handler: Log(handler),
	}
}
