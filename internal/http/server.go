package http

import (
	"context"
	"fmt"
	"net/http"
)

// Server — структура, описывающая HTTP-сервер.
// Включает контекст для graceful shutdown, порт, хост и экземпляр http.Handler.
type Server struct {
	context  context.Context // Контекст для управления жизненным циклом сервера
	host     string          // Хост (обычно "0.0.0.0" или "localhost")
	port     int             // Порт, на котором будет запущен сервер
	handlers http.Handler    // Интерфейс для работы хэндлеров
}

// NewServer — конструктор для создания нового сервера.
// Принимает контекст, порт, хост и экземпляр http.Handler, возвращает *Server.
func NewServer(
	context context.Context,
	port int,
	host string,
	handlers http.Handler,
) *Server {
	return &Server{
		context:  context,
		port:     port,
		handlers: handlers,
	}
}

// MustRun — запускает HTTP-сервер и завершает приложение с фатальной ошибкой, если запуск невозможен.
// Формирует строку адреса (хост:порт) и запускает сервер.
func (s *Server) MustRun() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := http.ListenAndServe(addr, s.handlers)
	if err != nil {
		panic(err)
	}
}

// TODO: Сделать Gracefull Shutdown
func (s *Server) Stop() {}
