package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olezhek28/wether-service/internal/domain/models"
)

// WeatherService определяет контракт для сервиса погоды
// Интерфейс описывает методы, которые используются обработчиками HTTP
type WeatherService interface {
	AddWeather(ctx context.Context, name string, temperature float64, timestamp time.Time) error
	GetWeather(ctx context.Context, city string) (models.Weather, error)
}

// Handlers представляет слой обработчиков HTTP-запросов
// Содержит зависимости и маршрутизатор для обработки запросов
type Handlers struct {
	weatherService WeatherService // Сервис для работы с бизнес-логикой погоды
	r              *chi.Mux       // Маршрутизатор Chi для управления HTTP-маршрутами
}

// New создает новый экземпляр обработчиков с внедренными зависимостями
// Принимает маршрутизатор и сервис погоды для инициализации
func New(
	r *chi.Mux,
	weatherService WeatherService,
) *Handlers {
	return &Handlers{
		r:              r,
		weatherService: weatherService,
	}
}

// Init инициализирует маршруты и middleware для обработчиков
// Настраивает обработку HTTP-запросов и добавляет промежуточное ПО
func (h *Handlers) Init() {
	// Добавляем middleware для логирования всех запросов
	h.r.Use(middleware.Logger)

	// Регистрируем обработчик для GET запросов по пути /{city}
	// {city} - параметр маршрута, который будет извлекаться из URL
	h.r.Get("/{city}", h.getCity)
}

// getCity обрабатывает GET запрос для получения погоды по городу
// Это метод-обработчик, соответствующий интерфейсу http.HandlerFunc
func (h *Handlers) getCity(w http.ResponseWriter, r *http.Request) {
	// Получаем контекст из запроса для управления таймаутами и отменой
	ctx := r.Context()

	// Извлекаем параметр city из URL пути
	// Например, для /moscow вернет "moscow"
	city := chi.URLParam(r, "city")

	// Вызываем сервис для получения погодных данных
	// Делегируем бизнес-логику сервисному слою
	weather, err := h.weatherService.GetWeather(ctx, city)
	if err != nil {
		// В случае ошибки возвращаем статус 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching weather"))
		return // Важно: прекращаем выполнение после ошибки
	}

	// Преобразуем доменную модель в формат для HTTP-ответа
	// Метод ToResponse вероятно сериализует данные в JSON или другой формат
	raw, err := weather.ToResponse()
	if err != nil {
		// Если преобразование не удалось, возвращаем ошибку сервера
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Записываем сырые данные в тело ответа
	// По умолчанию статус 200 OK
	w.Write(raw)
}
