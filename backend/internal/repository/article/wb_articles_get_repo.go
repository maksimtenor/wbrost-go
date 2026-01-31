package article

import (
	"fmt"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/database/postgres"
)

type WBArticlesGetRepository struct {
	db *postgres.PostgresDB
}

func NewWBArticlesGetRepository(db *postgres.PostgresDB) *WBArticlesGetRepository {
	return &WBArticlesGetRepository{db: db}
}

// Create создает новую запись запроса карточек товаров
func (r *WBArticlesGetRepository) Create(article *entity.WBArticlesGet) error {
	query := `
		INSERT INTO wb_articles_get (id_user, status, last_error)
		VALUES ($1, $2, $3)
		RETURNING id, created, updated
	`

	return r.db.QueryRow(query,
		article.UserID,
		article.Status,
		article.LastError,
	).Scan(&article.ID, &article.Created, &article.Updated)
}

// GetByUserID получает запросы пользователя
func (r *WBArticlesGetRepository) GetByUserID(userID int) ([]entity.WBArticlesGet, error) {
	query := `
		SELECT id, id_user, status, created, updated, last_error
		FROM wb_articles_get 
		WHERE id_user = $1 
		ORDER BY created DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []entity.WBArticlesGet
	for rows.Next() {
		var a entity.WBArticlesGet
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Status,
			&a.Created,
			&a.Updated,
			&a.LastError,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

// GetPendingArticles возвращает запросы с статусом 0 (в обработке)
func (r *WBArticlesGetRepository) GetPendingArticles() ([]entity.WBArticlesGet, error) {
	query := `
		SELECT id, id_user, status, created, updated, last_error
		FROM wb_articles_get 
		WHERE status = $1 
		ORDER BY created ASC
	`

	rows, err := r.db.Query(query, entity.ArticlesStatusWait)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending articles: %w", err)
	}
	defer rows.Close()

	var articles []entity.WBArticlesGet
	for rows.Next() {
		var article entity.WBArticlesGet
		err := rows.Scan(
			&article.ID,
			&article.UserID,
			&article.Status,
			&article.Created,
			&article.Updated,
			&article.LastError,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// UpdateStatus обновляет статус запроса
func (r *WBArticlesGetRepository) UpdateStatus(articleID int, status int, errorMsg string) error {
	query := `
		UPDATE wb_articles_get 
		SET status = $1, last_error = $2, updated = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, status, errorMsg, time.Now(), articleID)
	if err != nil {
		return fmt.Errorf("failed to update article status: %w", err)
	}

	return nil
}
