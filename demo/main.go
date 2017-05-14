package main

import (
	"log"
	"net/http"

	"github.com/cafebazaar/healthz"
)

const addr = "127.0.0.1:8090"

func main() {
	h := healthz.NewHandler("Demo (v1.0.0)", true)

	componentMajorRedundant := h.RegisterSubcomponent("component-major-redundant", healthz.Major)
	componentMajorRedundant.SetGroupHealth(healthz.Redundant)
	componentMajorWarning := h.RegisterSubcomponent("component-major-warning", healthz.Major)
	componentMajorWarning.SetGroupHealth(healthz.Warning)
	componentUnspecifiedWarning := h.RegisterSubcomponent("component-unspecified-warning", healthz.Unspecified)
	componentUnspecifiedWarning.SetGroupHealth(healthz.Warning)
	componentMinorError := h.RegisterSubcomponent("component-minor-error", healthz.Minor)
	componentMinorError.SetGroupHealth(healthz.Error)

	componentComplex := h.RegisterSubcomponent("component-complex", healthz.Unspecified)
	subcomponent1 := componentComplex.RegisterSubcomponent("subcomponent1", healthz.Major)
	subcomponent1.SetGroupHealth(healthz.Redundant)
	subcomponent2 := componentComplex.RegisterSubcomponent("subcomponent2", healthz.Unspecified)
	subcomponent2.SetGroupHealth(healthz.Normal)
	subcomponent3 := componentComplex.RegisterSubcomponent("subcomponent3", healthz.Minor)
	subcomponent3.SetGroupHealth(healthz.Error)

	healthzServer := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	log.Printf("http://%s/\n", addr)
	err := healthzServer.ListenAndServe()
	if err != nil {
		log.Fatalln("Error while healthzServer.ListenAndServe():", err)
	}
}
