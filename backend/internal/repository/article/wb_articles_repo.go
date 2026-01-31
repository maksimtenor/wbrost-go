package article

import (
	"fmt"
	"strings"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/database/postgres"
)

type WBArticlesRepository struct {
	db *postgres.PostgresDB
}

func NewWBArticlesRepository(db *postgres.PostgresDB) *WBArticlesRepository {
	return &WBArticlesRepository{db: db}
}

// CreateOrUpdate создает или обновляет карточку товара
func (r *WBArticlesRepository) CreateOrUpdate(article *entity.WBArticles) error {
	// Проверяем существование
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM wb_articles WHERE id_user = $1 AND articule = $2)",
		article.UserID,
		article.Articule,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check article existence: %w", err)
	}

	now := time.Now()

	if exists {
		// Обновляем существующую запись
		query := `
			UPDATE wb_articles 
			SET name = $1, photo = $2, updated = $3, updated_at = $4,
			    rus_size = $5, eu_size = $6, chrt_id = $7, 
			    barcode = $8, internal_id = $9
			WHERE id_user = $10 AND articule = $11
		`
		_, err := r.db.Exec(query,
			article.Name,
			article.Photo,
			now,
			now,
			article.RusSize,
			article.EuSize,
			article.ChrtID,
			article.Barcode,
			article.InternalID,
			article.UserID,
			article.Articule,
		)
		return err
	} else {
		// Создаем новую запись
		query := `
			INSERT INTO wb_articles (
				id_user, articule, name, photo, cost_price, created, 
				updated, updated_at, rus_size, eu_size, chrt_id, 
				barcode, internal_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING id
		`
		return r.db.QueryRow(query,
			article.UserID,
			article.Articule,
			article.Name,
			article.Photo,
			article.CostPrice,
			now,
			now,
			now,
			article.RusSize,
			article.EuSize,
			article.ChrtID,
			article.Barcode,
			article.InternalID,
		).Scan(&article.ID)
	}
}

// GetByUserID получает карточки товаров пользователя
func (r *WBArticlesRepository) GetByUserID(userID int, page, pageSize int) ([]entity.WBArticles, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, id_user, articule, name, photo, cost_price, 
		       created, updated, updated_at, rus_size, eu_size,
		       chrt_id, barcode, internal_id
		FROM wb_articles 
		WHERE id_user = $1 
		ORDER BY updated DESC NULLS LAST, created DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []entity.WBArticles
	for rows.Next() {
		var a entity.WBArticles
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Articule,
			&a.Name,
			&a.Photo,
			&a.CostPrice,
			&a.Created,
			&a.Updated,
			&a.UpdatedAt,
			&a.RusSize,
			&a.EuSize,
			&a.ChrtID,
			&a.Barcode,
			&a.InternalID,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

// GetCountByUserID получает общее количество карточек пользователя
func (r *WBArticlesRepository) GetCountByUserID(userID int) (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM wb_articles WHERE id_user = $1",
		userID,
	).Scan(&count)

	return count, err
}

// UpdateCostPrice обновляет себестоимость товара
func (r *WBArticlesRepository) UpdateCostPrice(userID int, articule, costPrice string) error {
	query := `
		UPDATE wb_articles 
		SET cost_price = $1, updated = $2, updated_at = $3
		WHERE id_user = $4 AND articule = $5
	`

	_, err := r.db.Exec(query, costPrice, time.Now(), time.Now(), userID, articule)
	return err
}

// SearchArticles поиск карточек по названию или артикулу
func (r *WBArticlesRepository) SearchArticles(userID int, search string, page, pageSize int) ([]entity.WBArticles, error) {
	offset := (page - 1) * pageSize
	searchPattern := "%" + strings.ToLower(search) + "%"

	query := `
		SELECT id, id_user, articule, name, photo, cost_price, 
		       created, updated, updated_at, rus_size, eu_size,
		       chrt_id, barcode, internal_id
		FROM wb_articles 
		WHERE id_user = $1 
		  AND (LOWER(name) LIKE $2 OR articule::text LIKE $2)
		ORDER BY updated DESC NULLS LAST, created DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, userID, searchPattern, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []entity.WBArticles
	for rows.Next() {
		var a entity.WBArticles
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Articule,
			&a.Name,
			&a.Photo,
			&a.CostPrice,
			&a.Created,
			&a.Updated,
			&a.UpdatedAt,
			&a.RusSize,
			&a.EuSize,
			&a.ChrtID,
			&a.Barcode,
			&a.InternalID,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}
