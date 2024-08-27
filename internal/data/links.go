package data

import (
	"database/sql"
	"math/rand"
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
	  INSERT INTO links (url, code)
	  VALUES ($1, $2)
		RETURNING code;
	`

	alphabets := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	alphabetsLen := len(alphabets)
	buf := make([]byte, 7)

	for i := 0; i < 7; i++ {
		randInt := rand.Intn(alphabetsLen)
		buf[i] = alphabets[randInt]
	}

	return l.DB.QueryRow(query, link.Url, string(buf)).Scan(&link.Code)
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

type Models struct {
	Links LinksModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Links: LinksModel{DB: db},
	}
}
