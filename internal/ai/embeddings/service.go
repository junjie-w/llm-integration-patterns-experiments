package embeddings

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/junjie-w/llm-integration-patterns-experiments/internal/store/document"
	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	client *openai.Client
	cache  map[string][]float32
}

func NewService(cfg *config.Config) *Service {
	client := openai.NewClient(cfg.OpenAIKey)
	return &Service{
		client: client,
		cache:  make(map[string][]float32),
	}
}

func (s *Service) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}

func (s *Service) IndexDocument(ctx context.Context, doc document.Document) error {
	embedding, err := s.GetEmbedding(ctx, doc.Content)
	if err != nil {
		return fmt.Errorf("failed to get embedding for document %s: %w", doc.ID, err)
	}

	s.cache[doc.ID] = embedding
	return nil
}

type SimilarityResult struct {
	Document  document.Document
	Score     float32
	Embedding []float32
}

func (s *Service) FindSimilarDocuments(ctx context.Context, query string, docs []document.Document, limit int) ([]SimilarityResult, error) {
	queryEmbedding, err := s.GetEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query embedding: %w", err)
	}

	var results []SimilarityResult

	for _, doc := range docs {
		if _, exists := s.cache[doc.ID]; !exists {
			if err := s.IndexDocument(ctx, doc); err != nil {
				return nil, err
			}
		}

		docEmbedding := s.cache[doc.ID]
		score := cosineSimilarity(queryEmbedding, docEmbedding)

		results = append(results, SimilarityResult{
			Document:  doc,
			Score:     score,
			Embedding: docEmbedding,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func cosineSimilarity(a, b []float32) float32 {
	var dotProduct float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt(normA) * sqrt(normB))
}

func sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
