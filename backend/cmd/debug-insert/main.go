package main

import (
	"fmt"
	"log"
	"time"
	"wbrost-go/internal/config"
	"wbrost-go/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// 1. Показываем структуру таблицы
	fmt.Println("📋 Структура таблицы stat:")
	fmt.Println("==========================")

	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable, ordinal_position
		FROM information_schema.columns 
		WHERE table_name = 'stat' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatalf("Ошибка запроса структуры: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var colName, dataType, isNullable string
		var pos int
		rows.Scan(&colName, &dataType, &isNullable, &pos)
		fmt.Printf("%2d. %-25s %-15s NULL: %s\n", pos, colName, dataType, isNullable)
	}

	// 2. Проверяем точное количество столбцов
	var colCount int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'stat'").Scan(&colCount)
	if err != nil {
		log.Printf("Ошибка подсчета столбцов: %v", err)
	} else {
		fmt.Printf("\n📊 Всего столбцов: %d\n", colCount)
	}

	// 3. Тестируем INSERT с разным количеством полей
	fmt.Println("\n🔧 Тестируем INSERT запросы:")

	// Тест 1: Минимальные поля (должно работать)
	fmt.Println("\nТест 1: 3 поля (hash_info, user_id, created_at)")
	query1 := `INSERT INTO stat (hash_info, user_id, created_at) VALUES ($1, $2, $3) RETURNING id`

	var id1 int64
	err = db.QueryRow(query1, "test-hash-1", 1, time.Now()).Scan(&id1)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Успешно, ID: %d\n", id1)
		db.Exec("DELETE FROM stat WHERE id = $1", id1) // очистка
	}

	// Тест 2: 5 полей
	fmt.Println("\nТест 2: 5 полей")
	query2 := `INSERT INTO stat (hash_info, user_id, created_at, nm_id, ppvz_for_pay) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id2 int64
	err = db.QueryRow(query2, "test-hash-2", 1, time.Now(), 123456, "100.00").Scan(&id2)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Успешно, ID: %d\n", id2)
		db.Exec("DELETE FROM stat WHERE id = $1", id2)
	}

	// Тест 3: 10 полей
	fmt.Println("\nТест 3: 10 полей")
	query3 := `
		INSERT INTO stat (
			hash_info, user_id, created_at, nm_id, ppvz_for_pay,
			supplier_oper_name, delivery_rub, penalty, additional_payment, storage_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id
	`

	var id3 int64
	err = db.QueryRow(query3,
		"test-hash-3", 1, time.Now(), 123456, "100.00",
		1, 50.00, 0.00, 0.00, 0.00,
	).Scan(&id3)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Успешно, ID: %d\n", id3)
		db.Exec("DELETE FROM stat WHERE id = $1", id3)
	}

	// Тест 4: Все поля как в нашем запросе (50 полей)
	fmt.Println("\nТест 4: Считаем, сколько полей в нашем запросе")
	fmt.Println("В INSERT запросе у нас 50 полей. Давайте проверим соответствие:")

	// Список полей из нашего INSERT
	fields := []string{
		"hash_info", "user_id", "nm_id", "ppvz_for_pay", "supplier_oper_name",
		"delivery_rub", "penalty", "additional_payment", "storage_fee",
		"rebill_logistic_cost", "acquiring_fee", "acquiring_percent",
		"ppvz_sales_commission", "deduction", "ppvz_spp_prc", "ppvz_kvw_prc_base",
		"ppvz_kvw_prc", "acceptance", "dlv_prc", "created_at", "rr_dt",
		"shk_id", "sticker_id", "gi_id", "realizationreport_id", "barcode",
		"bonus_type_name", "last_error", "brand_name", "ppvz_office_id",
		"assembly_id", "sa_name", "ppvz_vw_nds", "ppvz_vw", "gi_box_type_name",
		"subject_name", "ts_name", "quantity", "retail_price", "retail_amount",
		"commission_percent", "office_name", "order_dt", "sale_dt",
		"delivery_amount", "return_amount", "report_type", "srid", "rid",
	}

	fmt.Printf("В нашем INSERT: %d полей\n", len(fields))
	fmt.Println("\nСравниваем с таблицей...")

	// Проверяем каждое поле
	missingFields := []string{}
	for _, field := range fields {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'stat' AND column_name = $1
			)
		`, field).Scan(&exists)

		if err != nil {
			fmt.Printf("Ошибка проверки поля %s: %v\n", field, err)
		} else if !exists {
			missingFields = append(missingFields, field)
			fmt.Printf("❌ Отсутствует в таблице: %s\n", field)
		} else {
			fmt.Printf("✅ Присутствует: %s\n", field)
		}
	}

	if len(missingFields) > 0 {
		fmt.Printf("\n⚠️ Отсутствуют поля: %v\n", missingFields)
	} else {
		fmt.Println("\n✅ Все поля присутствуют в таблице!")
	}
}
