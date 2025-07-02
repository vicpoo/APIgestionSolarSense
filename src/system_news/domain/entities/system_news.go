// system_news.go
package entities

import "time"

type SystemNews struct {
	ID        int32     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"column:title;type:varchar(255);not null"`
	Content   string    `json:"content" gorm:"column:content;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	AuthorID  int32     `json:"author_id" gorm:"column:author_id"`
}

// Setters
func (s *SystemNews) SetID(id int32) {
	s.ID = id
}

func (s *SystemNews) SetTitle(title string) {
	s.Title = title
}

func (s *SystemNews) SetContent(content string) {
	s.Content = content
}

func (s *SystemNews) SetCreatedAt(createdAt time.Time) {
	s.CreatedAt = createdAt
}

func (s *SystemNews) SetAuthorID(authorID int32) {
	s.AuthorID = authorID
}

// Getters
func (s *SystemNews) GetID() int32 {
	return s.ID
}

func (s *SystemNews) GetTitle() string {
	return s.Title
}

func (s *SystemNews) GetContent() string {
	return s.Content
}

func (s *SystemNews) GetCreatedAt() time.Time {
	return s.CreatedAt
}

func (s *SystemNews) GetAuthorID() int32 {
	return s.AuthorID
}

// Constructor
func NewSystemNews(title, content string, authorID int32) *SystemNews {
	return &SystemNews{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
}
