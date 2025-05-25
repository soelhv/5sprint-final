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

// Parse разбирает строку вида "<шаги>,<длительность>".
// При любых ошибках возвращает ненулевую ошибку:
//   - неверный формат строки
//   - лишние пробелы внутри полей
//   - некорректное количество шагов
//   - отсутствие единицы измерения времени
//   - отрицательная или нулевая длительность
func (ds *DaySteps) Parse(datastring string) error {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return errors.New("Error: wrong string format") // комментарии исправил на английский
	}

	stepStr := parts[0]
	durStr := parts[1]

	// шаги
	stepStr = strings.TrimPrefix(stepStr, "+")
	steps, err := strconv.Atoi(stepStr)
	if err != nil {
		// возвращаем ошибку преобразования
		return fmt.Errorf("Error: failed to parse step count %q: %w", stepStr, err)
	}
	if steps <= 0 {
		// проверяем, что количество шагов положительное
		return errors.New("Error: step count must be a positive integer")
	}
	ds.Steps = steps

	// парсинг длительности
	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return errors.New("Error: invalid walk duration")
	}
	if duration <= 0 {
		return errors.New("Error: duration must be greater than zero")
	}
	ds.Duration = duration

	return nil
}

// ActionInfo формирует строку с информацией о прогулке: количество шагов, дистанция и сожжённые калории.
func (ds DaySteps) ActionInfo() (string, error) {
	if ds.Duration == 0 {
		return "", errors.New("Error: walk duration is zero")
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
