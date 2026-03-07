package controller

import (
	"encoding/json"
	"go-book-api/model"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

const dataFile = "data/books.json"

type BookStore struct {
	mu    sync.RWMutex
	books map[string]model.Book
}

var store = &BookStore{
	books: make(map[string]model.Book),
}

func init() {
	loadBooks()
}

func loadBooks() {
	store.mu.Lock()
	defer store.mu.Unlock()

	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var books []model.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return
	}

	for _, b := range books {
		store.books[b.ID] = b
	}
}

func saveBooks() {
	store.mu.RLock()
	defer store.mu.RUnlock()

	books := make([]model.Book, 0, len(store.books))
	for _, b := range store.books {
		books = append(books, b)
	}

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(dataFile, data, 0644)
}

type BookController struct{}

func (ctrl *BookController) Index(c *echo.Context) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	books := make([]model.Book, 0, len(store.books))

	author := c.QueryParam("author")

	for _, b := range store.books {
		if author != "" && !strings.Contains(strings.ToLower(b.Author), strings.ToLower(author)) {
			continue
		}
		books = append(books, b)
	}

	sort.Slice(books, func(i, j int) bool {
		return strings.ToLower(books[i].Title) < strings.ToLower(books[j].Title)
	})

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	end := start + limit

	if start > len(books) {
		books = []model.Book{}
	} else {
		if end > len(books) {
			end = len(books)
		}
		books = books[start:end]
	}

	return c.JSON(http.StatusOK, books)
}

func (ctrl *BookController) Detail(c *echo.Context) error {
	id := c.Param("id")

	store.mu.RLock()
	defer store.mu.RUnlock()

	book, ok := store.books[id]
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "Data not found"})
	}

	return c.JSON(http.StatusOK, book)
}

func (ctrl *BookController) Create(c *echo.Context) error {
	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	if input.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "title is required"})
	}
	if input.Author == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "author is required"})
	}
	if input.Year == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "year is required"})
	}

	book := model.Book{
		ID:        uuid.New().String(),
		Title:     strings.TrimSpace(input.Title),
		Author:    strings.TrimSpace(input.Author),
		Year:      input.Year,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	store.books[book.ID] = book

	go saveBooks()

	return c.JSON(http.StatusCreated, book)
}

func (ctrl *BookController) Update(c *echo.Context) error {
	id := c.Param("id")

	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	if input.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "title is required"})
	}
	if input.Author == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "author is required"})
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	book := model.Book{
		ID:        uuid.New().String(),
		Title:     strings.TrimSpace(input.Title),
		Author:    strings.TrimSpace(input.Author),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	existing, ok := store.books[id]
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "book not found"})
	}

	book.ID = existing.ID
	book.CreatedAt = existing.CreatedAt
	book.UpdatedAt = time.Now()
	store.books[id] = book

	go saveBooks()

	return c.JSON(http.StatusOK, input)
}

func (ctrl *BookController) Delete(c *echo.Context) error {
	id := c.Param("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.books[id]; !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "book not found"})
	}

	delete(store.books, id)
	go saveBooks()

	return c.JSON(http.StatusOK, map[string]any{"message": "book deleted"})
}
