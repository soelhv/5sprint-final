package spentenergy

import (
	"errors"
	"time"
)

const (
	mInKm                      = 1000.0
	minInH                     = 60.0
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func Distance(steps int, height float64) float64 {
	if steps < 0 || height <= 0 {
		return 0
	}
	return float64(steps) * (height * stepLengthCoefficient) / mInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	return Distance(steps, height) / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("invalid step value")
	}
	if weight <= 0 {
		return 0, errors.New("invalid weight value")
	}
	if height <= 0 {
		return 0, errors.New("invalid height value")
	}
	if duration <= 0 {
		return 0, errors.New("invalid duration value")
	}

	speed := MeanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("invalid step value")
	}
	if weight <= 0 {
		return 0, errors.New("invalid weight value")
	}
	if height <= 0 {
		return 0, errors.New("invalid height value")
	}
	if duration <= 0 {
		return 0, errors.New("invalid duration value")
	}

	calories, _ := RunningSpentCalories(steps, weight, height, duration)
	return calories * walkingCaloriesCoefficient, nil
}