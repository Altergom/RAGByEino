package prompt

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var Template *prompt.DefaultChatTemplate

func NewTemplate() *prompt.DefaultChatTemplate {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.MessagesPlaceholder("history_key", false),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我{task}",
		},
	)

	return template
}

func SetTemplate(role, task string, history []*schema.Message) map[string]any {
	return map[string]any{
		"role":        role,
		"task":        task,
		"history_key": history,
	}
}
