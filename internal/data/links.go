package data

import (
	"database/sql"
	"math/rand"

	"github.com/abdelkd/linkly/internal/util"
)

type CreateLinkRequest struct {
	Location string `json:"location"`
}

type CreateLinkResponse struct {
	Link string `json:"location"`
	Code string `json:"code"`
}

type Link struct {
	Code string
	Url  string
}

type LinksModel struct {
	DB *sql.DB
}

func (l *LinksModel) Add(link *Link) error {
	query := `
	  INSERT INTO links (url, code, hash)
	  VALUES ($1, $2, $3)
		RETURNING code;
	`

	alphabets := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	alphabetsLen := len(alphabets)
	buf := make([]byte, 7)

	for i := 0; i < 7; i++ {
		randInt := rand.Intn(alphabetsLen)
		buf[i] = alphabets[randInt]
	}

	args := []interface{}{link.Url, string(buf), util.HashString(link.Url)}

	return l.DB.QueryRow(query, args...).Scan(&link.Code)
}

func (l *LinksModel) Get(code string) (string, error) {
	query := `
	  SELECT (url)
	  FROM links
	  WHERE code = $1;
	`

	var url string
	err := l.DB.QueryRow(query, code).Scan(&url)

	return url, err
}

func (l *LinksModel) GetByHashCode(hashCode string) (*Link, error) {
	query := `
	  SELECT url, code
	  FROM links
	  WHERE hash = $1;
	`

	var link Link
	err := l.DB.QueryRow(query, hashCode).Scan(&link.Url, &link.Code)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

type Models struct {
	Links LinksModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Links: LinksModel{DB: db},
	}
}
