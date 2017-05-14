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
	var res Health = 100
	for _, c := range c.Subcomponents {
		if c.Severity >= Major {
			componentHealth := c.OverallHealth()
			if res > componentHealth {
				res = componentHealth
			}
		}
	}
	if res == 100 {
		res = Unknown
	}
	for _, c := range c.Subcomponents {
		if c.Severity == Unspecified {
			componentHealth := c.OverallHealth()
			if res > (componentHealth + 1) {
				res = componentHealth + 1
			}
		}
	}
	return res
}

func (c *componentGroup) reportComponents() *reportComponents {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	rc := &reportComponents{
		Severity:      c.Severity,
		Subcomponents: make([]*reportComponents, 0, len(c.Subcomponents)),
	}
	if c.state == componentStateStandalone || len(c.Subcomponents) == 0 {
		rc.OverallHealth = c.Health
		return rc
	}
	for name, c := range c.Subcomponents {
		subcomponentRC := c.reportComponents()
		subcomponentRC.Name = name
		rc.Subcomponents = append(rc.Subcomponents, subcomponentRC)
	}
	// TODO: how can we refactor this algorithm which also appears in OverallHealth()?
	rc.OverallHealth = 100
	for _, c := range rc.Subcomponents {
		if c.Severity >= Major && rc.OverallHealth > c.OverallHealth {
			rc.OverallHealth = c.OverallHealth
		}
	}
	if rc.OverallHealth == 100 {
		rc.OverallHealth = Unknown
	}
	for _, c := range rc.Subcomponents {
		if c.Severity == Unspecified && rc.OverallHealth > (c.OverallHealth+1) {
			rc.OverallHealth = c.OverallHealth + 1
		}
	}
	sort.Sort(rc)
	return rc
}
