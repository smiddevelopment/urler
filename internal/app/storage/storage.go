package storage

import "math/rand"

var EncodedUrls []EncodedUrl

type EncodedUrl struct {
	Id  string
	URL string
}

var urlRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func Add(url string) string {
	newURL := EncodedUrl{
		Id:  generateRandomID(),
		URL: url,
	}

	if EncodedUrls == nil {
		EncodedUrls = []EncodedUrl{newURL}

		return newURL.Id
	}

	for i := 0; i < len(EncodedUrls); i++ {
		if url == EncodedUrls[i].URL {
			return EncodedUrls[i].Id
		}
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

func generateRandomID() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = urlRunes[rand.Intn(len(urlRunes))]
	}
	return string(b)
}
