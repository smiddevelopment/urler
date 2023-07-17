package storage

import "math/rand"

// EncodedURLs хранилище структур ссылок
var EncodedURLs []EncodedURL

// EncodedURL структура для хранения ссылки
type EncodedURL struct {
	ID  string
	URL string
}

// Массив для генерации ID ссылок
var urlRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// Add добавление новой ссылки или поиск существующей по URL
func Add(url string) string {
	newURL := EncodedURL{
		ID:  generateRandomID(),
		URL: url,
	}

	if EncodedURLs == nil {
		EncodedURLs = []EncodedURL{newURL}

		return newURL.ID
	}

	for i := 0; i < len(EncodedURLs); i++ {
		if url == EncodedURLs[i].URL {
			return EncodedURLs[i].ID
		}
	}

	_ = append(EncodedURLs, newURL)

	return newURL.ID
}

// Get поиск существующей в хранилище ссылки по ID
func Get(id string) string {
	for i := 0; i < len(EncodedURLs); i++ {
		if id == EncodedURLs[i].ID {
			return EncodedURLs[i].URL
		}
	}

	return "InvalidURL!"
}

// generateRandomID генерация случайного ID в пределах 8 символов для ссылки
func generateRandomID() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = urlRunes[rand.Intn(len(urlRunes))]
	}

	return string(b)
}
