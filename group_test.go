package healthz

import "testing"

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

var overallHealthTestCases = []struct {
	in  [3]Health
	out Health
}{
	{
		in:  [3]Health{Unknown, Unknown, Unknown},
		out: Unknown,
	},
	{
		in:  [3]Health{Unknown, Unknown, Redundant},
		out: Unknown,
	},
	{
		in:  [3]Health{Unknown, Redundant, Unknown},
		out: Unknown,
	},
	{
		in:  [3]Health{Redundant, Unknown, Unknown},
		out: Normal,
	},
	{
		in:  [3]Health{Redundant, Redundant, Unknown},
		out: Redundant,
	},
	{
		in:  [3]Health{Redundant, Normal, Unknown},
		out: Redundant,
	},
	{
		in:  [3]Health{Redundant, Redundant, Warning},
		out: Redundant,
	},
	{
		in:  [3]Health{Redundant, Error, Redundant},
		out: Warning,
	},
}

func TestOverallHealth(t *testing.T) {
	c := newComponentGroup(Major, Unknown)
	if o := c.OverallHealth(); o != Unknown {
		t.Fatal("unexpected OverallHealth:", o)
	}

	cs0 := c.RegisterSubcomponent("major", Major)
	cs1 := c.RegisterSubcomponent("unspecified", Unspecified)
	cs2 := c.RegisterSubcomponent("minor", Minor)
	if o := c.OverallHealth(); o != Unknown {
		t.Fatal("unexpected OverallHealth:", o)
	}

	for i, tCase := range overallHealthTestCases {
		cs0.SetGroupHealth(tCase.in[0])
		cs1.SetGroupHealth(tCase.in[1])
		cs2.SetGroupHealth(tCase.in[2])
		if o := c.OverallHealth(); o != tCase.out {
			t.Fatalf("case #%d: unexpected OverallHealth: %v", i, o)
		}
	}
}
