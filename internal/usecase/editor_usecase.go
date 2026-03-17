package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/magomedcoder/gen/internal/domain"
)

type EditorUseCase struct {
	llmRepo domain.LLMRepository
}

func NewEditorUseCase(llmRepo domain.LLMRepository) *EditorUseCase {
	return &EditorUseCase{
		llmRepo: llmRepo,
	}
}

func (e *EditorUseCase) Transform(ctx context.Context, model string, text string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("пустой текст")
	}

	sessionID := time.Now().UnixNano()
	system := "Ты - редактор текста. Задача: исправь орфографию, пунктуацию и грамматику.\n" +
		"Правила:\n" +
		"- Верни ТОЛЬКО итоговый отредактированный текст, без пояснений.\n" +
		"- Сохраняй смысл; не добавляй новых фактов.\n" +
		"- Имена, числа, даты и сущности не меняй (кроме явных опечаток).\n" +
		"- Сохраняй переносы строк и структуру по смыслу.\n"

	messages := []*domain.Message{
		domain.NewMessage(sessionID, system, domain.MessageRoleSystem),
		domain.NewMessage(sessionID, wrapUserText(text), domain.MessageRoleUser),
	}

	ch, err := e.llmRepo.SendMessage(ctx, sessionID, model, messages)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for chunk := range ch {
		b.WriteString(chunk)
	}

	return strings.TrimSpace(b.String()), nil
}

func wrapUserText(text string) string {
	return "Текст:\n\n```\n" + text + "\n```"
}
