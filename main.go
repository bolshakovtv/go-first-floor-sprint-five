package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре
)

// Training — общая структура для всех тренировок.
type Training struct {
	TrainingType string        // тип тренировки
	Action       int           // количество повторов (шаги, гребки при плавании)
	LenStep      float64       // длина одного шага или гребка в м
	Duration     time.Duration // продолжительность тренировки
	Weight       float64       // вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
// Формула расчёта: количество_повторов * длина_шага / метров_в_км.
func (train Training) distance() float64 {
	// Пример простой проверки.
	if train.Action <= 0 || train.LenStep <= 0 {
		return 0
	}
	return float64(train.Action) * train.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость (км/ч).
func (train Training) meanSpeed() float64 {
	// Пример проверки на продолжительность.
	if train.Duration.Hours() == 0 {
		return 0
	}
	return train.distance() / train.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке.
// Пока возвращаем 0, так как этот метод будет переопределяться для каждого типа тренировки.
func (train Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведённой тренировке.
type InfoMessage struct {
	TrainingType string        // тип тренировки
	Duration     time.Duration // длительность тренировки
	Distance     float64       // расстояние, которое преодолел пользователь
	MeanSpeed    float64       // средняя скорость (км/ч)
	Calories     float64       // количество потраченных килокалорий
}

// TrainingInfo возвращает структуру InfoMessage с информацией о тренировке.
func (train Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: train.TrainingType,
		Duration:     train.Duration,
		Distance:     train.distance(),
		MeanSpeed:    train.meanSpeed(),
		Calories:     train.Calories(),
	}
}

// String возвращает строку с информацией о тренировке.
func (infoMsg InfoMessage) String() string {
	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.0f мин\nДистанция: %.2f км\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		infoMsg.TrainingType,
		infoMsg.Duration.Minutes(),
		infoMsg.Distance,
		infoMsg.MeanSpeed,
		infoMsg.Calories,
	)
}

// CaloriesCalculator — интерфейс для расчёта калорий.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчёта калорий при беге.
const (
	CaloriesMeanSpeedMultiplier = 18   // множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // коэффициент изменения средней скорости
)

// Running описывает тренировку «Бег».
type Running struct {
	Training
}

// Calories возвращает количество потраченных килокалорий при беге.
// Формула:
// ((18 * ср_скорость_в_км_ч + 1.79) * вес / метров_в_км * время_в_часах * мин_в_часе).
func (run Running) Calories() float64 {
	if run.meanSpeed() == 0 || run.Weight <= 0 {
		return 0
	}
	return (CaloriesMeanSpeedMultiplier*run.meanSpeed() + CaloriesMeanSpeedShift) *
		run.Weight / MInKm * float64(run.Duration.Hours()*MinInHours)
}

// TrainingInfo возвращает структуру InfoMessage (переопределяет метод из Training).
func (run Running) TrainingInfo() InfoMessage {
	return run.Training.TrainingInfo()
}

// Константы для расчёта калорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // коэффициент для роста
	KmHInMsec                     = 0.278 // коэффициент для перевода км/ч в м/с
)

// Walking описывает тренировку «Ходьба».
type Walking struct {
	Training
	Height float64 // рост пользователя в см
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
// Формула:
// ((0.035 * вес + ((скорость_в_м_с^2) / (рост_в_м)) * 0.029 * вес ) * время_в_часах * мин_в_часе).
func (walk Walking) Calories() float64 {
	if walk.Weight <= 0 || walk.Height <= 0 || walk.meanSpeed() == 0 {
		return 0
	}
	// Переводим скорость в м/с: (скорость_км_ч * 1000) / 3600 = скорость_км_ч * 0.277...
	speedInMs := walk.meanSpeed() * KmHInMsec
	heightInM := walk.Height / CmInM

	return ((CaloriesWeightMultiplier*walk.Weight +
		((math.Pow(speedInMs, 2) / heightInM) * CaloriesSpeedHeightMultiplier * walk.Weight)) *
		walk.Duration.Hours() * MinInHours)
}

// TrainingInfo возвращает структуру InfoMessage (переопределяет метод из Training).
func (walk Walking) TrainingInfo() InfoMessage {
	return walk.Training.TrainingInfo()
}

// Константы для расчёта калорий при плавании.
const (
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэфф. изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // множитель веса
)

// Swimming описывает тренировку «Плавание».
type Swimming struct {
	Training
	LengthPool int // длина бассейна в метрах
	CountPool  int // количество пересечений бассейна
}

// meanSpeed возвращает среднюю скорость при плавании (км/ч).
// Формула: (длина_бассейна * кол-во_пересечений) / метров_в_км / время_в_часах.
func (swim Swimming) meanSpeed() float64 {
	if swim.LengthPool <= 0 || swim.CountPool <= 0 || swim.Duration.Hours() == 0 {
		return 0
	}
	return float64(swim.LengthPool) * float64(swim.CountPool) / MInKm / swim.Duration.Hours()
}

// Calories возвращает количество потраченных калорий при плавании.
// Формула:
// (ср_скорость_км_ч + 1.1) * 2 * вес * время_в_часах.
func (swim Swimming) Calories() float64 {
	if swim.Weight <= 0 || swim.meanSpeed() == 0 {
		return 0
	}
	return (swim.meanSpeed() + SwimmingCaloriesMeanSpeedShift) *
		SwimmingCaloriesWeightMultiplier * swim.Weight * swim.Duration.Hours()
}

// TrainingInfo возвращает информацию о тренировке (переопределяет метод из Training).
func (swim Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: swim.TrainingType,
		Duration:     swim.Duration,
		Distance:     swim.distance(),
		MeanSpeed:    swim.meanSpeed(),
		Calories:     swim.Calories(),
	}
}

// ReadData возвращает итоговую строку с информацией о тренировке.
func ReadData(training CaloriesCalculator) string {
	// Получаем калории.
	spentCalories := training.Calories()
	// Получаем подробную информацию о тренировке.
	infoMsg := training.TrainingInfo()
	infoMsg.Calories = spentCalories
	return fmt.Sprint(infoMsg)
}

func main() {
	// Тренировка «Плавание».
	swimTraining := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}
	fmt.Println(ReadData(swimTraining))

	// Тренировка «Ходьба».
	walkTraining := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}
	fmt.Println(ReadData(walkTraining))

	// Тренировка «Бег».
	runTraining := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}
	fmt.Println(ReadData(runTraining))
}
