package utils

// StringInSlice checks whether a slice contains a value
func StringInSlice(needle string, haystack []string) bool {

	valueMap := map[string]bool{}

	for _, value := range haystack {
		valueMap[value] = true
	}

	if _, ok := valueMap[needle]; ok {
		return true
	}

	return false

}

// StringSliceIntersection gets the intersection of multiple string slices
func StringSliceIntersection(slices [][]string) []string {

	valueMap := map[string]int{}
	result := []string{}

	for _, singleSlice := range slices {

		for _, sliceValue := range singleSlice {

			valueMap[sliceValue] = valueMap[sliceValue] + 1

		}

	}

	for value, count := range valueMap {

		if count == len(slices) {
			result = append(result, value)
		}

	}

	return result

}
