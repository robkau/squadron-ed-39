package physics

func deletePlatforms(platforms *[]*platform) {
	// todo: debug
	j := 0
	for _, p := range *platforms {
		if p.health > 0 {
			(*platforms)[j] = p
			j++
		}
	}
	*platforms = (*platforms)[:j]

	/*
		todo: compare performance
		// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
		b := (*platforms)[:0]
		for _, p := range *platforms {
			if p.health >0 {
				b = append(b,
			}
		}
	*/
}
