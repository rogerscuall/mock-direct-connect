package bgp

func checkASN(asn int) bool {
	if asn < 1 || asn > 2147483647 {
		return false
	}
	return true
}
