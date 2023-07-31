package storage

import (
	"bufio"
	"encoding/json"
	"math/rand"
	"os"

	"github.com/smiddevelopment/urler.git/internal/app/config"
)

// EncodedURLs хранилище структур ссылок
var EncodedURLs []URLEncoded

//// Database хранилище структур ссылок
//type Database struct {
//	EncodedURLs []URLEncoded
//}
//
//func (box *Database) AddItem(item URLEncoded) []URLEncoded {
//	box.EncodedURLs = append(box.EncodedURLs, item)
//	return box.EncodedURLs
//}

// URLEncoded структура для хранения ссылки
type URLEncoded struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}

// Массив для генерации ID ссылок
var urlRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// InitDB инициализация и чтение хранилища
func InitDB() {
	r, _ := newReadURL()
	r.scanner.Scan()
	if len(r.scanner.Bytes()) > 0 {
		err := json.Unmarshal(r.scanner.Bytes(), &EncodedURLs)
		if err != nil {
			panic(err)

		}
	}
}

// Add добавление новой ссылки или поиск существующей по URL
func Add(url string) string {
	if url == "" {
		return "Invalid URL!"
	}

	// Поиск существующей ссылки по URL в массиве
	for i := 0; i < len(EncodedURLs); i++ {
		if url == EncodedURLs[i].URL {
			return EncodedURLs[i].Result
		}
	}

	newURL := URLEncoded{
		Result: generateRandomID(),
		URL:    url,
	}
	// Добавление новой ссылки, если она не была найдена ранее
	EncodedURLs = append(EncodedURLs, newURL)
	w, err := newWriteURL()
	if err != nil {
		return err.Error()
	}
	byteURL, err := json.Marshal(EncodedURLs)
	if err != nil {
		return err.Error()
	}
	w.writer.Write(byteURL)
	w.writer.Flush()
	return newURL.Result
}

// Get поиск существующей в хранилище ссылки по ID
func Get(id string) string {
	// Поиск существующей ссылки по ID в массиве
	for i := 0; i < len(EncodedURLs); i++ {
		if id == EncodedURLs[i].Result {
			return EncodedURLs[i].URL
		}
	}

	return "Invalid ID!"
}

// generateRandomID генерация случайного ID в пределах 8 символов для ссылки
func generateRandomID() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = urlRunes[rand.Intn(len(urlRunes))]
	}

	return string(b)
}

type readURL struct {
	file *os.File
	// добавляем reader в readURL
	scanner *bufio.Scanner
}

func newReadURL() (*readURL, error) {
	file, err := os.OpenFile(config.ServerConfig.URLFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &readURL{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *readURL) ReadEvent() (*URLEncoded, error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	event := URLEncoded{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *readURL) Close() error {
	return c.file.Close()
}

type writeURL struct {
	file *os.File
	// добавляем writer во writeURL
	writer *bufio.Writer
}

func newWriteURL() (*writeURL, error) {
	file, err := os.OpenFile(config.ServerConfig.URLFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &writeURL{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *writeURL) WriteEvent(event *URLEncoded) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *writeURL) Close() error {
	// закрываем файл
	return p.file.Close()
}
