package clients

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// geocodingUrl - шаблон URL для Geocoding API Open-Meteo
// Параметры:
// - name=%s: название города для поиска
// - count=1: возвращать только первый результат
// - language=ru: язык возвращаемых данных (русский)
// - format=json: формат ответа (JSON)
const geocodingUrl = "https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json"

// GeocodingResponse представляет структуру ответа от Geocoding API
// Содержит информацию о найденном городе и его координатах
type GeocodingResponse struct {
	Name      string  `json:"name"`      // Название города
	Country   string  `json:"country"`   // Название страны
	Latitude  float64 `json:"latitude"`  // Географическая широта
	Longitude float64 `json:"longitude"` // Географическая долгота
}

// Geocoding - клиент для работы с Geocoding API
// Инкапсулирует логику взаимодействия с сервисом геокодинга
type Geocoding struct {
	httpClient *http.Client // HTTP-клиент для выполнения запросов
}

// NewGeocoding создает новый экземпляр клиента Geocoding
// Принимает готовый HTTP-клиент для переиспользования соединений
func NewGeocoding(httpClient *http.Client) *Geocoding {
	return &Geocoding{
		httpClient: httpClient,
	}
}

// GetCoordinate выполняет запрос к Geocoding API для получения координат города
// Возвращает информацию о городе и его координаты или ошибку в случае неудачи
func (g *Geocoding) GetCoordinate(city string) (GeocodingResponse, error) {
	// Формируем URL запроса с подстановкой названия города
	res, err := g.httpClient.Get(
		fmt.Sprintf(geocodingUrl, city),
	)
	if err != nil {
		slog.Error(err.Error())
		return GeocodingResponse{}, err // Возвращаем ошибку сети или таймаута
	}

	// Гарантируем закрытие тела ответа при выходе из функции
	// Важно для предотвращения утечек ресурсов
	defer res.Body.Close()

	// Проверяем HTTP-статус ответа
	// Ожидаем статус 200 OK, иначе считаем запрос неудачным
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code %d", res.StatusCode)
		slog.Error(err.Error())
		return GeocodingResponse{}, fmt.Errorf("status code %d", res.StatusCode)
	}

	// Структура для парсинга JSON ответа
	// API возвращает объект с массивом results, даже для одного элемента
	var geoResp struct {
		Results []GeocodingResponse `json:"results"` // Массив найденных городов
	}

	// Декодируем JSON из тела ответа в структуру
	// json.NewDecoder более эффективен для потокового чтения
	err = json.NewDecoder(res.Body).Decode(&geoResp)
	if err != nil {
		slog.Error(err.Error())
		return GeocodingResponse{}, err // Возвращаем ошибку парсинга JSON
	}

	// Возвращаем первый (и единственный) результат из массива
	// Так как в запросе указано count=1, массив содержит 0 или 1 элемент
	return geoResp.Results[0], nil
}
