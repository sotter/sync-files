package base


func RSHash(str []byte) uint64 {
	b := 378551
	a := 63689
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = hash*uint64(a) + uint64(str[i])
		a = a * b
	}

	return hash
}

func DEKHash(str []byte) uint64{
	hash := uint64(len(str))
	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) ^ (hash >> 27)) ^ uint64(str[i])
	}

	return hash
}
