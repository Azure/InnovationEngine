package engine

import "fmt"

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

// / Create a new engine instance.
func NewEngine() *Engine {
	return &Engine{}
}

// / Executes a scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	fmt.Println(titleStyle.Render(scenario.Name))
	ExecuteAndRenderSteps(scenario.Steps, scenario.Environment)
	fmt.Printf("---Generated script---\n%s", scenario.ToShellScript())
	return nil
}
