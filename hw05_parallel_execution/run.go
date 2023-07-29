package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorsCount struct {
	Count int32
	Limit int32
}

func (er *ErrorsCount) incr() {
	atomic.AddInt32(&er.Count, 1)
}

func (er *ErrorsCount) check() bool {
	return atomic.LoadInt32(&er.Limit) <= 0 || atomic.LoadInt32(&er.Count) < atomic.LoadInt32(&er.Limit)
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Создаём буферизованный канал с задачами
	ch := make(chan Task, len(tasks))

	// Создаём структуру для контроля за кол-вом ошибок и превышением лимита
	errorsCount := ErrorsCount{
		Count: 0,
		Limit: int32(m),
	}
	// Открываем waitGroup на кол-во горутин
	wg := &sync.WaitGroup{}
	wg.Add(n)

	// Запускаем горутины
	for i := 0; i < n; i++ {
		go worker(ch, &errorsCount, wg)
	}

	// Заполняем канал задачами и закрываем его
	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	wg.Wait() // Ожидаем завершения работы всех горутин

	// Возвращаем ошибку, если превысили лимит
	if !errorsCount.check() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(tasks <-chan Task, errorsCount *ErrorsCount, wg *sync.WaitGroup) {
	defer wg.Done() // Сообщаем, что завершили горутину

	for task := range tasks {
		// Если задача выполнилась с ошибкой - инкрементим счётчик ошибок
		if task() != nil {
			errorsCount.incr()
		}
		// Прерываем работу горутины, если превысили лимит ошибок
		if !errorsCount.check() {
			break
		}
	}
}
