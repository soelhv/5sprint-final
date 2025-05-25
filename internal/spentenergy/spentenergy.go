package spentenergy

import (
	"errors"
	"time"
)

const (
	// mInKm — количество метров в одном километре
	mInKm = 1000.0
	// minInH — количество минут в одном часе
	minInH = 60.0
	// stepLengthCoefficient — коэффициент длины шага от роста
	stepLengthCoefficient = 0.45
	// walkingCaloriesCoefficient — коэффициент для вычисления калорий при ходьбе
	walkingCaloriesCoefficient = 0.5
)

// Distance рассчитывает пройденное расстояние в километрах по количеству шагов и росту
// Возвращает 0 для некорректных входных данных
func Distance(steps int, height float64) float64 {
	if steps < 0 || height <= 0 {
		return 0
	}
	// distance (км) = шаги * (рост * коэффициент) / 1000
	return float64(steps) * (height * stepLengthCoefficient) / mInKm
}

// MeanSpeed вычисляет среднюю скорость в км/ч по шагам, росту и времени
// Возвращает 0 для некорректных входных данных
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	// speed = distance (км) / время (ч)
	return Distance(steps, height) / duration.Hours()
}

// RunningSpentCalories оценивает калории, потраченные на бег
// steps — количество шагов, weight — вес в кг, height — рост в м, duration — длительность
// Формула: вес * скорость * время (мин) / 60
// Возвращает ошибку при неверных входных данных
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Error: invalid step value")
	}
	if weight <= 0 {
		return 0, errors.New("Error: invalid weight value")
	}
	if height <= 0 {
		return 0, errors.New("Error: invalid height value")
	}
	if duration <= 0 {
		return 0, errors.New("Error: invalid duration value")
	}

	// вычисляем среднюю скорость
	speed := MeanSpeed(steps, height, duration)
	// вычисляем калории: weight * speed * duration.Minutes() / minInH
	calories := (weight * speed * duration.Minutes()) / minInH
	return calories, nil
}

// WalkingSpentCalories оценивает калории при ходьбе
// Использует результат RunningSpentCalories и применяет walkingCaloriesCoefficient
// Возвращает ошибку при неверных входных данных
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Error: invalid step value")
	}
	if weight <= 0 {
		return 0, errors.New("Error: invalid weight value")
	}
	if height <= 0 {
		return 0, errors.New("Error: invalid height value")
	}
	if duration <= 0 {
		return 0, errors.New("Error: invalid duration value")
	}

	// получаем калории для бега и умножаем на коэффициент для ходьбы
	calories, _ := RunningSpentCalories(steps, weight, height, duration)
	return calories * walkingCaloriesCoefficient, nil
}