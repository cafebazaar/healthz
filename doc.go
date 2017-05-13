/*
TODO

	healthzHandler := healthz.NewHandler("Application V1.0.0", true)
	healthzHandler.RegisterComponent("main", healthz.Critical)
	healthzServer := &http.Server{
		Addr:    "127.0.0.1:8090",
		Handler: healthzHandler,
	}
	go func() {
		err := healthzServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("Error while healthzServer.ListenAndServe():", err)
		}
	}()
	...
	healthzHandler.SetHealth(healthz.Normal)

*/
package healthz
