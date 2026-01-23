package generator

import (
	"generate_api_docs_mLua/pkg/document"
	"strings"
	"testing"
)

func TestRenderHandlerDoc(t *testing.T) {
	typeLinks := make(TypeLinkInfo)

	tests := []struct {
		name           string
		handler        document.HandlerDoc
		expectBadge    string
		expectExtraRow bool
		extraRowText   string
	}{
		{
			name: "Entity badge only",
			handler: document.HandlerDoc{
				Name:             "OnEntityEvent",
				Description:      "Handler for entity events",
				ExecSpace:        "ServerOnly",
				EventSenderType:  "Entity",
				EventSenderValue: "ignored",
			},
			expectBadge:    "Entity",
			expectExtraRow: false,
		},
		{
			name: "Model badge only",
			handler: document.HandlerDoc{
				Name:             "OnModelEvent",
				Description:      "Handler for model events",
				ExecSpace:        "ClientOnly",
				EventSenderType:  "Model",
				EventSenderValue: "ignored",
			},
			expectBadge:    "Model",
			expectExtraRow: false,
		},
		{
			name: "Logic badge with extra info",
			handler: document.HandlerDoc{
				Name:             "OnLogicEvent",
				Description:      "Handler for logic events",
				ExecSpace:        "Server",
				EventSenderType:  "Logic",
				EventSenderValue: "MyLogicName",
			},
			expectBadge:    "Logic",
			expectExtraRow: true,
			extraRowText:   "MyLogicName",
		},
		{
			name: "Service badge with extra info",
			handler: document.HandlerDoc{
				Name:             "OnServiceEvent",
				Description:      "Handler for service events",
				ExecSpace:        "Client",
				EventSenderType:  "Service",
				EventSenderValue: "MyServiceName",
			},
			expectBadge:    "Service",
			expectExtraRow: true,
			extraRowText:   "MyServiceName",
		},
		{
			name: "LocalPlayer badge only",
			handler: document.HandlerDoc{
				Name:            "OnPlayerEvent",
				Description:     "Handler for player events",
				ExecSpace:       "All",
				EventSenderType: "LocalPlayer",
			},
			expectBadge:    "LocalPlayer",
			expectExtraRow: false,
		},
		{
			name: "Self badge only",
			handler: document.HandlerDoc{
				Name:            "OnSelfEvent",
				Description:     "Handler for self events",
				EventSenderType: "Self",
			},
			expectBadge:    "Self",
			expectExtraRow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := renderHandlerDoc(tt.handler, typeLinks)

			// Check if badge is present
			if !strings.Contains(html, tt.expectBadge) {
				t.Errorf("Expected badge %s not found in output", tt.expectBadge)
			}

			// Check for extra row with Logic/Service name
			if tt.expectExtraRow {
				if !strings.Contains(html, tt.extraRowText) {
					t.Errorf("Expected extra row text %s not found in output", tt.extraRowText)
				}
			}
		})
	}
}

func TestRenderHandlerDocWithoutEventSender(t *testing.T) {
	typeLinks := make(TypeLinkInfo)

	handler := document.HandlerDoc{
		Name:        "OnSimpleEvent",
		Description: "Simple handler without EventSender",
		ExecSpace:   "ServerOnly",
	}

	html := renderHandlerDoc(handler, typeLinks)

	// Should not contain EventSender badges
	eventSenderBadges := []string{"Entity", "Model", "Logic", "Service", "LocalPlayer", "Self"}
	for _, badge := range eventSenderBadges {
		// Check if the badge image is in the output (not just the text)
		badgeImg := `badge/` + badge
		if strings.Contains(html, badgeImg) {
			t.Errorf("Unexpected EventSender badge %s found in output", badge)
		}
	}

	// Should contain ExecSpace badge
	if !strings.Contains(html, "ServerOnly") {
		t.Error("Expected ServerOnly badge not found in output")
	}
}
