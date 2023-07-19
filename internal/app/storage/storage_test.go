package storage

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name  string
		value string
		want  string
	}{
		{
			name:  "Test empty URL", // описываем каждый тест:
			value: "",               // значения, которые будет принимать функция,
			want:  "Invalid URL!",   // и ожидаемый результат
		},
	}

	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			if val := Add(test.value); val != test.want {
				t.Errorf("Add() = " + val + ", want " + test.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name  string
		value string
		want  string
	}{
		{
			name:  "Test empty ID", // описываем каждый тест:
			value: "",              // значения, которые будет принимать функция,
			want:  "Invalid ID!",   // и ожидаемый результат
		},
	}

	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			if val := Get(test.value); val != test.want {
				t.Errorf("Get() = " + val + ", want " + test.want)
			}
		})
	}
}
