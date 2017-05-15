package healthz

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

// GroupReport is a copied status of a group and its subcomponents
type GroupReport struct {
	Name          string
	Severity      Severity
	OverallHealth Health
	Subcomponents []*GroupReport `json:",omitempty"`
}

// ComponentGroup represents a component or a group of components. You can set
// health level of the group by calling `SetGroupHealth`, OR, by creating
// subcomponents. Note that you can't mix these two mechanism, it will cause
// panic!
type ComponentGroup interface {
	// SetGroupHealth sets the health level of the specified component.
	SetGroupHealth(health Health)

	// RegisterSubcomponent creates a subcomponent if it wasn't registered
	// before, and sets the severity level of the subcomponent to the given
	// value.
	RegisterSubcomponent(name string, severity Severity) ComponentGroup

	// UnregisterSubcomponent removes the subcomponent from the group, and the
	// calculation of the `OverallHealth`
	UnregisterSubcomponent(name string)

	// OverallHealth is the specified value set by `SetGroupHealth`, or if the
	// this instance contains one or more subcomponents, it's the minimum value
	// of these two:
	//	* Minimum health level of all the subcomponents with severity=Major
	//	* 1 + Minimum health level of all the components with severity=Unspecified
	// If no Major or Unspecified component is registered, OverallHealth
	// returns Unknown.
	// Otherwise, if no Major component is registered, the result will be
	// capped at Normal
	OverallHealth() Health

	// GroupReport copies the current status of the group and its subcomponents,
	// and returns the copied object
	GroupReport() *GroupReport
}
