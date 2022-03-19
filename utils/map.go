package utils

func RemoveFromStringMap(m map[string]string, s []string) {
	for _, v := range s {
		delete(m, v)
	}
}
