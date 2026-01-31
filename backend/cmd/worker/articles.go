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

	flag.BoolVar(&runOnce, "once", false, "Запустить один раз и выйти")
	flag.IntVar(&interval, "interval", 0, "Интервал в секундах между запусками (0 = использовать из конфига)")
	flag.Parse()

	// Загружаем конфиг
	cfg := config.Load()

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
	articlesService := wb.NewWBService(userRepo, statsGetRepo, statRepo, articlesGetRepo, articleRepo)

	// Определяем интервал
	if interval == 0 {
		interval = cfg.Worker.ArticlesInterval // Добавьте это поле в конфиг
	}

	if runOnce {
		// Запускаем один раз
		fmt.Println("🚀 Запуск обработки карточек товаров...")
		if err := articlesService.ProcessPendingArticles(); err != nil {
			log.Printf("❌ Ошибка обработки: %v", err)
			os.Exit(1)
		}
		fmt.Println("✅ Обработка карточек завершена")
		os.Exit(0)
	}

	// Запускаем как демон
	fmt.Printf("🔄 Запуск воркера карточек с интервалом %d секунд...\n", interval)

	// Канал для сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Таймер для интервалов
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Первый запуск сразу
	fmt.Println("🎯 Первоначальная обработка карточек...")
	if err := articlesService.ProcessPendingArticles(); err != nil {
		log.Printf("⚠️ Ошибка при первоначальной обработке: %v", err)
	}

	// Основной цикл
	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n⏰ Запуск обработки карточек в %s\n", time.Now().Format("2006-01-02 15:04:05"))
			if err := articlesService.ProcessPendingArticles(); err != nil {
				log.Printf("⚠️ Ошибка обработки карточек: %v", err)
			}
			fmt.Printf("💤 Следующий запуск через %d секунд...\n", interval)

		case sig := <-sigChan:
			fmt.Printf("\n🛑 Получен сигнал: %v. Завершение работы...\n", sig)
			return
		}
	}
}
