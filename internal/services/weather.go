package services

import (
	"context"
	"time"

	"github.com/olezhek28/wether-service/internal/domain/models"
)

// WeatherSaver определяет контракт для сохранения погодных данных
// Это интерфейс, который абстрагирует конкретную реализацию хранилища
type WeatherSaver interface {
	CreateWeatherCity(ctx context.Context, name string, temperature float64, timestamp time.Time) error
}

// WeatherProvider определяет контракт для получения погодных данных
// Интерфейс позволяет работать с разными источниками данных (БД, API, кэш и т.д.)
type WeatherProvider interface {
	ReadWeatherByCity(ctx context.Context, city string) (models.WeatherDTO, error)
}

// WeatherService представляет сервисный слой для работы с погодными данными
// Реализует бизнес-логику приложения, используя внедренные зависимости
type WeatherService struct {
	weatherSaver    WeatherSaver    // зависимость для сохранения данных
	weatherProvider WeatherProvider // зависимость для получения данных
}

// New создает новый экземпляр WeatherService с внедренными зависимостями
// Принимает реализации интерфейсов WeatherSaver и WeatherProvider
// Это пример Dependency Injection (DI) - принцип инверсии зависимостей
func New(weatherSaver WeatherSaver, weatherProvider WeatherProvider) *WeatherService {
	return &WeatherService{
		weatherSaver:    weatherSaver,
		weatherProvider: weatherProvider,
	}
}

// AddWeather добавляет новые погодные данные для города
// Делегирует операцию сохранения реализации WeatherSaver
// Является фасадом над методом хранилища, может содержать дополнительную бизнес-логику
func (w *WeatherService) AddWeather(ctx context.Context, name string, temperature float64, timestamp time.Time) error {
	return w.weatherSaver.CreateWeatherCity(ctx, name, temperature, timestamp)
}

// GetWeather получает погодные данные для указанного города
// Возвращает данные в формате доменной модели Weather
// Преобразует DTO (Data Transfer Object) в доменную модель
func (w *WeatherService) GetWeather(ctx context.Context, city string) (models.Weather, error) {
	var weather models.Weather // Доменная модель для возврата

	// Получаем данные через провайдер в формате DTO
	dto, err := w.weatherProvider.ReadWeatherByCity(ctx, city)
	if err != nil {
		return models.Weather{}, err // Возвращаем ошибку если данные не получены
	}

	// Преобразуем DTO в доменную модель
	// Метод ToWeather вероятно заполняет поля структуры Weather
	dto.ToWeather(&weather)

	return weather, nil // Возвращаем доменную модель
}
