package healthz

import "testing"

func TestSetHealth(t *testing.T) {
	h := NewHandler("Test", true).(*handler)
	h.SetHealth("test", Warning)
	c, found := h.components["test"]
	if !found {
		t.Fatal("component test not found")
	}
	if c.Health != Warning {
		t.Fatal("unexpected health:", c.Health)
	}
}

func TestRegister(t *testing.T) {
	h := NewHandler("Test", true).(*handler)
	h.RegisterComponent("test", Major)
	c, found := h.components["test"]
	if !found {
		t.Fatal("component test not found")
	}
	if c.Severity != Major {
		t.Fatal("unexpected severity:", c.Severity)
	}

	h.UnregisterComponent("test")
	_, found = h.components["test"]
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
	h := NewHandler("Test", true).(*handler)
	if o := h.OverallHealth(); o != Unknown {
		t.Fatal("unexpected OverallHealth:", o)
	}

	h.RegisterComponent("major", Major)
	h.RegisterComponent("unspecified", Unspecified)
	h.RegisterComponent("minor", Minor)
	if o := h.OverallHealth(); o != Unknown {
		t.Fatal("unexpected OverallHealth:", o)
	}

	for i, tCase := range overallHealthTestCases {
		h.SetHealth("major", tCase.in[0])
		h.SetHealth("unspecified", tCase.in[1])
		h.SetHealth("minor", tCase.in[2])
		if o := h.OverallHealth(); o != tCase.out {
			t.Fatalf("case #%d: unexpected OverallHealth: %v", i, o)
		}
	}

}
