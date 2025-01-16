package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	metersInKm              = 1000.0 // Количество метров в одном километре
	minutesInHour           = 60     // Количество минут в одном часе
	defaultStepLength       = 0.65   // Длина одного шага в метрах
	cmInMeter               = 100.0  // Количество сантиметров в одном метре
	caloriesSpeedMultiplier = 18.0
	caloriesSpeedShift      = 1.79
)

const (
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2.0  // множитель веса пользователя
)

// Training содержит данные о тренировке.
type Training struct {
	Type       string        // Тип тренировки
	Action     int           // Количество повторов (шаги, гребки при плавании)
	StepLength float64       // Длина одного шага или гребка в метрах
	Duration   time.Duration // Продолжительность тренировки
	Weight     float64       // Вес пользователя в килограммах
}

// distance вычисляет дистанцию в километрах.
func (t Training) distance() float64 {
	if t.StepLength <= 0 {
		return 0
	}
	return float64(t.Action) * t.StepLength / metersInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч.
func (t Training) meanSpeed() float64 {
	if t.Duration <= 0 {
		return 0
	}
	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий (переопределяется).
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о тренировке.
type InfoMessage struct {
	Type      string        // Тип тренировки
	Duration  time.Duration // Длительность тренировки
	Distance  float64       // Дистанция в километрах
	MeanSpeed float64       // Средняя скорость в км/ч
	Calories  float64       // Количество потраченных калорий
}

// TrainingInfo собирает информацию о тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		Type:      t.Type,
		Duration:  t.Duration,
		Distance:  t.distance(),
		MeanSpeed: t.meanSpeed(),
		Calories:  t.Calories(),
	}
}

// String форматирует сообщение о тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.0f мин\nДистанция: %.2f км\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.Type,
		i.Duration.Minutes(),
		i.Distance,
		i.MeanSpeed,
		i.Calories,
	)
}

// Running описывает тренировку Бег.
type Running struct {
	Training
}

// Calories рассчитывает потраченные калории при беге.
func (r Running) Calories() float64 {
	if r.Weight <= 0 || r.meanSpeed() <= 0 {
		return 0
	}
	return (caloriesSpeedMultiplier*r.meanSpeed() + caloriesSpeedShift) * r.Weight / metersInKm * r.Duration.Hours() * minutesInHour
}

// Walking описывает тренировку Ходьба.
type Walking struct {
	Training
	Height float64 // Рост пользователя в сантиметрах
}

// Calories рассчитывает потраченные калории при ходьбе.
func (w Walking) Calories() float64 {
	if w.Weight <= 0 || w.Height <= 0 || w.meanSpeed() <= 0 {
		return 0
	}
	const (
		weightMultiplier      = 0.035
		speedHeightMultiplier = 0.029
	)
	speedMetersPerSecond := w.meanSpeed() * 1000 / 3600 // Перевод км/ч в м/с
	return ((weightMultiplier * w.Weight) + (math.Pow(speedMetersPerSecond, 2) / (w.Height / cmInMeter) * speedHeightMultiplier * w.Weight)) * w.Duration.Hours() * minutesInHour
}

// Swimming описывает тренировку Плавание.
type Swimming struct {
	Training
	PoolLength int // Длина бассейна в метрах
	PoolCount  int // Количество пересечений бассейна
}

// meanSpeed вычисляет среднюю скорость при плавании.
func (s Swimming) meanSpeed() float64 {
	if s.PoolLength <= 0 || s.Duration <= 0 {
		return 0
	}
	return float64(s.PoolLength*s.PoolCount) / metersInKm / s.Duration.Hours()
}

// Calories рассчитывает потраченные калории при плавании.
func (s Swimming) Calories() float64 {
	if s.Weight <= 0 || s.meanSpeed() <= 0 {
		return 0
	}
	const (
		caloriesSpeedShift   = SwimmingCaloriesMeanSpeedShift
		caloriesWeightFactor = SwimmingCaloriesWeightMultiplier
	)
	return (s.meanSpeed() + caloriesSpeedShift) * caloriesWeightFactor * s.Weight * s.Duration.Hours()
}

// CaloriesCalculator интерфейс для всех типов тренировок.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// ReadData возвращает информацию о тренировке.
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	return info.String()
}

func main() {
	swimming := Swimming{
		Training: Training{
			Type:       "Плавание",
			Action:     2000,
			StepLength: SwimmingLenStep,
			Duration:   90 * time.Minute,
			Weight:     85,
		},
		PoolLength: 50,
		PoolCount:  5,
	}
	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			Type:       "Ходьба",
			Action:     20000,
			StepLength: defaultStepLength,
			Duration:   3*time.Hour + 45*time.Minute,
			Weight:     85,
		},
		Height: 185,
	}
	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			Type:       "Бег",
			Action:     5000,
			StepLength: defaultStepLength,
			Duration:   30 * time.Minute,
			Weight:     85,
		},
	}
	fmt.Println(ReadData(running))
}
