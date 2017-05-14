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
	for _, c := range c.Subcomponents {
		componentHealth := c.OverallHealth()
		if c.Severity >= Major && componentHealth < res {
			res = componentHealth
		} else if c.Severity == Unspecified && (componentHealth+1) < res {
			res = componentHealth + 1
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
	rc.OverallHealth = Redundant
	for name, c := range c.Subcomponents {
		subcomponentRC := c.reportComponents()
		subcomponentRC.Name = name
		rc.Subcomponents = append(rc.Subcomponents, subcomponentRC)

		componentHealth := subcomponentRC.OverallHealth
		if c.Severity >= Major && componentHealth < rc.OverallHealth {
			rc.OverallHealth = componentHealth
		} else if c.Severity == Unspecified && (componentHealth+1) < rc.OverallHealth {
			rc.OverallHealth = componentHealth + 1
		}
	}
	sort.Sort(rc)
	return rc
}
