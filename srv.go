package dnsutils

import (
	"math/rand"
	"net"
	"sort"
)

type SRVRecords []*net.SRV

func (s SRVRecords) Len() int {
	return len(s)
}
func (s SRVRecords) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SRVRecords) Less(i, j int) bool {
	// Order on priority and then on weighted random
	if s[i].Priority == s[j].Priority {

	}
	return s[i].Priority < s[j].Priority
}

// OrderedSRV returns a sorted int slice for the entries in the map. To use in the correct order:
//
// count, orderedSRV, err := OrderedSRV(addrs)
// i := 1
// for  i <= count {
//   srv := orderedSRV[i]
//   // Do something such as dial this SRV. If fails move on the the next
//   i += 1
// }
func OrderedSRV(service, proto, name string) (int, map[int]*net.SRV, error) {
	_, addrs, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return 0, make(map[int]*net.SRV), err
	}
	index, os := orderSRV(addrs)
	return index, os, nil
}

func orderSRV(addrs []*net.SRV) (int, map[int]*net.SRV) {
	// Initialise the ordered map
	var o int
	osrv := make(map[int]*net.SRV)

	prioMap := make(map[int][]*net.SRV)
	for _, srv := range addrs {
		prioMap[int(srv.Priority)] = append(prioMap[int(srv.Priority)], srv)
	}

	priorities := make([]int, 0)
	for p, _ := range prioMap {
		priorities = append(priorities, p)
	}
	sort.Ints(priorities)

	var count int
	sort.Ints(priorities)
	for _, p := range priorities {
		tos := weightedOrder(prioMap[p])
		for i, s := range tos {
			count = o + i
			osrv[count] = s
		}
		o += len(tos)
	}
	return count, osrv
}

func weightedOrder(srvs []*net.SRV) map[int]*net.SRV {
	// Get the total weight
	var tw int
	for _, s := range srvs {
		tw += int(s.Weight)
	}

	// Initialise the ordered map
	o := 1
	osrv := make(map[int]*net.SRV)

	// Whilst there are still entries to be ordered
	l := len(srvs)
	for l > 0 {
		i := rand.Intn(l)
		s := srvs[i]
		var rw int
		if tw > 0 {
			// Greater the weight the more likely this will be zero or less
			rw = rand.Intn(tw) - int(s.Weight)
		}
		if rw <= 0 {
			// Put entry in position
			osrv[o] = s
			if len(srvs) > 1 {
				// Remove the entry from the source slice by swapping with the last entry and truncating
				srvs[len(srvs)-1], srvs[i] = srvs[i], srvs[len(srvs)-1]
				srvs = srvs[:len(srvs)-1]
				l = len(srvs)
			} else {
				l = 0
			}
			o += 1
			tw = tw - int(s.Weight)
		}
	}
	return osrv
}
