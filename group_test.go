package healthz

import (
	"fmt"
	"testing"
)

func TestComponentGroupSetGroupHealth(t *testing.T) {
	c := newComponentGroup(Major, Unknown)
	if c.Health != Unknown {
		t.Fatal("unexpected health:", c.Health)
	}
	if c.OverallHealth() != Unknown {
		t.Fatal("unexpected OverallHealth():", c.Health)
	}
	c.SetGroupHealth(Warning)
	if c.Health != Warning {
		t.Fatal("unexpected health:", c.Health)
	}
	if c.OverallHealth() != Warning {
		t.Fatal("unexpected OverallHealth():", c.Health)
	}
}

func TestRegister(t *testing.T) {
	c := newComponentGroup(Minor, Unknown)
	c.RegisterSubcomponent("test", Major)
	sc, found := c.Subcomponents["test"]
	if !found {
		t.Fatal("component test not found")
	}
	if sc.Severity != Major {
		t.Fatal("unexpected severity:", c.Severity)
	}

	c.UnregisterSubcomponent("test")
	_, found = c.Subcomponents["test"]
	if found {
		t.Fatal("unregistered component test was found")
	}
}

type hs struct {
	h Health
	s Severity
}

var overallHealthTestCases = []struct {
	in  []hs
	out Health
}{
	{
		in:  []hs{{Unknown, Major}, {Unknown, Unspecified}, {Unknown, Minor}},
		out: Unknown,
	},
	{
		in:  []hs{{Unknown, Major}, {Unknown, Unspecified}, {Redundant, Minor}},
		out: Unknown,
	},
	{
		in:  []hs{{Unknown, Major}, {Redundant, Unspecified}, {Unknown, Minor}},
		out: Unknown,
	},
	{
		in:  []hs{{Redundant, Major}, {Unknown, Unspecified}, {Unknown, Minor}},
		out: Normal,
	},
	{
		in:  []hs{{Redundant, Major}, {Redundant, Unspecified}, {Unknown, Minor}},
		out: Redundant,
	},
	{
		in:  []hs{{Redundant, Major}, {Normal, Unspecified}, {Unknown, Minor}},
		out: Redundant,
	},
	{
		in:  []hs{{Redundant, Major}, {Redundant, Unspecified}, {Warning, Minor}},
		out: Redundant,
	},
	{
		in:  []hs{{Redundant, Major}, {Error, Unspecified}, {Redundant, Minor}},
		out: Warning,
	},
	{
		in:  []hs{{Unknown, Unspecified}},
		out: Unknown,
	},
	{
		in:  []hs{{Error, Minor}},
		out: Unknown,
	},
}

func TestOverallHealth(t *testing.T) {
	{
		c := newComponentGroup(Major, Unknown)
		if o := c.OverallHealth(); o != Unknown {
			t.Fatal("unexpected OverallHealth:", o)
		}

		c.RegisterSubcomponent("major", Major)
		c.RegisterSubcomponent("unspecified", Unspecified)
		c.RegisterSubcomponent("minor", Minor)
		if o := c.OverallHealth(); o != Unknown {
			t.Fatal("unexpected OverallHealth:", o)
		}
	}

	for i, tCase := range overallHealthTestCases {
		c := newComponentGroup(Major, Unknown)
		if o := c.OverallHealth(); o != Unknown {
			t.Fatal("unexpected OverallHealth:", o)
		}
		for i, p := range tCase.in {
			c.RegisterSubcomponent(fmt.Sprintf("c#%d", i), p.s).SetGroupHealth(p.h)
		}
		if o := c.OverallHealth(); o != tCase.out {
			t.Fatalf("case #%d: unexpected OverallHealth: %v", i, o)
		}
	}
}
