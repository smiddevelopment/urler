package storage

import "math/rand"

var EncodedURLs []EncodedURL

type EncodedURL struct {
	Id  string
	URL string
}

var urlRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func Add(url string) string {
	newURL := EncodedURL{
		Id:  generateRandomID(),
		URL: url,
	}

	if EncodedURLs == nil {
		EncodedURLs = []EncodedURL{newURL}

		return newURL.Id
	}

	for i := 0; i < len(EncodedURLs); i++ {
		if url == EncodedURLs[i].URL {
			return EncodedURLs[i].Id
		}
	}

	_ = append(EncodedURLs, newURL)

	return newURL.Id
}

func Get(id string) string {
	for i := 0; i < len(EncodedURLs); i++ {
		if id == EncodedURLs[i].Id {
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
