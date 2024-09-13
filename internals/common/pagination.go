package common

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Paginable interface {
	GetLimit() int
	GetOffset() int
	GetCurrentURL() string
}

type StandardPaginable struct {
	Limit      int
	Offset     int
	CurrentURL string
}

func (s StandardPaginable) GetOffset() int {
	return s.Offset
}

func (s StandardPaginable) GetLimit() int {
	return s.Limit
}

func (s StandardPaginable) GetCurrentURL() string {
	return s.CurrentURL
}

func GetCurrentURL(c *fiber.Ctx) string {
	return c.BaseURL() + c.Path()
}

type Pagination struct {
	Count          int64  `json:"count"`
	NextOffset     int    `json:"-"`
	PreviousOffset int    `json:"-"`
	Next           string `json:"next"`
	Previous       string `json:"previous"`
}

type PageParams struct {
	Limit  int
	Offset int
	Count  int64
	URL    string
}

func getNext(p *PageParams) int {
	next := p.Offset + p.Limit
	if int64(next) >= p.Count {
		next = 0
	}

	return next
}

func getPrevious(p *PageParams) int {
	pc := p.Offset - p.Limit
	if pc <= 0 {
		return 0
	} else {
		return pc
	}
}

func CreatePagination(pp *PageParams) Pagination {
	next := getNext(pp)
	prev := getPrevious(pp)

	var nextURL string
	var prevURL string

	if next > 0 {
		nextURL = fmt.Sprintf("%s?limit=%d&offset=%d", pp.URL, pp.Limit, next)
	}

	if prev >= 0 {
		prevURL = fmt.Sprintf("%s?limit=%d&offset=%d", pp.URL, pp.Limit, prev)
	}

	return Pagination{
		NextOffset:     next,
		PreviousOffset: prev,
		Count:          pp.Count,
		Next:           nextURL,
		Previous:       prevURL,
	}
}
