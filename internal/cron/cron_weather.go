package cron

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/olezhek28/wether-service/internal/clients"
)

// city - константа с названием города для которого собирается погода
const city = "moscow"

// WeatherService определяет контракт для сохранения погодных данных
// Используется для внедрения зависимости в cron-сервис
type WeatherService interface {
	AddWeather(ctx context.Context, name string, temperature float64, timestamp time.Time) error
}

// CronWeather представляет сервис для периодического сбора погодных данных
// Выполняет запланированные задачи по сбору температуры через внешние API
type CronWeather struct {
	client          *http.Client       // HTTP-клиент для выполнения запросов к API
	scheduler       gocron.Scheduler   // Планировщик задач для cron-выполнения
	geocodingClient *clients.Geocoding // Клиент для получения координат города
	openMeteo       *clients.OpenMeteo // Клиент для получения погодных данных
	weatherService  WeatherService     // Сервис для сохранения данных в хранилище
}

// New создает новый экземпляр CronWeather с инициализированными зависимостями
// Принимает планировщик задач и сервис погоды для внедрения зависимостей
func New(sheduler gocron.Scheduler, weatherService WeatherService) *CronWeather {
	// Создаем HTTP-клиент с таймаутом для предотвращения зависаний
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Инициализируем клиенты для внешних API
	geocodingClient := clients.NewGeocoding(client)
	openMeteo := clients.NewOpenMeteo(client)

	return &CronWeather{
		client:          client,
		scheduler:       sheduler,
		geocodingClient: geocodingClient,
		openMeteo:       openMeteo,
		weatherService:  weatherService,
	}
}

// Init инициализирует cron-задачи и возвращает список созданных jobs
// Создает периодическую задачу, которая выполняется каждые 10 секунд
func (c *CronWeather) Init(ctx context.Context) ([]gocron.Job, error) {
	// Создаем новую задачу в планировщике:
	// - DurationJob(10*time.Second) - задача выполняется каждые 10 секунд
	// - NewTask(c.cronTask, ctx) - выполняемая функция с контекстом
	job, err := c.scheduler.NewJob(gocron.DurationJob(
		10*time.Second,
	), gocron.NewTask(c.cronTask, ctx))
	if err != nil {
		slog.Error(err.Error())
		panic(err) // Паника в случае ошибки создания задачи (можно заменить на логирование)
	}

	return []gocron.Job{job}, nil
}

// cronTask - основная функция, выполняемая по расписанию
// Собирает данные о погоде и сохраняет их в хранилище
func (c *CronWeather) cronTask(ctx context.Context) {
	// 1. Получаем координаты города через геокодинг API
	geocodingRes, err := c.geocodingClient.GetCoordinate(city)
	if err != nil {
		slog.Error(err.Error())
		return // В случае ошибки просто выходим (можно добавить логирование)
	}

	// 2. Получаем температуру по координатам через OpenMeteo API
	openmeteoRes, err := c.openMeteo.GetTemperature(geocodingRes.Latitude, geocodingRes.Longitude)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// 3. Парсим временную метку из строкового формата
	// Формат "2006-01-02T15:04" - стандартный для Go (RFC 3339)
	timestamp, err := time.Parse("2006-01-02T15:04", openmeteoRes.Current.Time)
	if err != nil {
		slog.Error(err.Error()) // Логируем ошибку парсинга времени
		return
	}

	// 4. Сохраняем полученные данные в хранилище через сервис
	err = c.weatherService.AddWeather(ctx, city, openmeteoRes.Current.Temperature2m, timestamp)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
