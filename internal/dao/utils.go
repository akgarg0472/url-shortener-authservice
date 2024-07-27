package dao

func getStringOrNil(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func getInt64OrNil(i *int64) int64 {
	if i != nil {
		return *i
	}

	return -1
}
