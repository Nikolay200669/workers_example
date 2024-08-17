package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func createWorker(ctx context.Context, id int, done chan struct{}) {
	defer ctx.Done()
	fmt.Println("Worker", id, "started")

	// создаем канал ожидания ответа
	respChan := make(chan string)

	go func() {
		select {
		case <-ctx.Done():
			return
		case resp := <-respChan:
			fmt.Println("Received response from worker", id, ": ", resp)
			return
		}
	}()

	// имитируем выполнение задачи с ошибкой
	if err := doSomeWork(); err != nil {
		close(respChan) // отправляем ошибку через канал
	} else {
		respChan <- "Success!" // отправляем ответ
	}

	fmt.Println("Worker", id, "finished")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// создаем 5 го-рутин
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			createWorker(ctx, id, make(chan struct{}))
			wg.Done()
		}(i)
	}

	// ожидаем завершения всех го-рутин
	wg.Wait()

	fmt.Println("All workers finished")
}

func doSomeWork() error {
	// имитируем выполнение задачи с ошибкой
	time.Sleep(5 * time.Second) // имитируем выполнение задачи
	if rand.Intn(2) == 0 {      // имитируем ошибку
		return fmt.Errorf("worker %d failed", 42)
	}
	return nil
}
