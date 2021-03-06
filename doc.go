/*
Package healthz is a status reporter library for microservices. It models the
system as a hierarchicy subcomponents, and calculate the overall health usng the
overall health of the direct subcomponets of the root component. The handler
exposes the health reports in multiple formats:

	* HTML Report (`/`)
	* JSON Report (`/json`)
	* Liveness (`/liveness`): HTTP 200 only if the overall health is > `Error`
	* Readiness (`/readiness`): HTTP 200 only if the overall health is >= `Normal`

Check the example, or run the demo:
	$GOPATH/src/github.com/cafebazaar/healthz/demo$ go run main.go
*/
package healthz
