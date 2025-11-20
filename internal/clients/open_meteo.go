package clients

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// openMeteoUrl - шаблон URL для Open-Meteo Weather API
// Параметры:
// - latitude=%f: географическая широта (подставляется как float)
// - longitude=%f: географическая долгота (подставляется как float)
// - current=temperature_2m: запрос текущей температуры на высоте 2 метра
const openMeteoUrl = "https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m"

// OpenMeteoResponse представляет структуру ответа от Open-Meteo API
// Содержит текущие погодные данные для запрошенных координат
type OpenMeteoResponse struct {
	Current struct {
		Time          string  `json:"time"`           // Временная метка измерения в формате ISO 8601
		Temperature2m float64 `json:"temperature_2m"` // Температура воздуха на высоте 2 метра в градусах Цельсия
	}
}

// OpenMeteo - клиент для работы с Open-Meteo Weather API
// Инкапсулирует логику получения текущих погодных данных по координатам
type OpenMeteo struct {
	httpClient *http.Client // HTTP-клиент для выполнения запросов
}

// NewOpenMeteo создает новый экземпляр клиента OpenMeteo
// Принимает готовый HTTP-клиент для переиспользования соединений
func NewOpenMeteo(httpClient *http.Client) *OpenMeteo {
	return &OpenMeteo{
		httpClient: httpClient,
	}
}

// GetTemperature выполняет запрос к Open-Meteo API для получения текущей температуры
// Принимает географические координаты (широту и долготу)
// Возвращает структуру с температурой и временем измерения или ошибку
func (c *OpenMeteo) GetTemperature(lat, long float64) (OpenMeteoResponse, error) {
	// Формируем URL запроса с подстановкой координат
	// fmt.Sprintf с %f форматирует float значения в строку
	res, err := c.httpClient.Get(
		fmt.Sprintf(openMeteoUrl, lat, long),
	)
	if err != nil {
		slog.Error(err.Error())
		return OpenMeteoResponse{}, err // Возвращаем ошибки сети, таймаута и т.д.
	}

	// Гарантируем закрытие тела ответа для предотвращения утечек ресурсов
	defer res.Body.Close()

	// Проверяем успешность HTTP-запроса
	// Статус 200 OK указывает на успешное выполнение запроса
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code %d", res.StatusCode)
		slog.Error(err.Error())
		return OpenMeteoResponse{}, fmt.Errorf("status code %d", res.StatusCode)
	}

	// Создаем структуру для парсинга JSON ответа
	var response OpenMeteoResponse

	// Декодируем JSON ответ непосредственно из потока тела ответа
	// Это более эффективно чем чтение всего тела в память и затем парсинг
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		slog.Error(err.Error())
		return OpenMeteoResponse{}, err // Возвращаем ошибки парсинга JSON
	}

	return response, nil
}
