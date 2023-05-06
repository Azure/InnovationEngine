package engine

type EngineConfiguration struct {
	supportedLanguages []string `yaml:"supportedLanguages"`
}

type Engine struct {
	Configuration EngineConfiguration
}

func LoadConfiguration() EngineConfiguration {
	return EngineConfiguration{
		supportedLanguages: []string{"bash", "azurecli", "azurecli-interactive", "terraform"},
	}
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	ExecuteAndRenderSteps(scenario.Steps, scenario.Environment)
	return nil
}
