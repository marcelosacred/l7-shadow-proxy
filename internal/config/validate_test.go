package config

import "testing"

func TestPathMatchValidate(t *testing.T) {
	cases := []struct {
		name  string
		path  PathMatch
		wantErr bool
	}  {
		{"only exact", 	 PathMatch{Exact: "/api"},   		   false},
		{"only prefix",  PathMatch{Prefix: "/api/"},           false},
		{"only regex", 	 PathMatch{Regex: "^/v[0-9]+$"}, 	   false},
		{"empty", 		 PathMatch{}, 						   true},
		{"two at once",  PathMatch{Exact: "/a", Prefix: "/b"}, true},
		{"broken regex", PathMatch{Regex: "^([0-9"},           true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.path.validate()
			gotErr := err != nil

			if gotErr != tc.wantErr {
				t.Errorf("validate() error = %v, wantErr = %v", err, tc.wantErr)			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	valid := &Config{
		Server: ServerConfig{Listen: ":8080"},
		Upstreams: map[string]UpstreamConfig{
			"api": {Targets: []string{"http://localhost:8081"}},
		},
		Routes: []RouteConfig {
			{Name: "r1", Match: MatchConfig{Path: PathMatch{Prefix: "/"}}, Upstream: "api"},
		},
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}

	badRef := &Config{
		Server:    ServerConfig{Listen: ":8080"},
		Upstreams: map[string]UpstreamConfig{"api": {Targets: []string{"http://x"}}},
		Routes: []RouteConfig{
			{Name: "r1", Match: MatchConfig{Path: PathMatch{Prefix: "/"}}, Upstream: "nope"},
		},
	}
	if err := badRef.Validate(); err == nil {
			t.Errorf("expected error for unknown upstream, got nil")
		}
}