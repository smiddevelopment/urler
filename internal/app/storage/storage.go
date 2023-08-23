package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"os"

	"github.com/smiddevelopment/urler.git/internal/app/config"
)

// EncodedURLs хранилище структур ссылок
var EncodedURLs []URLEncoded

// URLEncoded структура для хранения ссылки
type URLEncoded struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}

type URLArrayIncoming struct {
	ID  string `json:"correlation_id,omitempty"`
	URL string `json:"original_url,omitempty"`
}

type URLArrayOutgoing struct {
	ID    string `json:"correlation_id,omitempty"`
	Short string `json:"short_url,omitempty"`
}

// Массив для генерации ID ссылок
var urlRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// InitStore инициализация хранилищ
func InitStore() {
	if config.ServerConfig.DBURL != "" {
		err := createTableInDB()
		if err != nil {
			panic(err)
		}
	}

}

// Add проверка наличия и добавление новой ссылки
func Add(url string) string {
	if url == "" {
		return "Invalid URL!"
	}

	// Поиск существующей ссылки по URL в базе данных, файле или массиве оперативной памяти
	if config.ServerConfig.DBURL != "" {
		var res = getIDFromDB(url)
		if res != "" {
			return res
		}
	} else if config.ServerConfig.URLFile != "" {
		var URLsInFile []URLEncoded
		r, _ := newReadURL()
		r.scanner.Scan()

		if len(r.scanner.Bytes()) > 0 {
			err := json.Unmarshal(r.scanner.Bytes(), &URLsInFile)
			if err != nil {
				panic(err)

			}

			for i := 0; i < len(URLsInFile); i++ {
				if url == URLsInFile[i].URL {
					return URLsInFile[i].Result
				}
			}

		}
	} else {
		// Поиск существующей ссылки по URL в массиве оперативной памяти
		for i := 0; i < len(EncodedURLs); i++ {
			if url == EncodedURLs[i].URL {
				return EncodedURLs[i].Result
			}
		}
	}

	newURL := URLEncoded{
		Result: generateRandomID(),
		URL:    url,
	}
	// Добавление новой ссылки, если она не была найдена ранее
	if config.ServerConfig.DBURL != "" {
		err := addURLToDB(newURL)
		if err != nil {
			return "Can't add url to database!"
		}
	} else if config.ServerConfig.URLFile != "" {
		var URLsInFile []URLEncoded
		r, _ := newReadURL()
		r.scanner.Scan()

		if len(r.scanner.Bytes()) > 0 {
			err := json.Unmarshal(r.scanner.Bytes(), &URLsInFile)
			if err != nil {
				panic(err)

			}

		}
		URLsInFile = append(URLsInFile, newURL)

		w, err := newWriteURL()
		if err != nil {
			return err.Error()
		}
		byteURL, err := json.Marshal(URLsInFile)
		if err != nil {
			return err.Error()
		}
		w.writer.Write(byteURL)
		w.writer.Flush()
	} else {
		EncodedURLs = append(EncodedURLs, newURL)
	}

	return newURL.Result
}

// AddRange проверка наличия и добавление новой ссылки
func AddRange(urls []URLArrayIncoming) []URLArrayOutgoing {
	var res []URLArrayOutgoing
	if len(urls) > 0 {
		if config.ServerConfig.DBURL != "" {
			db, err := sql.Open("postgres", config.ServerConfig.DBURL)
			if err != nil {
				panic(err)
			}

			var URLsToAdd []URLEncoded
			for i := 0; i < len(urls); i++ {
				var temp URLArrayOutgoing
				temp.ID = urls[i].ID

				var result string
				row := db.QueryRowContext(context.Background(),
					`SELECT result FROM urls WHERE url = $1 LIMIT 1`, urls[i].URL)
				// готовим переменную для чтения результата
				err = row.Scan(&result) // разбираем результат
				if err != nil {
					newURL := URLEncoded{
						Result: generateRandomID(),
						URL:    urls[i].URL,
					}
					URLsToAdd = append(URLsToAdd, newURL)
					temp.Short = config.ServerConfig.ResURL + "/" + newURL.Result
				} else {
					temp.Short = config.ServerConfig.ResURL + "/" + result
				}
				res = append(res, temp)
			}

			if len(URLsToAdd) > 0 {
				tx, errT := db.Begin()
				if errT != nil {
					panic(errT)
				}

				for _, v := range URLsToAdd {
					// все изменения записываются в транзакцию
					_, errC := tx.ExecContext(context.Background(),
						`INSERT INTO urls (url, result) VALUES ($1,$2)`, v.URL, v.Result)
					if errC != nil {
						// если ошибка, то откатываем изменения
						errC = tx.Rollback()
						if errC != nil {
							panic(errC)
						}
					}
				}

				errCommit := tx.Commit()
				if errCommit != nil {
					panic(errCommit)
				}
			}
		} else if config.ServerConfig.URLFile != "" {
			var URLsInFile []URLEncoded
			r, _ := newReadURL()
			r.scanner.Scan()

			if len(r.scanner.Bytes()) > 0 {
				err := json.Unmarshal(r.scanner.Bytes(), &URLsInFile)
				if err != nil {
					panic(err)

				}
				for i := 0; i < len(urls); i++ {
					var temp URLArrayOutgoing
					temp.ID = urls[i].ID
					for j := 0; j < len(URLsInFile); j++ {
						if urls[i].URL == URLsInFile[j].URL {
							temp.Short = config.ServerConfig.ResURL + "/" + URLsInFile[j].Result
						}
					}
					if temp.Short == "" {
						newURL := URLEncoded{
							Result: generateRandomID(),
							URL:    urls[i].URL,
						}
						URLsInFile = append(URLsInFile, newURL)
						temp.Short = config.ServerConfig.ResURL + "/" + newURL.Result
					}
					res = append(res, temp)
				}

				if len(URLsInFile) > 0 {

					w, errW := newWriteURL()
					if errW != nil {
						panic(errW.Error())

					}

					byteURL, err := json.Marshal(URLsInFile)
					if err != nil {
						panic(err.Error())
					}
					w.writer.Write(byteURL)

					w.writer.Flush()

				}

			}
		} else {
			for i := 0; i < len(urls); i++ {
				var temp URLArrayOutgoing
				temp.ID = urls[i].ID
				for j := 0; j < len(EncodedURLs); j++ {
					if urls[i].URL == EncodedURLs[j].URL {
						temp.Short = config.ServerConfig.ResURL + "/" + EncodedURLs[j].Result
					}
				}
				if temp.Short == "" {
					newURL := URLEncoded{
						Result: generateRandomID(),
						URL:    urls[i].URL,
					}
					temp.Short = config.ServerConfig.ResURL + "/" + newURL.Result
					EncodedURLs = append(EncodedURLs, newURL)

				}

				res = append(res, temp)
			}
		}
	}

	return res
}

// Get поиск существующей ссылки по ID
func Get(id string) string {
	if config.ServerConfig.DBURL != "" {
		var res = getURLFromDB(id)
		if res != "" {
			return res
		}
	} else if config.ServerConfig.URLFile != "" {
		var URLsInFile []URLEncoded
		r, _ := newReadURL()
		r.scanner.Scan()

		if len(r.scanner.Bytes()) > 0 {
			err := json.Unmarshal(r.scanner.Bytes(), &URLsInFile)
			if err != nil {
				panic(err)

			}

			for i := 0; i < len(URLsInFile); i++ {
				if id == URLsInFile[i].Result {
					return URLsInFile[i].URL
				}
			}

		}
	} else {
		// Поиск существующей ссылки по URL в массиве оперативной памяти
		for i := 0; i < len(EncodedURLs); i++ {
			if id == EncodedURLs[i].Result {
				return EncodedURLs[i].URL
			}
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

// region LocalFile

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

// endregion LocalFile

// region DataBase

func createTableInDB() error {
	db, err := sql.Open("postgres", config.ServerConfig.DBURL)
	if err != nil {
		return err
	}

	defer db.Close()

	_, errC := db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS urls (url VARCHAR (255) UNIQUE NOT NULL, result VARCHAR (255) UNIQUE NOT NULL)`,
	)

	if errC != nil {
		return errC
	}

	return nil
}

func addURLToDB(urls URLEncoded) error {
	db, err := sql.Open("postgres", config.ServerConfig.DBURL)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO urls (url, result) VALUES ($1,$2)`, urls.URL, urls.Result)
	if err != nil {
		return err
	}
	return nil
}

func getIDFromDB(URLOrigin string) string {
	var result string

	if URLOrigin != "" {
		db, err := sql.Open("postgres", config.ServerConfig.DBURL)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		row := db.QueryRowContext(context.Background(),
			`SELECT result FROM urls WHERE url = $1 LIMIT 1`, URLOrigin)
		// готовим переменную для чтения результата
		err = row.Scan(&result) // разбираем результат
		if err != nil {
			return result
		}

	}

	return result
}

func getURLFromDB(URLId string) string {
	var result string

	if URLId != "" {
		db, err := sql.Open("postgres", config.ServerConfig.DBURL)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		row := db.QueryRowContext(context.Background(),
			`SELECT url FROM urls WHERE result = $1 LIMIT 1`, URLId)
		// готовим переменную для чтения результата
		err = row.Scan(&result) // разбираем результат
		if err != nil {
			return result
		}

	}

	return result
}

// endregion DataBase
