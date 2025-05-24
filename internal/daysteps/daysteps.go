package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps содержит данные о дневной прогулке пользователя.
type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

// При любых ошибках возвращает ненулевую ошибку:
//   - неверный формат строки
//   - лишние пробелы внутри полей
//   - некорректное количество шагов
//   - отсутствие единицы измерения времени
//   - отрицательная или нулевая длительность
func (ds *DaySteps) Parse(datastring string) error {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return errors.New("неверный формат строки")
	}

	stepStr := parts[0]
	durStr := parts[1]

	// пробелы вокруг полей недопустимы
	if stepStr != strings.TrimSpace(stepStr) || durStr != strings.TrimSpace(durStr) {
		return errors.New("неверный формат строки — лишние пробелы")
	}
	stepStr = strings.TrimSpace(stepStr)
	durStr = strings.TrimSpace(durStr)

	// шаги
	stepStr = strings.TrimPrefix(stepStr, "+")
	steps, err := strconv.Atoi(stepStr)
	if err != nil || steps <= 0 {
		return errors.New("некорректное количество шагов")
	}
	ds.Steps = steps

	// проверка единицы измерения в конце: h, m или s
	if !(strings.HasSuffix(durStr, "h") || strings.HasSuffix(durStr, "m") || strings.HasSuffix(durStr, "s")) {
		return errors.New("отсутствует единица измерения времени")
	}

	// парсинг длительности
	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return errors.New("некорректная длительность прогулки")
	}
	if duration <= 0 {
		return errors.New("длительность должна быть больше нуля")
	}
	ds.Duration = duration

	return nil
}

// ActionInfo формирует строку с информацией о прогулке:
// количество шагов, дистанция и сожжённые калории.
func (ds DaySteps) ActionInfo() (string, error) {
	if ds.Duration == 0 {
		return "", errors.New("длительность прогулки равна 0")
	}

	dist := spentenergy.Distance(ds.Steps, ds.Height)
	cals, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		ds.Steps, dist, cals,
	), nil
}
