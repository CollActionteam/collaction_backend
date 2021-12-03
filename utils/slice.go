package utils

func IndexOf(s []string, v string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == v {
			return i
		}
	}
	return -1
}

func Remove(s *[]string, i int) {
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}
