package ptrconv

func String(s string) *string {
	return &s
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Ptr[T comparable](v T) *T {
	return &v
}

func SafeValue[T comparable](v *T) T {
	if v == nil {
		temp := new(T)
		return *temp
	}
	return *v
}
