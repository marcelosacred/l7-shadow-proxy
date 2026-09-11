package config

import "time"


// КОРЕНЬ
type Config struct {
	Server 	  ServerConfig				`yaml:"server"`
	Metrics   MetricsConfig				`yaml:"metrics"`
	Upstreams map[string]UpstreamConfig `yaml:"upstreams"`
	Routes    []RouteConfig				`yaml:"routes"`
	Fallback  *FallbackConfig			`yaml"fallback"`
}

type ServerConfig struct {
	Listen 			  string	    `yaml:"listen"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	MaxHeaderBytes    int			`yaml:"max_header_bytes"`
	MaxBodyBytes      int64			`yaml:"max_body_bytes"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
}

type MetricsConfig struct {
	Listen string `yaml"listen"`
}

type UpstreamConfig struct {
	Targets 		[]string 	  `yaml:"targets"`
	Timeout 		time.Duration `yaml:"timeout"`
	MaxConnsPerHost int 		  `yaml:"max_conns_per_host"`
}

type RouteConfig struct {
	Name     string      `yaml:"name"`
	Match    MatchConfig `yaml:"match"`
	Upstream string      `yaml:"upstream"`
}

type MatchConfig struct {
	Path PathMatch `yaml:"path"`
}

type PathMatch struct {
	Exact  string `yaml:"exact"`
	Prefix string `yaml:"prefix"`
	Regex  string `yaml:"regex"`
}

type FallbackConfig struct {
	Upstream string `yaml:"upstream"`
}