package dnsutils

import (
	"fmt"
	"math/rand"
	"net"
	"os"
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
// index, orderedSRV, err := OrderedSRV(addrs)
// for _, i := range index {
//   srv := orderedSRV[i]
//   // Do something such as dial this SRV. If fails move on the the next
// }
func OrderedSRV(service, proto, name string) ([]int, map[int]*net.SRV, error) {
	_, addrs, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return []int{}, make(map[int]*net.SRV), err
	}
	index, os := orderSRV(addrs)
	return index, os, nil
}

func orderSRV(addrs []*net.SRV) ([]int, map[int]*net.SRV) {
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

	index := make([]int, 0)
	sort.Ints(priorities)
	for _, p := range priorities {
		tos := weightedOrder(prioMap[p])
		for i, s := range tos {
			osrv[o+i] = s
			index = append(index, o+i)
		}
		o += len(tos)
	}
	sort.Ints(index)
	return index, osrv
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
		rw := rand.Intn(tw)
		for i, s := range srvs {
			// Greater the weight the more likely this will be zero or less
			//fmt.Fprintf(os.Stderr, "rw: %d\n", rw)
			rw = rw - int(s.Weight)
			//fmt.Fprintf(os.Stderr, "tw: %d, rw: %d, w: %d\n", tw, rw, int(s.Weight))
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
				break
			}
		}
	}
	return osrv
}
