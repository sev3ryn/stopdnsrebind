package stopdnsrebind

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/miekg/dns"
	"inet.af/netaddr"
)

type Stopdnsrebind struct {
	Next       plugin.Handler
	PublicNets *netaddr.IPSet
}

// ServeDNS implements the plugin.Handler interface.
func (a Stopdnsrebind) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	nw := nonwriter.New(w)

	rcode, err := plugin.NextOrFailure(a.Name(), a.Next, ctx, nw, r)
	if err != nil {
		return rcode, err
	}

	for _, ans := range nw.Msg.Answer {
		var ip net.IP

		switch ans.Header().Rrtype {
		case dns.TypeA:
			ip = ans.(*dns.A).A
		case dns.TypeAAAA:
			ip = ans.(*dns.AAAA).AAAA
		default:
			// we only care about A and AAA
			continue
		}

		// check if private
		if a.isForbidden(ip) {
			m := new(dns.Msg)
			m.SetRcode(r, dns.RcodeRefused)
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}
	}

	w.WriteMsg(nw.Msg)

	return 0, nil
}

func (a *Stopdnsrebind) isForbidden(ip net.IP) bool {
	ipaddr, ok := netaddr.FromStdIP(ip)
	if !ok {
		return true
	}

	if ipaddr.IsLoopback() || ipaddr.IsPrivate() || ipaddr.IsUnspecified() || a.PublicNets.Contains(ipaddr) {
		return true
	}

	return false
}

// Name implements the Handler interface.
func (a Stopdnsrebind) Name() string { return "stopdnsrebind" }
