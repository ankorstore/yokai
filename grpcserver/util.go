package grpcserver

// Contains returns true if a given string can be found in a given slice of strings.
func Contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

// Unique remove duplicated elements from a slice of strings.
func Unique(list []string) []string {
	uniqueSlice := make([]string, 0, len(list))
	seen := make(map[string]bool, len(list))

	for _, element := range list {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}

	return uniqueSlice
}
