package stopdnsrebind

import (
	"testing"

	"github.com/coredns/caddy"
)

func Test_setup(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			"public net",
			`stopdnsrebind {
				public_nets 8.8.8.8/24
			}`,
			false,
		},
		{
			"public net multiple",
			`stopdnsrebind {
				public_nets 8.8.8.0/10 9.9.9.0/16
			}`,
			false,
		},
		{
			"non supported op",
			`stopdnsrebind {
				anything internal.example.org.
			}`,
			true,
		},
		{
			"not a valid range",
			`stopdnsrebind {
				public_nets 8.8.8.z
			}`,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := caddy.NewTestController("dns", tt.config)
			if err := setup(ctr); (err != nil) != tt.wantErr {
				t.Errorf("setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
