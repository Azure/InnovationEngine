package environments

type ScenarioConfigurations struct {
	Permissions []string `json:"permissions"`
	// These are not being picked up yet but would contain variables that are
	// found within the document and can be configured.
	Variables []string `json:"variables"`
}

type ScenarioMetadata struct {
	Key              string                 `json:"key"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	ExtraDetails     string                 `json:"extraDetails"`
	BulletPoints     []string               `json:"bulletPoints"`
	SourceURL        string                 `json:"sourceURL"`
	DocumentationURL string                 `json:"documentationURL"`
	Configurations   ScenarioConfigurations `json:"configurations"`
}

type LocalizedScenarioMetadata struct {
	Key      string                      `json:"key"`
	IsActive bool                        `json:"isActive"`
	Locales  map[string]ScenarioMetadata `json:"locales"`
}

type ScenarioMetadataCollection []LocalizedScenarioMetadata

// Resulting structure looks like:
// [
//    "key": "scenario-key",
//    "isActive": true,
//    "locales: {
//      "en": {
//        "key": "scenario-key",
//        "title": "Scenario Title",
//        "description": "Scenario Description",
//        "extraDetails": "Extra Details",
//        "bulletPoints": ["Bullet Point 1", "Bullet Point 2"],
//        "sourceURL": "https://source.url",
//        "documentationURL": "https://documentation.url",
//        "configurations": {
//          "permissions": ["permission1", "permission2"],
//          "variables": ["variable1", "variable2"]
//        }
//      }
//    }
//  }
// ]
//
