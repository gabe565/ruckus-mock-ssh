package cli

func longestCommonPrefix(strs []string) string {
	longestPrefix := ""
	endPrefix := false

	if len(strs) > 0 {
		first := strs[0]
		last := strs[len(strs)-1]

		for i := 0; i < len(first); i++ {
			if !endPrefix && string(last[i]) == string(first[i]) {
				longestPrefix += string(last[i])
			} else {
				endPrefix = true
			}
		}
	}
	return longestPrefix
}
