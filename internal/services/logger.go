package services

import (
	"io"
	"log"
	"os"
	"sync"
)

type LoggerService struct {
	file  *os.File
	mutex sync.Mutex
}

// Конструктор для LoggerService
func NewLoggerService(logFilePath string) (*LoggerService, error) {
	// Создаем файл для логирования
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// Создаем MultiWriter для записи как в файл, так и в консоль
	multiWriter := io.MultiWriter(file, os.Stdout)

	log.SetOutput(multiWriter) // Устанавливаем MultiWriter как output для стандартного логгера

	return &LoggerService{
		file: file,
	}, nil
}

// Метод для записи логов
func (l *LoggerService) Log(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	log.Println(message)
}

// Метод для форматированного логирования
func (l *LoggerService) Logf(format string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	log.Printf(format, args...)
}

// Метод для закрытия файла
func (l *LoggerService) Close() error {
	return l.file.Close()
}
