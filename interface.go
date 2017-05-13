package healthz

import "net/http"

// Severity specifies the seriousness of a component
type Severity int8

// Health specifies if a component is ready or not
type Health int8

const (
	// Major means this component's failour disruptes some major
	// functionalities of the system.
	Major Severity = 1

	// Unspecified is the default value of Severity
	Unspecified Severity = 0

	// Minor means this component's failour causes non-critical loss of some
	// functionalities of the system.
	Minor Severity = -1

	// Redundant means there are multiple instances of this components ready to serve
	Redundant Health = 2

	// Normal means there is at least one instance of this component ready to serve
	Normal Health = 1

	// Unknown is the default value of Health
	Unknown Health = 0

	// Warning means there are some problems with this component
	Warning Health = -1

	// Error means that this components is unable to serve as expected
	Error Health = -2
)

var (
	healthToTitle   = make(map[Health]string)
	severityToTitle = make(map[Severity]string)
)

func init() {
	healthToTitle[Redundant] = "Redundant"
	healthToTitle[Normal] = "Normal"
	healthToTitle[Unknown] = "Unknown"
	healthToTitle[Warning] = "Warning"
	healthToTitle[Error] = "Error"

	severityToTitle[Major] = "Major"
	severityToTitle[Minor] = "Minor"
	severityToTitle[Unspecified] = "Unspecified"
}

// Handler is a http handler which is notified about health level of the
// different components in the system, and reports them through http
type Handler interface {

	// ServeHTTP is implemented to return reports as expected
	ServeHTTP(http.ResponseWriter, *http.Request)

	// RegisterComponent creates the component if it wasn't registered before,
	// and sets the severity level of the component to the given value
	RegisterComponent(name string, severity Severity)

	// UnregisterComponent removes the component from the report, and the
	// calculation of the `OverallHealth`
	UnregisterComponent(name string)

	// SetHealth sets the health level of the specified component
	SetHealth(component string, health Health)

	// OverallHealth of a system is minimum value of these two:
	//	* Minimum health level of all the components with severity=Major
	//	* 1 + Minimum health level of all the components with severity=Unspecified
	OverallHealth() Health
}
