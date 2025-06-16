package moderation

import (
	"strings"
	"sync"
)

// CensorService представляет сервис цензурирования
type CensorService struct {
	forbiddenWords []string
	mu             sync.RWMutex
}

// NewCensorService создает новый экземпляр сервиса цензурирования
func NewCensorService(initialWords []string) *CensorService {
	return &CensorService{
		forbiddenWords: initialWords,
	}
}

// AddForbiddenWord добавляет новое запрещенное слово
func (cs *CensorService) AddForbiddenWord(word string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	word = strings.ToLower(strings.TrimSpace(word))
	if word != "" {
		cs.forbiddenWords = append(cs.forbiddenWords, word)
	}
}

// CheckText проверяет текст на наличие запрещенных слов
func (cs *CensorService) CheckText(text string) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	lowerText := strings.ToLower(text)
	for _, word := range cs.forbiddenWords {
		if strings.Contains(lowerText, word) {
			return false
		}
	}
	return true
}
