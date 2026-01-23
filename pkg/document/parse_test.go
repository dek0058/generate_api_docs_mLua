package document

import (
	"testing"
)

func TestEventSenderParsing(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedType      string
		expectedValue     string
		expectedBadgeType string
	}{
		{
			name: "Entity with ignored value",
			input: `@Logic
---@description "Test handler"
@EventSender("Entity", "IgnoredValue")
handler TestHandler()`,
			expectedType:      "Entity",
			expectedValue:     "IgnoredValue",
			expectedBadgeType: "Entity",
		},
		{
			name: "Model with ignored value",
			input: `@Logic
---@description "Test handler"
@EventSender("Model", "IgnoredValue")
handler TestHandler()`,
			expectedType:      "Model",
			expectedValue:     "IgnoredValue",
			expectedBadgeType: "Model",
		},
		{
			name: "Logic with LogicName",
			input: `@Logic
---@description "Test handler"
@EventSender("Logic", "MyLogicName")
handler TestHandler()`,
			expectedType:      "Logic",
			expectedValue:     "MyLogicName",
			expectedBadgeType: "Logic",
		},
		{
			name: "Service with ServiceName",
			input: `@Logic
---@description "Test handler"
@EventSender("Service", "MyServiceName")
handler TestHandler()`,
			expectedType:      "Service",
			expectedValue:     "MyServiceName",
			expectedBadgeType: "Service",
		},
		{
			name: "LocalPlayer without second parameter",
			input: `@Logic
---@description "Test handler"
@EventSender("LocalPlayer")
handler TestHandler()`,
			expectedType:      "LocalPlayer",
			expectedValue:     "",
			expectedBadgeType: "LocalPlayer",
		},
		{
			name: "Self without second parameter",
			input: `@Logic
---@description "Test handler"
@EventSender("Self")
handler TestHandler()`,
			expectedType:      "Self",
			expectedValue:     "",
			expectedBadgeType: "Self",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if len(doc.Handlers) == 0 {
				t.Fatalf("Expected at least one handler, got none")
			}

			handler := doc.Handlers[0]
			if handler.EventSenderType != tt.expectedType {
				t.Errorf("EventSenderType = %v, want %v", handler.EventSenderType, tt.expectedType)
			}

			if handler.EventSenderValue != tt.expectedValue {
				t.Errorf("EventSenderValue = %v, want %v", handler.EventSenderValue, tt.expectedValue)
			}
		})
	}
}

func TestEventSenderWithoutAnnotation(t *testing.T) {
	input := `@Logic
---@description "Test handler"
handler TestHandler()`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(doc.Handlers) == 0 {
		t.Fatalf("Expected at least one handler, got none")
	}

	handler := doc.Handlers[0]
	if handler.EventSenderType != "" {
		t.Errorf("EventSenderType = %v, want empty string", handler.EventSenderType)
	}

	if handler.EventSenderValue != "" {
		t.Errorf("EventSenderValue = %v, want empty string", handler.EventSenderValue)
	}
}
