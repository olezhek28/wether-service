package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/olezhek28/wether-service/internal/domain/models"
)

// Weather представляет слой доступа к данным для работы с погодными данными
// Содержит подключение к базе данных для выполнения операций
type Weather struct {
	db *pgx.Conn // Подключение к PostgreSQL через драйвер pgx
}

// New создает и возвращает новый экземпляр Weather с переданным подключением к БД
// Используется для инициализации хранилища в основном приложении
func New(db *pgx.Conn) *Weather {
	return &Weather{
		db: db,
	}
}

// CreateWeatherCity создает новую запись о погоде для указанного города
// Принимает контекст для управления таймаутами и отменой, название города,
// температуру и временную метку измерения
// Возвращает ошибку в случае неудачи операции
func (w *Weather) CreateWeatherCity(ctx context.Context, name string, temperature float64, timestamp time.Time) error {
	// SQL-запрос для вставки данных в таблицу reading
	// Используются позиционные параметры $1, $2, $3 для защиты от SQL-инъекций
	query := "insert into reading (name, temperature, timestamp) values ($1, $2, $3)"

	// Выполнение SQL-запроса с передачей параметров
	rows, err := w.db.Exec(ctx, query, name, temperature, timestamp)
	if err != nil {
		return err // Возвращаем ошибку если запрос не выполнился
	}

	// Проверяем, что хотя бы одна строка была затронута операцией
	// Это гарантирует, что данные действительно были добавлены
	if rows.RowsAffected() == 0 {
		return errors.New("Weather was not add")
	}
	return nil
}

// ReadWeatherByCity возвращает последние погодные данные для указанного города
// Выполняет поиск самой свежей записи по временной метке
// Возвращает структуру WeatherDTO с данными или ошибку если город не найден
func (w *Weather) ReadWeatherByCity(ctx context.Context, city string) (models.WeatherDTO, error) {
	var weatherDto models.WeatherDTO // Структура для хранения результата

	// SQL-запрос для выборки последней записи погоды по городу
	// ORDER BY timestamp DESC - сортировка по убыванию времени
	// LIMIT 1 - берем только самую свежую запись
	query := "select name, timestamp,temperature from reading where name = $1 order by timestamp desc limit 1"

	// Выполнение запроса и сканирование результата в структуру
	err := w.db.QueryRow(ctx, query, city).Scan(&weatherDto.Name, &weatherDto.Timestamp, &weatherDto.Temperature)
	if err != nil {
		// Обработка случая когда город не найден в базе данных
		if errors.Is(err, pgx.ErrNoRows) {
			return models.WeatherDTO{}, errors.New("No city with same name")
		}
		// Возвращаем другие ошибки (проблемы с подключением, синтаксисом и т.д.)
		return models.WeatherDTO{}, err
	}

	return weatherDto, nil // Возвращаем успешно найденные данные
}
