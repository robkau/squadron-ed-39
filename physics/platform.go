package physics

func deletePlatforms(platforms *[]*platform) {
	for i, p := range *platforms {
		if p != nil && p.health <= 0 {
			(*platforms)[i] = (*platforms)[len(*platforms)-1]
			// dereference dead platform pointer
			(*platforms)[len(*platforms)-1] = nil
			// shrink slice by 1
			*platforms = (*platforms)[:len(*platforms)-1]
		}
	}
}
