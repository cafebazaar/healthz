package healthz

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"time"
)

var (
	htmlTemplate *template.Template
)

func init() {
	var err error
	t := template.New("webpage").Funcs(map[string]interface{}{
		"HealthTitle": func(h Health) string {
			return healthToTitle[h]
		},
		"SeverityTitle": func(s Severity) string {
			return severityToTitle[s]
		},
	})
	htmlTemplate, err = t.Parse(tpl)
	if err != nil {
		panic(err)
	}
}

type reportComponents []*reportComponent
type reportComponent struct {
	Name     string
	Severity Severity
	Health   Health
}
type report struct {
	ServiceSignature string
	OverallHealth    Health
	Uptime           time.Duration    `json:",omitempty"`
	Hostname         string           `json:",omitempty"`
	Components       reportComponents `json:",omitempty"`
}

func (a reportComponents) Len() int      { return len(a) }
func (a reportComponents) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a reportComponents) Less(i, j int) bool {
	if a[i].Health == a[j].Health {
		return a[i].Severity > a[j].Severity
	}
	return a[i].Health < a[j].Health
}

func (h *handler) report() *report {
	rpt := &report{
		ServiceSignature: h.serviceSignature,
		OverallHealth:    h.OverallHealth(),
	}
	if h.details {
		uptime := time.Since(h.startTime)
		rpt.Uptime = uptime
		rpt.Hostname = h.hostname
		h.mutex.RLock()
		components := make([]*reportComponent, 0, len(h.components))
		for name, c := range h.components {
			components = append(components, &reportComponent{
				Name:     name,
				Health:   c.Health,
				Severity: c.Severity,
			})
		}
		h.mutex.RUnlock()
		rpt.Components = reportComponents(components)
		sort.Sort(rpt.Components)
	}
	return rpt
}

func (h *handler) reportLiveness(w http.ResponseWriter, r *http.Request) {
	overallHealth := h.OverallHealth()
	if overallHealth >= Unknown {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, healthToTitle[overallHealth], http.StatusServiceUnavailable)
	}
}

func (h *handler) reportReadiness(w http.ResponseWriter, r *http.Request) {
	overallHealth := h.OverallHealth()
	if overallHealth >= Normal {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, healthToTitle[overallHealth], http.StatusServiceUnavailable)
	}
}

func (h *handler) reportJSON(w http.ResponseWriter, r *http.Request) {
	rpt := h.report()
	jData, _ := json.Marshal(rpt)
	w.Header().Set("Overall-Health", healthToTitle[rpt.OverallHealth])
	w.Header().Set("Overall-Health-Code", fmt.Sprintf("%d", rpt.OverallHealth))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func (h *handler) reportHTML(w http.ResponseWriter, r *http.Request) {
	rpt := h.report()
	w.Header().Set("Overall-Health", healthToTitle[rpt.OverallHealth])
	w.Header().Set("Overall-Health-Code", fmt.Sprintf("%d", rpt.OverallHealth))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := htmlTemplate.Execute(w, rpt)
	if err != nil {
		w.Write([]byte("\n\n<h1>error while rendering the report:\n"))
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n</h1>"))
	}
}
