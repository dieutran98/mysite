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
