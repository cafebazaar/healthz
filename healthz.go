package healthz

import (
	"net/http"
	"os"
	"sync"
	"time"
)

type component struct {
	Severity Severity
	Health   Health
}

type handler struct {
	details          bool
	serviceSignature string
	startTime        time.Time
	hostname         string
	components       map[string]*component

	mux   *http.ServeMux
	mutex sync.RWMutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) registerComponent(name string, severity Severity, health Health) {
	h.mutex.Lock()
	h.components[name] = &component{
		Severity: severity,
		Health:   health,
	}
	h.mutex.Unlock()
}

func (h *handler) RegisterComponent(name string, severity Severity) {
	h.registerComponent(name, severity, Unknown)
}

func (h *handler) UnregisterComponent(name string) {
	h.mutex.Lock()
	delete(h.components, name)
	h.mutex.Unlock()
}

func (h *handler) SetHealth(name string, health Health) {
	h.mutex.RLock()
	if c, found := h.components[name]; found {
		c.Health = health
		h.mutex.RUnlock()
		return
	}
	h.mutex.RUnlock()
	h.registerComponent(name, Unspecified, health)
}

func (h *handler) OverallHealth() Health {
	h.mutex.RLock()
	if len(h.components) == 0 {
		h.mutex.RUnlock()
		return Unknown
	}
	res := Redundant
	for _, c := range h.components {
		if c.Severity >= Major && c.Health < res {
			res = c.Health
		} else if c.Severity == Unspecified && (c.Health+1) < res {
			res = c.Health + 1
		}
	}
	h.mutex.RUnlock()
	return res
}

// NewHandler creates a Handler and returns it. if details is true, these
// values are also included in the applicable reports:
// * Uptime: since call of this function
// * Service Signature: which is passed to this function
// * Hostname: which is extracted using `os.Hostname()`
// * Health level and Severity level for all the registered components
func NewHandler(serviceSignature string, details bool) Handler {
	mux := http.NewServeMux()
	hostname, _ := os.Hostname()
	h := &handler{
		details:          details,
		serviceSignature: serviceSignature,
		startTime:        time.Now(),
		hostname:         hostname,
		components:       make(map[string]*component),

		mux: mux,
	}

	mux.HandleFunc("/", h.reportHTML)
	mux.HandleFunc("/liveness", h.reportLiveness)
	mux.HandleFunc("/readiness", h.reportReadiness)
	mux.HandleFunc("/min.css", h.reportMinCSS)
	mux.HandleFunc("/json", h.reportJSON)

	return h
}
