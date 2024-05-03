package hook

// Contains returns true if a provided list contains a provided item.
func Contains(list []string, item string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
