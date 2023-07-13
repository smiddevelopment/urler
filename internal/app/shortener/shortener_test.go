package shortener

import "testing"

func TestEncodeString(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name  string
		value string
		want  string
	}{
		{
			name:  "Test encode #`1",              // описываем каждый тест:
			value: "https://practicum.yandex.ru/", // значения, которые будет принимать функция,
			want:  "12345",                        // и ожидаемый результат
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			if val := EncodeString(test.value); val != test.want {
				t.Errorf("EncodeString() = " + val + ", want " + test.want)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name  string
		value string
		want  string
	}{
		{
			name:  "Test decode #`1",              // описываем каждый тест:
			value: "12345",                        // значения, которые будет принимать функция,
			want:  "https://practicum.yandex.ru/", // и ожидаемый результат
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			if val := DecodeString(test.value); val != test.want {
				t.Errorf("DecodeString() = " + val + ", want " + test.want)
			}
		})
	}
}
