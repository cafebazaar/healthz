package main

import (
	"log"
	"net/http"

	"github.com/cafebazaar/healthz"
)

const addr = "127.0.0.1:8090"

func main() {
	h := healthz.NewHandler("Demo (v1.0.0)", true)

	h.RegisterComponent("component-major-redundant", healthz.Major)
	h.SetHealth("component-major-redundant", healthz.Redundant)
	h.RegisterComponent("component-major-warning", healthz.Major)
	h.SetHealth("component-major-warning", healthz.Warning)
	h.RegisterComponent("component-unspecified-warning", healthz.Unspecified)
	h.SetHealth("component-unspecified-warning", healthz.Warning)
	h.RegisterComponent("component-minor-error", healthz.Minor)
	h.SetHealth("component-minor-error", healthz.Error)

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
