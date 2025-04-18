package cors

import "github.com/rednek21/SimpleChat/config"

type OriginChecker interface {
	IsAllowed(origin string) bool
}

type Checker struct {
	allowed map[string]struct{}
}

func NewChecker(cfg config.Cors) *Checker {
	m := make(map[string]struct{}, len(cfg.Origins))
	for _, origin := range cfg.Origins {
		m[origin] = struct{}{}
	}
	return &Checker{allowed: m}
}

func (c *Checker) IsAllowed(origin string) bool {
	_, ok := c.allowed[origin]
	return ok
}
