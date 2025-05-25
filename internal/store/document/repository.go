package document

import (
	"fmt"
	"strings"
	"sync"
)

type Document struct {
	ID      string            `json:"id"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Tags    []string          `json:"tags"`
	Meta    map[string]string `json:"meta,omitempty"`
}

type Repository struct {
	documents map[string]Document
	mu        sync.RWMutex
	counter   int
}

func NewRepository() *Repository {
	return &Repository{
		documents: make(map[string]Document),
	}
}

func (r *Repository) Add(doc Document) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if doc.ID == "" {
		r.counter++
		doc.ID = fmt.Sprintf("doc_%d", r.counter)
	}

	r.documents[doc.ID] = doc
	return doc.ID, nil
}

func (r *Repository) Get(id string) (Document, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doc, exists := r.documents[id]
	return doc, exists
}

func (r *Repository) List() []Document {
	r.mu.RLock()
	defer r.mu.RUnlock()

	docs := make([]Document, 0, len(r.documents))
	for _, doc := range r.documents {
		docs = append(docs, doc)
	}
	return docs
}

func (r *Repository) SearchByKeyword(query string) []Document {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []Document
	for _, doc := range r.documents {
		if containsIgnoreCase(doc.Title, query) || containsIgnoreCase(doc.Content, query) {
			results = append(results, doc)
		}
	}
	return results
}

func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
