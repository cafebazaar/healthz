package healthz_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cafebazaar/healthz"
)

func NewComponent(args string, h healthz.Handler) {
}

func ExampleNewHandler() {
	// In the main package
	healthzHandler := healthz.NewHandler("Application V1.0.0", true)
	healthzHandler.RegisterComponent("main", healthz.Major)
	healthzServer := &http.Server{
		Addr:    "127.0.0.1:8090",
		Handler: healthzHandler,
	}
	go func() {
		err := healthzServer.ListenAndServe()
		if err != nil {
			log.Fatalln("Error while healthzServer.ListenAndServe():", err)
		}
	}()

	fmt.Println("Overall Health:", healthzHandler.OverallHealth())

	// Pass the handler to the components
	NewComponent("component-args...", healthzHandler)

	// Mark "main" as ready when all the components are initialized
	healthzHandler.SetHealth("main", healthz.Normal)
	fmt.Println("Overall Health:", healthzHandler.OverallHealth())
	// Output: Overall Health: 0
	// Overall Health: 1
}
