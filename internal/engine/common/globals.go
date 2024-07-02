package common

import tea "github.com/charmbracelet/bubbletea"

// TODO: Ideally we won't need a global program variable. We should
// refactor this in the future such that each tea program is localized to the
// function that creates it and ExecuteCodeBlockSync doesn't mutate the global
// program variable.
var Program *tea.Program = nil
