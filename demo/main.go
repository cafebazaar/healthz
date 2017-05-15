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
	subcomponent11 := componentComplex.RegisterSubcomponent("subcomponent11", healthz.Major)
	subcomponent11.SetGroupHealth(healthz.Redundant)
	subcomponent12 := componentComplex.RegisterSubcomponent("subcomponent12", healthz.Unspecified)
	subcomponent12.SetGroupHealth(healthz.Normal)
	subcomponent13 := componentComplex.RegisterSubcomponent("subcomponent13", healthz.Minor)
	subcomponent13.SetGroupHealth(healthz.Error)

	componentComplexNoMajor := h.RegisterSubcomponent("component-complex-no-major", healthz.Unspecified)
	subcomponent21 := componentComplexNoMajor.RegisterSubcomponent("subcomponent21", healthz.Unspecified)
	subcomponent21.SetGroupHealth(healthz.Redundant)

	componentComplexNoMajorNoUnspecified := h.RegisterSubcomponent("component-complex-no-major-no-unspecified", healthz.Unspecified)
	subcomponent31 := componentComplexNoMajorNoUnspecified.RegisterSubcomponent("subcomponent31", healthz.Minor)
	subcomponent31.SetGroupHealth(healthz.Redundant)

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
