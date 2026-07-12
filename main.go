package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 1. Запускаем виртуальный Redis в памяти (без скачивания программ!)
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("Не удалось запустить miniredis: %v", err)
	}
	defer mr.Close()

	// 2. Создаем стандартного клиента Go-Redis и подключаем его к нашему виртуальному серверу
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(), // берем адрес нашего miniredis
	})

	fmt.Println("[GO] Успешно подключились к Redis!")

	// 3. Симулируем добавление задач в очередь (Аналог LPUSH в CLI)
	// Первый аргумент — контекст, второй — имя ключа (очереди), третий — значение
	err = rdb.LPush(ctx, "my_queue", "Отправить приветственное письмо").Err()
	if err != nil {
		log.Fatalf("Ошибка LPUSH: %v", err)
	}
	rdb.LPush(ctx, "my_queue", "Сгенерировать PDF отчет")

	fmt.Println("[GO] Закинули 2 задачи в очередь 'my_queue'...")

	// 4. Симулируем работу воркера — забираем задачу с другого конца (Аналог RPOP)
	// Метод Result() возвращает само значение строки и ошибку, если она есть
	task, err := rdb.RPop(ctx, "my_queue").Result()
	if err != nil {
		log.Fatalf("Ошибка RPOP: %v", err)
	}

	fmt.Printf("[WORKER] Забрал задачу из Redis: %s\n", task)
}
