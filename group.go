package healthz

import (
	"sort"
	"sync"
)

type componentState uint8

const (
	componentStateUnknown componentState = iota
	componentStateStandalone
	componentStateGroup
)

type componentGroup struct {
	Severity      Severity
	Health        Health
	Subcomponents map[string]*componentGroup

	state componentState
	mutex sync.RWMutex
}

func newComponentGroup(severity Severity, health Health) *componentGroup {
	return &componentGroup{
		Severity:      severity,
		Health:        health,
		Subcomponents: make(map[string]*componentGroup),
		state:         componentStateUnknown,
	}
}

func (c *componentGroup) RegisterSubcomponent(name string, severity Severity) ComponentGroup {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.state == componentStateStandalone {
		panic("You should avoid calling both SetGroupHealth and RegisterSubcomponent on same ComponentGroup")
	}
	c.state = componentStateGroup
	c.Subcomponents[name] = newComponentGroup(severity, Unknown)
	return c.Subcomponents[name]
}

func (c *componentGroup) UnregisterSubcomponent(name string) {
	c.mutex.Lock()
	delete(c.Subcomponents, name)
	c.mutex.Unlock()
}

func (c *componentGroup) SetGroupHealth(health Health) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.state == componentStateGroup {
		panic("You should avoid calling both SetGroupHealth and RegisterSubcomponent on same ComponentGroup")
	}
	c.state = componentStateStandalone
	c.Health = health
}

func (c *componentGroup) OverallHealth() Health {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.state == componentStateStandalone {
		return c.Health
	}
	if len(c.Subcomponents) == 0 {
		return Unknown
	}
	res := Redundant
	var majorIsInGroup, unspecifiedIsInGroup bool
	for _, c := range c.Subcomponents {
		if c.Severity >= Major {
			majorIsInGroup = true
			componentHealth := c.OverallHealth()
			if res > componentHealth {
				res = componentHealth
			}
		}
	}
	if !majorIsInGroup && res > Normal {
		res = Normal
	}
	for _, c := range c.Subcomponents {
		if c.Severity == Unspecified {
			unspecifiedIsInGroup = true
			componentHealth := c.OverallHealth()
			if majorIsInGroup || componentHealth < Normal {
				componentHealth++
			}
			if res > componentHealth {
				res = componentHealth
			}
		}
	}
	if !majorIsInGroup && !unspecifiedIsInGroup {
		res = Unknown
	}
	return res
}

func (c *componentGroup) GroupReport() *GroupReport {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	rc := &GroupReport{
		Severity:      c.Severity,
		Subcomponents: make([]*GroupReport, 0, len(c.Subcomponents)),
	}
	if c.state == componentStateStandalone || len(c.Subcomponents) == 0 {
		rc.OverallHealth = c.Health
		return rc
	}
	// TODO: how can we refactor this algorithm which also appears in OverallHealth()?
	rc.OverallHealth = Redundant
	var majorIsInGroup, unspecifiedIsInGroup bool
	for name, c := range c.Subcomponents {
		subcomponentRC := c.GroupReport()
		subcomponentRC.Name = name
		rc.Subcomponents = append(rc.Subcomponents, subcomponentRC)
		if c.Severity >= Major {
			majorIsInGroup = true
			if rc.OverallHealth > subcomponentRC.OverallHealth {
				rc.OverallHealth = subcomponentRC.OverallHealth
			}

		}
	}
	if !majorIsInGroup && rc.OverallHealth > Normal {
		rc.OverallHealth = Normal
	}
	for _, subcomponentRC := range rc.Subcomponents {
		if subcomponentRC.Severity == Unspecified {
			unspecifiedIsInGroup = true
			componentHealth := subcomponentRC.OverallHealth
			if majorIsInGroup || componentHealth < Normal {
				componentHealth++
			}
			if rc.OverallHealth > componentHealth {
				rc.OverallHealth = componentHealth
			}
		}
	}
	if !majorIsInGroup && !unspecifiedIsInGroup {
		rc.OverallHealth = Unknown
	}
	sort.Sort(rc)
	return rc
}
