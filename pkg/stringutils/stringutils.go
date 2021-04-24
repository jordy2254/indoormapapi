package stringutils

func IsNilPtrOrEmpty(strPtr *string) bool{
	return strPtr == nil || *strPtr == ""
}