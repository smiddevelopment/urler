package storage

import "math/rand"

var EncodedURLs []EncodedURL

type EncodedURL struct {
	ID  string
	URL string
}

var urlRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

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

func Get(id string) string {
	for i := 0; i < len(EncodedURLs); i++ {
		if id == EncodedURLs[i].ID {
			return EncodedURLs[i].URL
		}
	}

	return "InvalidURL!"
}

func generateRandomID() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = urlRunes[rand.Intn(len(urlRunes))]
	}

	return string(b)
}
