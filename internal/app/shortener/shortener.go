package shortener

import (
	"encoding/json"
	"log"
	"os"
)

type encodedUrl struct {
	Id  string
	URL string
}

func EncodeString(val string) string {
	return checkData(val, false).Id
}

func DecodeString(val string) string {
	return checkData(val, true).URL
}

func checkData(val string, isFind bool) encodedUrl {
	var data []encodedUrl
	filePath := "./static/urls.json"

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(data); i++ {
		if val == data[i].URL && !isFind {
			return data[i]
		}
		if val == data[i].Id && isFind {
			return data[i]
		}
	}

	if !isFind {
		newURL := encodedUrl{Id: "1wq1324", URL: val}
		data = append(data, newURL)
		content, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(filePath, content, 0644)
		if err != nil {
			log.Fatal(err)
		}

		return newURL
	}
	return encodedUrl{Id: "", URL: ""}
}

// TODO: сделать реальный сокращатель ссылок? сделать базу данных? Сделать временное хранение данных? Почему этого нет в курсе?
