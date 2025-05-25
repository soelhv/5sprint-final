package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Training хранит данные об одной тренировке
type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

// Parse разбирает строку "steps, type, duration"
// и заполняет поля структуры. Проверяет:
// - формат из трёх частей через запятую,
// - отсутствие пробелов в количестве шагов,
// - положительный int-количество шагов,
// - наличие и корректность единицы измерения в duration,
// - неотрицательность и ненулевую длительность.
// При ошибке возвращает её, при успехе — nil.

func (t *Training) Parse(datastring string) error {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return errors.New("invalid input format")
	}

	// Шаги
	stepStr := parts[0]
	// проверяем, что строка с количеством шагов не содержит пробелов
	if stepStr != strings.TrimSpace(stepStr) {
		// возвращаем ошибку, если в stepStr есть пробелы
		return errors.New("количество шагов не должно содержать пробелов")
	}
	stepStr = strings.TrimPrefix(stepStr, "+")
	steps, err := strconv.Atoi(stepStr)
	if err != nil {
		// сохраняем оригинальную ошибку преобразования
		return fmt.Errorf("failed to parse step count %q: %w", stepStr, err)
	}
	if steps <= 0 {
		// проверяем, что количество шагов положительное
		return errors.New("step count must be a positive integer")
	}
	t.Steps = steps

	// Тип тренировки
	t.TrainingType = strings.TrimSpace(parts[1])
	if t.TrainingType == "" {
		return errors.New("training type not specified")
	}

	// Длительность
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return errors.New("invalid duration format")
	}
	if duration <= 0 {
		return errors.New("duration must be greater than zero")
	}
	t.Duration = duration

	return nil
}

// ActionInfo формирует подробный отчёт о тренировке.
// Если внутри есть ошибка (например, Duration==0 или неизвестный тип),
// она возвращается в error. Иначе — текстовый отчёт и nil.
func (t Training) ActionInfo() (string, error) {
	if t.Personal.Weight <= 0 || t.Personal.Height <= 0 {
		return "", errors.New("invalid user weight or height")
	}

	distance := spentenergy.Distance(t.Steps, t.Personal.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Personal.Height, t.Duration)

	var calories float64
	var err error

	switch strings.ToLower(t.TrainingType) {
	case "бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Personal.Weight, t.Personal.Height, t.Duration)
		if err != nil {
			return "", fmt.Errorf("error calculating running calories: %w", err)
		}
	case "ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Personal.Weight, t.Personal.Height, t.Duration)
		if err != nil {
			return "", fmt.Errorf("error calculating walking calories: %w", err)
		}
	default:
		return "", errors.New("unknown training type")
	}

	report := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		t.TrainingType,
		t.Duration.Hours(),
		distance,
		speed,
		calories,
	)
	return report, nil
}
