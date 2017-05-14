package healthz

import (
	"net/http"
	"os"
	"time"
)

// Handler is a http handler which is notified about health level of the
// different components in the system, and reports them through http
type Handler struct {
	details          bool
	serviceSignature string
	startTime        time.Time
	hostname         string

	rootComponent *componentGroup

	mux *http.ServeMux
}

// RegisterSubcomponent creates a subcomponent if it wasn't registered before,
// and sets the severity level of the component to the given value
func (h *Handler) RegisterSubcomponent(name string, severity Severity) ComponentGroup {
	return h.rootComponent.RegisterSubcomponent(name, severity)
}

// UnregisterSubcomponent removes the subcomponent from the group, the reports,
// and the calculation of the `OverallHealth`
func (h *Handler) UnregisterSubcomponent(name string) {
	h.rootComponent.UnregisterSubcomponent(name)
}

// SetGroupHealth sets the health level of the root component.
// Note that you can't mix SetGroupHealth and RegisterSubcomponent
func (h *Handler) SetGroupHealth(health Health) {
	h.rootComponent.SetGroupHealth(health)
}

// OverallHealth is the overall health level of the handler. Its calculation
// is described in the doc of `ComponentGroup.OverallHealth`
func (h *Handler) OverallHealth() Health {
	return h.rootComponent.OverallHealth()
}

// ServeHTTP is implemented to response to the report requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// NewHandler creates a Handler and returns it. if details is true, these
// values are also included in the applicable reports:
//
//	* Uptime: since call of this function
//	* Service Signature: which is passed to this function
//	* Hostname: which is extracted using `os.Hostname()`
//	* Health level and Severity level for all the registered components
func NewHandler(serviceSignature string, details bool) *Handler {
	mux := http.NewServeMux()
	hostname, _ := os.Hostname()
	h := &Handler{
		details:          details,
		serviceSignature: serviceSignature,
		startTime:        time.Now(),
		hostname:         hostname,
		rootComponent:    newComponentGroup(Major, Unknown),

		mux: mux,
	}

	mux.HandleFunc("/", h.reportHTML)
	mux.HandleFunc("/liveness", h.reportLiveness)
	mux.HandleFunc("/readiness", h.reportReadiness)
	mux.HandleFunc("/min.css", h.reportMinCSS)
	mux.HandleFunc("/json", h.reportJSON)

	return h
}
