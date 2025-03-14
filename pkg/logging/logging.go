package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

// Определяем структуру, которая реализует интерфейс logrus.Hook
type writeHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire метод выполняется, когда создается новая запись в журнале.
// Он записывает эту запись в каждый Writer, связанный с хуком.
func (hook *writeHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

// Levels возвращает все уровни, на которых хук будет активирован.
func (hook *writeHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Глобальная переменная для экземпляра Logger
var e *logrus.Entry

// Logger - это обертка над logrus.Entry
type Logger struct {
	*logrus.Entry
}

// GetLogger возвращает текущий экземпляр логгера
func GetLogger() *Logger {
	return &Logger{e}
}

// GetLoggerWithField возвращает новый логгер с дополнительным полем
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

// init инициализирует логгер при запуске программы
func init() {

	// Создаем новый экземпляр logrus.Logger
	l := logrus.New()

	// Включаем вывод информации о месте вызова (файл и строка)
	l.SetReportCaller(true)

	// Устанавливаем формат вывода записей лога
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// Выбираем короткое имя файла и номер строки для вывода
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		ForceColors: true, // Принудительное использование цвета даже если вывод не терминал
		//DisableColors: false, // Отключаем цветной вывод
		FullTimestamp: true, // Включаем вывод полного времени
	}

	// Создаем директорию для логов
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	// Открываем файл для записи логов
	allFile, err := os.OpenFile("logs/all.logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	// Перенаправляем стандартный вывод в "/dev/null"
	l.SetOutput(io.Discard)

	// Добавляем хук, который записывает логи в файл и в стандартный вывод
	l.AddHook(&writeHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	// Устанавливаем уровень логирования
	l.SetLevel(logrus.TraceLevel)

	// Создаем новый экземпляр Entry для логгера
	e = logrus.NewEntry(l)
}
