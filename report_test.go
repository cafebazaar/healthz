package healthz

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var livenessAndReadinessReportsTestCases = []struct {
	in           Health
	livenessOut  int
	readinessOut int
}{
	{in: Redundant, livenessOut: http.StatusOK, readinessOut: http.StatusOK},
	{in: Normal, livenessOut: http.StatusOK, readinessOut: http.StatusOK},
	{in: Unknown, livenessOut: http.StatusOK, readinessOut: http.StatusServiceUnavailable},
	{in: Warning, livenessOut: http.StatusServiceUnavailable, readinessOut: http.StatusServiceUnavailable},
	{in: Error, livenessOut: http.StatusServiceUnavailable, readinessOut: http.StatusServiceUnavailable},
}

func TestLivenessAndReadinessReports(t *testing.T) {
	h := NewHandler("Application V1.0.0", true)

	for i, tCase := range livenessAndReadinessReportsTestCases {
		h.SetGroupHealth(tCase.in)
		{
			r := httptest.NewRequest("GET", "http://127.0.0.1/liveness", strings.NewReader(""))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tCase.livenessOut {
				t.Fatalf("case #%d: unexpected liveness response code: %d", i, w.Code)
			}
		}
		{
			r := httptest.NewRequest("GET", "http://127.0.0.1/readiness", strings.NewReader(""))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tCase.readinessOut {
				t.Fatalf("case #%d: unexpected readiness response code: %d", i, w.Code)
			}
		}
	}

}

//TODO: test reportJSON
