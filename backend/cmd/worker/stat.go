package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wbrost-go/internal/config"
	"wbrost-go/internal/repository/article"
	"wbrost-go/internal/repository/database/postgres"
	"wbrost-go/internal/repository/stat"
	"wbrost-go/internal/repository/user"
	"wbrost-go/internal/service/wb"
)

func main() {
	// Флаги командной строки
	var runOnce bool
	var interval int
	var slowMode bool // Новый флаг!

	flag.BoolVar(&runOnce, "once", false, "Запустить один раз и выйти")
	flag.IntVar(&interval, "interval", 60, "Интервал в секундах между запусками (по умолчанию 300 = 5 минут)")
	flag.BoolVar(&slowMode, "slow", false, "Медленный режим для больших периодов")
	flag.Parse()

	// Загружаем конфиг
	cfg := config.Load()
	if slowMode {
		fmt.Println("🐌 Включен медленный режим для работы с большими периодами")
		// Здесь можно добавить дополнительную логику для медленного режима
	}
	// Инициализируем БД
	db, err := postgres.NewPostgresDB(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	fmt.Println("✓ Подключение к БД установлено")

	// Инициализируем репозитории
	userRepo := user.NewUserRepository(db)
	statsGetRepo := stat.NewWBStatsGetRepository(db)
	statRepo := stat.NewStatRepository(db)
	articlesGetRepo := article.NewWBArticlesGetRepository(db)
	articleRepo := article.NewWBArticlesRepository(db)

	// Инициализируем сервис
	wbService := wb.NewWBService(userRepo, statsGetRepo, statRepo, articlesGetRepo, articleRepo)

	// Определяем интервал
	if interval == 0 {
		interval = cfg.Worker.Interval
	}

	if runOnce {
		start := time.Now() // Начало отсчета
		fmt.Printf("💤 Следующий запуск через %d секунд...\n", interval)
		// Запускаем один раз
		fmt.Println("🚀 Запуск обработки отчетов...")
		if err := wbService.ProcessPendingOrders(); err != nil {
			log.Printf("❌ Ошибка обработки: %v", err)
			os.Exit(1)
		}
		duration := time.Since(start) // Конец отсчета
		fmt.Printf("✅ Обработка завершена - заняло по времени: %v", duration)
		os.Exit(0)
	}

	// Запускаем как демон
	fmt.Printf("🔄 Запуск воркера с интервалом %d секунд...\n", interval)

	// Канал для сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Таймер для интервалов
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Первый запуск сразу
	fmt.Println("🎯 Первоначальная обработка...")
	if err := wbService.ProcessPendingOrders(); err != nil {
		log.Printf("⚠️ Ошибка при первоначальной обработке: %v", err)
	}

	// Основной цикл
	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n⏰ Запуск обработки в %s\n", time.Now().Format("2006-01-02 15:04:05"))
			start := time.Now() // Начало отсчета
			if err := wbService.ProcessPendingOrders(); err != nil {
				log.Printf("⚠️ Ошибка обработки: %v", err)
			}
			duration := time.Since(start) // Конец отсчета
			fmt.Printf("✅ Обработка завершена за %v\n", duration)
			fmt.Printf("💤 Следующий запуск через %d секунд...\n", interval)

		case sig := <-sigChan:
			fmt.Printf("\n🛑 Получен сигнал: %v. Завершение работы...\n", sig)
			return
		}
	}
}
