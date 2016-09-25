package inform

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		t := time.Now().Format("02/Jan/2006 15:04:05")
		log.Printf("%s - - [%s] \"%s %s %s\"", r.RemoteAddr, t, r.Method, r.URL, r.Proto)
	})
}

type Device struct {
	Initialized  bool
	CurrentState bool
	DesiredState bool
	*sync.Mutex
}

type InformHandler struct {
	Codec *Codec
	ports map[string]map[int]*Device
	queue map[string]*list.List
	*sync.RWMutex
}

func NewInformHandler(c *Codec) *InformHandler {
	return &InformHandler{
		Codec:   c,
		ports:   make(map[string]map[int]*Device),
		queue:   make(map[string]*list.List),
		RWMutex: &sync.RWMutex{},
	}
}

func (h *InformHandler) AddPort(dev string, port int) {
	h.Lock()
	defer h.Unlock()

	_, ok := h.ports[dev]
	if !ok {
		h.ports[dev] = make(map[int]*Device)
	}

	_, ok = h.queue[dev]
	if !ok {
		log.Printf("Adding queue for %s", dev)
		h.queue[dev] = list.New()
	}

	log.Printf("Adding %s port %d", dev, port)
	h.ports[dev][port] = &Device{
		Mutex: &sync.Mutex{},
	}
}

func (h *InformHandler) getPort(dev string, port int) (*Device, error) {
	h.RLock()
	defer h.RUnlock()

	_, ok := h.ports[dev]
	if !ok {
		return nil, errors.New("No device found")
	}

	p, ok := h.ports[dev][port]
	if !ok {
		return nil, errors.New("No port found")
	}

	return p, nil
}

func (h *InformHandler) SetState(dev string, port int, state bool) error {
	p, err := h.getPort(dev, port)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	log.Printf("Set state to %t for %s port %d", state, dev, port)
	p.DesiredState = state
	return nil
}

func (h *InformHandler) buildCommands(dev string, pl *DeviceMessage) error {
	for _, o := range pl.Outputs {
		ds, err := h.getPort(dev, o.Port)
		if err != nil {
			return err
		}
		ds.Lock()

		// Get initial state
		if !ds.Initialized {
			ds.CurrentState = o.OutputState
			ds.Initialized = true
			return nil
		}

		// State didn't change at the sensor
		if ds.CurrentState == o.OutputState {
			if ds.DesiredState != o.OutputState {
				log.Printf("Toggle state %t for %s port %d", ds.DesiredState, dev, o.Port)
				// Generate change command
				// TODO: Don't lock the whole handler
				h.Lock()
				h.queue[dev].PushFront(NewOutputCommand(o.Port, ds.DesiredState, 0))
				h.Unlock()
			}
		} else { // Sensor caused the change, leave it alone
			log.Printf("Sensor state changed %s port %d", dev, o.Port)
			ds.DesiredState = o.OutputState
		}

		ds.CurrentState = o.OutputState
		ds.Unlock() // Don't hold the lock the entire loop
	}

	return nil
}

func (h *InformHandler) pop(dev string) *CommandMessage {
	// TODO: Don't lock the whole handler
	h.Lock()
	defer h.Unlock()

	q, ok := h.queue[dev]
	if !ok {
		log.Printf("No queue for %s", dev)
		return nil
	}

	e := q.Front()
	if e != nil {
		h.queue[dev].Remove(e)
		cmd := e.Value.(*CommandMessage)
		cmd.Freshen()
		return cmd
	}

	return nil
}

func (h *InformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := h.Codec.Unmarshal(r.Body)
	if err != nil {
		log.Printf("Unmarshal message: %s", err.Error())
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	pl, err := msg.UnmarshalPayload()
	if err != nil {
		log.Printf("Unmarshal payload: %s", err.Error())
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	dev := msg.FormattedMac()
	ret := NewInformWrapperResponse(msg)
	log.Printf("Inform from %s", dev)

	// Send a command until the queue is empty
	if cmd := h.pop(dev); cmd != nil {
		ret.UpdatePayload(cmd)
	} else {
		// Update internal state vs reality
		if err := h.buildCommands(msg.FormattedMac(), pl); err != nil {
			http.Error(w, "500 Server Error", 500)
			return
		}

		// If that generated a command send it
		if cmd = h.pop(dev); cmd != nil {
			ret.UpdatePayload(cmd)
		} else {
			// Otherwise noop
			ret.UpdatePayload(NewNoop(10))
		}
	}

	res, err := h.Codec.Marshal(ret)
	if err != nil {
		http.Error(w, "500 Server Error", 500)
		return
	}

	fmt.Fprintf(w, "%s", res)
}

// Create a new server, returns the mux so users can add other methods (for
// example, if they want to share a process to build a console that also
// accepts informs)
func NewServer(handler *InformHandler) (*http.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	mux.Handle("/inform", handler)

	return &http.Server{
		Addr:    ":6080",
		Handler: Log(mux),
	}, mux
}
