package pkg

func BoolPointer(b bool) *bool {
	return &b
}

func StringPointer(str string) *string {
	return &str
}

func IntPointer(val int) *int {
	return &val
}
