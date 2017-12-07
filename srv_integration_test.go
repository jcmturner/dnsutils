// +build integration
// To turn on this test use -tags=integration in go test command

package dnsutils

import "github.com/stretchr/testify/assert"

func TestResolveKDC(t *testing.T) {
	count, res, err := OrderedSRV("kerberos", "tcp", "test.gokrb5")
	if err != nil {
		t.Errorf("error resolving SRV DNS records: %v", err)
	}
	assert.Equal(t, 5, count, "Number of SRV records not as expected: %v", res)
	assert.Equal(t, count, len(res), "Map size does not match: %v", res)
	t.Logf("res: %v", res)
	expected := []string{
		"kdc.test.gokrb5:88",
		"kdc1a.test.gokrb5:88",
		"kdc2a.test.gokrb5:88",
		"kdc1b.test.gokrb5:88",
		"kdc2b.test.gokrb5:88",
	}
	for _, s := range expected {
		var found bool
		for _, v := range res {
			if s == v {
				found = true
				break
			}
		}
		assert.True(t, found, "Record %s not found in results", s)
	}
}
