package models

import (
	"encoding/json"
	"time"
)

// Weather представляет основную доменную модель погоды
// Используется в бизнес-логике приложения и для HTTP-ответов
type Weather struct {
	Name        string  `json:"name" db:"name"`               // Название города
	Temperature float64 `json:"temperature" db:"temperature"` // Температура в градусах
}

// ToResponse преобразует структуру Weather в JSON для HTTP-ответа
// Сериализует данные в формат, пригодный для отправки клиенту
// Возвращает байтовый массив с JSON данными или ошибку сериализации
func (w *Weather) ToResponse() ([]byte, error) {
	var raw []byte

	// Маршалинг структуры в JSON
	// json.Marshal преобразует Go-структуру в JSON байты
	raw, err := json.Marshal(w)
	if err != nil {
		return nil, err // Возвращаем ошибку если сериализация не удалась
	}

	return raw, nil
}

// WeatherDTO (Data Transfer Object) представляет модель данных для передачи между слоями
// Содержит дополнительные поля, необходимые для работы с хранилищем, но не для клиента
type WeatherDTO struct {
	Name        string    `json:"name" db:"name"`               // Название города
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`     // Временная метка измерения (из БД)
	Temperature float64   `json:"temperature" db:"temperature"` // Температура
}

// ToWeather преобразует WeatherDTO в доменную модель Weather
// Выполняет маппинг полей между DTO и доменной моделью
// Метод получает указатель на Weather для заполнения его полей
func (w *WeatherDTO) ToWeather(weather *Weather) {
	weather.Name = w.Name
	weather.Temperature = w.Temperature
	// Поле Timestamp не копируется, так как оно не нужно в доменной модели для клиента
}
