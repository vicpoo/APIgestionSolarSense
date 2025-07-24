//src/system_news/infrastructure mysql_system_news_repository.go
package infrastructure

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/vicpoo/apigestion-solar-go/src/core"
	repositories "github.com/vicpoo/apigestion-solar-go/src/system_news/domain"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type MySQLSystemNewsRepository struct {
	conn *sql.DB
}

func NewMySQLSystemNewsRepository() repositories.ISystemNews {
	conn := core.GetBD()
	return &MySQLSystemNewsRepository{conn: conn}
}

func (repo *MySQLSystemNewsRepository) Save(news *entities.SystemNews) error {
	query := `
		INSERT INTO system_news (title, content, author_id, created_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := repo.conn.Exec(query, news.Title, news.Content, news.AuthorID, news.CreatedAt)
	if err != nil {
		log.Println("Error al guardar noticia del sistema:", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error al obtener ID generado:", err)
		return err
	}
	news.ID = int32(id)

	return nil
}

func (repo *MySQLSystemNewsRepository) Update(news *entities.SystemNews) error {
	query := `
		UPDATE system_news
		SET title = ?, content = ?, author_id = ?
		WHERE id = ?
	`
	result, err := repo.conn.Exec(query, news.Title, news.Content, news.AuthorID, news.ID)
	if err != nil {
		log.Println("Error al actualizar noticia del sistema:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error al obtener filas afectadas:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("noticia con ID %d no encontrada", news.ID)
	}

	return nil
}

func (repo *MySQLSystemNewsRepository) Delete(id int32) error {
	query := "DELETE FROM system_news WHERE id = ?"
	result, err := repo.conn.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar la noticia del sistema:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error al obtener filas afectadas:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("noticia con ID %d no encontrada", id)
	}

	return nil
}

func (repo *MySQLSystemNewsRepository) GetById(id int32) (*entities.SystemNews, error) {
	query := `
		SELECT id, title, content, created_at, author_id
		FROM system_news
		WHERE id = ?
	`
	row := repo.conn.QueryRow(query, id)

	var news entities.SystemNews
	err := row.Scan(
		&news.ID,
		&news.Title,
		&news.Content,
		&news.CreatedAt,
		&news.AuthorID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("noticia con ID %d no encontrada", id)
		}
		log.Println("Error al buscar la noticia por ID:", err)
		return nil, err
	}

	return &news, nil
}

func (repo *MySQLSystemNewsRepository) GetAll() ([]entities.SystemNews, error) {
	query := `
		SELECT id, title, content, created_at, author_id
		FROM system_news
	`
	rows, err := repo.conn.Query(query)
	if err != nil {
		log.Println("Error al obtener todas las noticias:", err)
		return nil, err
	}
	defer rows.Close()

	var newsList []entities.SystemNews
	for rows.Next() {
		var news entities.SystemNews
		err := rows.Scan(
			&news.ID,
			&news.Title,
			&news.Content,
			&news.CreatedAt,
			&news.AuthorID,
		)
		if err != nil {
			log.Println("Error al escanear noticia:", err)
			return nil, err
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error después de iterar filas:", err)
		return nil, err
	}

	return newsList, nil
}
