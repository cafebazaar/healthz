package healthz

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

type reportComponents struct {
	Name          string
	Severity      Severity
	OverallHealth Health
	Subcomponents []*reportComponents `json:",omitempty"`
}
type report struct {
	ServiceSignature string
	Uptime           time.Duration `json:",omitempty"`
	Hostname         string        `json:",omitempty"`
	Root             *reportComponents
}

func (rc *reportComponents) Len() int {
	return len(rc.Subcomponents)
}
func (rc *reportComponents) Swap(i, j int) {
	rc.Subcomponents[i], rc.Subcomponents[j] = rc.Subcomponents[j], rc.Subcomponents[i]
}
func (rc *reportComponents) Less(i, j int) bool {
	if rc.Subcomponents[i].OverallHealth == rc.Subcomponents[j].OverallHealth {
		return rc.Subcomponents[i].Severity > rc.Subcomponents[j].Severity
	}
	return rc.Subcomponents[i].OverallHealth < rc.Subcomponents[j].OverallHealth
}

func (h *Handler) report() *report {
	rpt := &report{
		ServiceSignature: h.serviceSignature,
	}
	if h.details {
		uptime := time.Since(h.startTime)
		rpt.Uptime = uptime
		rpt.Hostname = h.hostname
		rpt.Root = h.rootComponent.reportComponents()
	}
	return rpt
}

func (h *Handler) reportLiveness(w http.ResponseWriter, r *http.Request) {
	overallHealth := h.OverallHealth()
	if overallHealth >= Unknown {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, healthToTitle[overallHealth], http.StatusServiceUnavailable)
	}
}

func (h *Handler) reportReadiness(w http.ResponseWriter, r *http.Request) {
	overallHealth := h.OverallHealth()
	if overallHealth >= Normal {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, healthToTitle[overallHealth], http.StatusServiceUnavailable)
	}
}

func (h *Handler) reportJSON(w http.ResponseWriter, r *http.Request) {
	rpt := h.report()
	jData, _ := json.Marshal(rpt)
	w.Header().Set("Overall-Health", healthToTitle[rpt.Root.OverallHealth])
	w.Header().Set("Overall-Health-Code", fmt.Sprintf("%d", rpt.Root.OverallHealth))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func (h *Handler) reportHTML(w http.ResponseWriter, r *http.Request) {
	rpt := h.report()
	w.Header().Set("Overall-Health", healthToTitle[rpt.Root.OverallHealth])
	w.Header().Set("Overall-Health-Code", fmt.Sprintf("%d", rpt.Root.OverallHealth))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := htmlTemplate.Execute(w, rpt)
	if err != nil {
		w.Write([]byte("\n\n<h1>error while rendering the report:\n"))
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n</h1>"))
	}
}
