package storage

import "math/rand"

var EncodedUrls []EncodedUrl

type EncodedUrl struct {
	Id  string
	URL string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func Add(url string) string {
	for i := 0; i < len(EncodedUrls); i++ {
		if url == EncodedUrls[i].URL {
			return EncodedUrls[i].Id
		}
	}
	newURL := EncodedUrl{
		Id:  generateRandomId(url),
		URL: url,
	}
	_ = append(EncodedUrls, newURL)
	return newURL.Id
}

func Get(id string) string {
	for i := 0; i < len(EncodedUrls); i++ {
		if id == EncodedUrls[i].Id {
			return EncodedUrls[i].URL
		}
	}
	return "InvalidURL!"
}

func generateRandomId(url string) string {
	b := make([]rune, len(url))
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
