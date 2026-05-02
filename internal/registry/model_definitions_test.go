package registry

import "testing"

func hasModel(models []*ModelInfo, modelID string) bool {
	for _, model := range models {
		if model != nil && model.ID == modelID {
			return true
		}
	}
	return false
}

func TestGitHubCopilotGeminiModelsAreChatOnly(t *testing.T) {
	models := GetGitHubCopilotModels()
	required := map[string]bool{
		"gemini-2.5-pro":         false,
		"gemini-3-pro-preview":   false,
		"gemini-3.1-pro-preview": false,
		"gemini-3-flash-preview": false,
	}

	for _, model := range models {
		if _, ok := required[model.ID]; !ok {
			continue
		}
		required[model.ID] = true
		if len(model.SupportedEndpoints) != 1 || model.SupportedEndpoints[0] != "/chat/completions" {
			t.Fatalf("model %q supported endpoints = %v, want [/chat/completions]", model.ID, model.SupportedEndpoints)
		}
	}

	for modelID, found := range required {
		if !found {
			t.Fatalf("expected GitHub Copilot model %q in definitions", modelID)
		}
	}
}

func TestCodexCatalogIncludesLatestCodexModels(t *testing.T) {
	cases := []struct {
		name   string
		models []*ModelInfo
	}{
		{name: "codex-free", models: GetCodexFreeModels()},
		{name: "codex-team", models: GetCodexTeamModels()},
		{name: "codex-plus", models: GetCodexPlusModels()},
		{name: "codex-pro", models: GetCodexProModels()},
		{name: "codex-static-channel", models: GetStaticModelDefinitionsByChannel("codex")},
	}

	requiredModels := []string{"gpt-5.5", "codex-auto-review"}
	for _, tc := range cases {
		for _, modelID := range requiredModels {
			if !hasModel(tc.models, modelID) {
				t.Fatalf("%s is missing %s", tc.name, modelID)
			}
		}
	}
}

func TestEmbeddedModelsCatalogIsValid(t *testing.T) {
	if err := loadModelsFromBytes(embeddedModelsJSON, "embed-test"); err != nil {
		t.Fatalf("embedded models catalog is invalid: %v", err)
	}
}
