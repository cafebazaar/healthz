package healthz_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cafebazaar/healthz"
)

type testComponent struct {
	h healthz.ComponentGroup
}

func NewImportantComponent(args string, h healthz.ComponentGroup) *testComponent {
	h.SetGroupHealth(healthz.Warning)
	return &testComponent{h: h}
}

func (c *testComponent) Fix() {
	c.h.SetGroupHealth(healthz.Redundant)
}

// Example code, to be placed in the main of the application
func ExampleNewHandler() {
	healthzHandler := healthz.NewHandler("Application V1.0.0", true)
	mainHealth := healthzHandler.RegisterSubcomponent("main", healthz.Major)

	// Start Status Server
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

	// Create and pass `ComponentGroup`s to the components
	componentHealth := healthzHandler.RegisterSubcomponent("important-component", healthz.Major)
	component := NewImportantComponent("init-component-in-warning-mode", componentHealth)

	fmt.Println("important-component Health:", componentHealth.OverallHealth())
	fmt.Println("Overall Health:", healthzHandler.OverallHealth())

	component.Fix()
	fmt.Println("important-component Health:", componentHealth.OverallHealth())
	fmt.Println("Overall Health:", healthzHandler.OverallHealth())

	// Mark "main" as ready when all the components are initialized
	mainHealth.SetGroupHealth(healthz.Normal)
	fmt.Println("Overall Health:", healthzHandler.OverallHealth())

	// Output: Overall Health: 0
	// important-component Health: -1
	// Overall Health: -1
	// important-component Health: 2
	// Overall Health: 0
	// Overall Health: 1
}
