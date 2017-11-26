package dnsutils

import (
	"net"
	"testing"
)

func TestOrderSRV(t *testing.T) {
	srv11 := net.SRV{
		Target:   "t11",
		Port:     1234,
		Priority: 1,
		Weight:   100,
	}
	srv12 := net.SRV{
		Target:   "t12",
		Port:     1234,
		Priority: 1,
		Weight:   101,
	}
	srv21 := net.SRV{
		Target:   "t21",
		Port:     1234,
		Priority: 2,
		Weight:   1,
	}

	addrs := []*net.SRV{
		&srv11, &srv21, &srv12,
	}
	index, orderedSRV := orderSRV(addrs)
	for _, i := range index {
		srv := orderedSRV[i]
		t.Logf("PRIO: %d WEIGHT: %d TARGET: %s", srv.Priority, srv.Weight, srv.Target)
	}
}
