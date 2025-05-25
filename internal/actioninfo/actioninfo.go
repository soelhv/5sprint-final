package actioninfo

import (
	"fmt"
)

type DataParser interface {
	// Parse выполняет разбор строки данных и возвращает ошибку при неудаче
	Parse(string) error
	// ActionInfo возвращает информацию о выполненном действии после парсинга
	ActionInfo() (string, error)
}

// Info проходит по набору строк данных и использует DataParser для парсинга и вывода информации
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		// пытаемся распарсить текущую строку
		err := dp.Parse(data)
		if err != nil {
			fmt.Println("Error: parsing error:", err) // исправил вывод с fmt и ошибки на англ.
			continue
		}

		// получаем информацию об действиях после успешного парсинга
		info, err := dp.ActionInfo()
		if err != nil {
			fmt.Println("Error: information output error:", err) // исправил вывод с fmt и ошибки на англ.
			continue
		}

		println(info)
	}
}
