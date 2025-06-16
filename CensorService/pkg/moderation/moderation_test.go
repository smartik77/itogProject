package moderation

import (
	"sync"
	"testing"
)

func TestNewCensorService(t *testing.T) {
	t.Run("creates service with initial words", func(t *testing.T) {
		words := []string{"bad", "word"}
		service := NewCensorService(words)

		if len(service.forbiddenWords) != 2 {
			t.Errorf("Expected 2 forbidden words, got %d", len(service.forbiddenWords))
		}
	})

	t.Run("creates empty service when no words provided", func(t *testing.T) {
		service := NewCensorService(nil)
		if len(service.forbiddenWords) != 0 {
			t.Error("Expected empty forbidden words list")
		}
	})
}

func TestAddForbiddenWord(t *testing.T) {
	t.Run("adds new word to list", func(t *testing.T) {
		service := NewCensorService([]string{"initial"})
		service.AddForbiddenWord("new")

		if len(service.forbiddenWords) != 2 {
			t.Error("Expected word to be added")
		}
	})

	t.Run("normalizes word before adding", func(t *testing.T) {
		service := NewCensorService(nil)
		service.AddForbiddenWord("  BAD  ")

		if service.forbiddenWords[0] != "bad" {
			t.Errorf("Expected normalized word, got '%s'", service.forbiddenWords[0])
		}
	})

	t.Run("ignores empty words", func(t *testing.T) {
		service := NewCensorService(nil)
		service.AddForbiddenWord("   ")

		if len(service.forbiddenWords) != 0 {
			t.Error("Expected empty word to be ignored")
		}
	})
}

func TestCheckText(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		text     string
		expected bool
	}{
		{"clean text passes", []string{"bad"}, "good text", true},
		{"detects forbidden word", []string{"bad"}, "this is bad", false},
		{"case insensitive", []string{"BaD"}, "this is BAD", false},
		{"partial word match", []string{"bad"}, "badminton", false},
		{"empty text", []string{"bad"}, "", true},
		{"no forbidden words", []string{}, "anything", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewCensorService(tt.words)
			result := service.CheckText(tt.text)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for text '%s'", tt.expected, result, tt.text)
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	service := NewCensorService([]string{"initial"})

	// Проверка на состояние гонки
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.AddForbiddenWord("concurrent")
			service.CheckText("some text")
		}()
	}
	wg.Wait()
}
