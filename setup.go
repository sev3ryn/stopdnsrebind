package stopdnsrebind

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"inet.af/netaddr"
)

func init() { plugin.Register("stopdnsrebind", setup) }

func setup(c *caddy.Controller) error {
	publicNets, err := parse(c)
	// parsing err
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Stopdnsrebind{Next: next, PublicNets: publicNets}
	})

	return nil
}

func parse(c *caddy.Controller) (*netaddr.IPSet, error) {
	var b netaddr.IPSetBuilder
	for c.Next() {
		for c.NextBlock() {
			if c.Val() != "public_nets" {
				return nil, plugin.Error("stopdnsrebind", c.Err("only public_nets operation is supported"))
			}

			for _, d := range c.RemainingArgs() {

				r, err := netaddr.ParseIPPrefix(d)
				if err != nil {
					return nil, plugin.Error("stopdnsrebind", c.Errf("%s is not a valid ip range: %s", d, err))
				}

				b.AddPrefix(r)
			}
		}
	}

	return b.IPSet()
}
