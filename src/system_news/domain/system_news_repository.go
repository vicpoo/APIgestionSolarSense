//api/src/system_news/domain system_news_repository.go
package domain

import (
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type ISystemNews interface {
	Save(news *entities.SystemNews) error
	Update(news *entities.SystemNews) error
	Delete(id int32) error
	GetById(id int32) (*entities.SystemNews, error)
	GetAll() ([]entities.SystemNews, error)
}
