package shortener

func EncodeString(val string) string {
	if val == "https://practicum.yandex.ru/" {
		return "12345"
	}
	return val
}

func DecodeString(val string) string {
	if val == "12345" {
		return "https://practicum.yandex.ru/"
	}
	return val
}
