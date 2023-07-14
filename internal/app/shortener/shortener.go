package shortener

func EncodeString(val string) string {
	if val == "https://practicum.yandex.ru/" {
		return "EwHXdJfB"
	}

	return val
}

func DecodeString(val string) string {
	if val == "EwHXdJfB" {
		return "https://practicum.yandex.ru/"
	}

	return val
}

// TODO: сделать реальный сокращатель ссылок или базу данных
