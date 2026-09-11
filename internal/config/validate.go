package config

import (
	"fmt"
	"regexp"
)

// суть валидейт проверять конфиг на внутренную согласованость
// надо вернуть первую ошибку или nil

func (c *Config) Validate() error {
	if c.Server.Listen == "" {
		return fmt.Errorf("server.listen is required")
	}
	if len(c.Upstreams) == 0 {
		return fmt.Errorf("at least one upstream is required")
	}

	for name, up := range c.Upstreams {
		if len(up.Targets) == 0 {
			return fmt.Errorf("upstream %q has no targets", name)
		}
	}

	seen := make(map[string]bool)
	for i, r := range c.Routes {
		if r.Name == "" {
			return fmt.Errorf("route #%d has empty name", i)
		}
		if seen[r.Name] {
			return fmt.Errorf("dublicate route name %q", r.Name)
		}
		seen[r.Name] = true

		if _, ok := c.Upstreams[r.Upstream]; !ok {
			return fmt.Errorf("route %q referensces unknown upstream %q", r.Name, r.Upstream)
		}
		if err := r.Match.Path.validate(); err != nil {
			return fmt.Errorf("route %q: %w", r.Name, err)
		}
	}

	if c.Fallback != nil {
		if _, ok := c.Upstreams[c.Fallback.Upstream]; !ok {
			return fmt.Errorf("fb references uknown upstream %q", c.Fallback.Upstream)
		}
	}
	return  nil
}

func (p PathMatch) validate() error {
	count := 0

	// ---------------------------------
	if p.Exact  != "" { count++ }
	if p.Prefix != "" { count++ }
	if p.Regex  != "" { count++ }

	if count == 0 { return fmt.Errorf("path match is empty") }
	if count >  1 { return fmt.Errorf("path match must have exactly one")}
	// ---------------------------------

	if p.Regex != "" {
		if _, err := regexp.Compile(p.Regex); err != nil {
			return fmt.Errorf("invalid regex %q: %w", p.Regex, err)
		}
	}

	return  nil
}