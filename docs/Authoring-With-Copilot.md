# Confgigure Copilot

* CTRL+SHIFT+P
* Type 'settings.json`
* Select `Preferences Open User Settings`

Add the following to the `settings.json` file that opens.

<!-- TODO: Ensure this is updated from current user settings -->

```json
// A set of instructions that will be added to Copilot requests that generate code.
// Instructions can come from: 
// - a file in the workspace: `{ "file": "fileName" }`
// - text in natural language: `{ "text": "Use underscore for field names." }`
// 
// Note: Keep your instructions short and precise. Poor instructions can degrade Copilot's quality and performance.
"github.copilot.chat.codeGeneration.instructions": [
    {
    "text": "When writing mardown files each section should have an introductory paragraph, and optional code block and a summary paragraph."
    },
],
```

# Use Copilot

Using Copilot to author docs is easy in Visual Studio Code.

## Initial Authoring

* Create a new document
* `CTRL+I`
* Type "Create an executable document which [Objective]"
* Copilot will attempt to create the document for your, providing heading titles and intro paragraphs
* Review the document, if any section is missing or needs adjustment position the cursor at that point, hit `CTRL-I`, give the instruction
* Work through the document creating the code blocks

## Testing

Once you have the document in good shape and you feel it will work you can test it with Innovation Engine.

* Hit CTRL-SHIFT-` to open a WSL terminal (Innovation Engine does not work in PowerShell)
* Type `ie test filename.md`
* The document will be executed in test mode, any failure will be reported in the terminal
* If you want Copilot assistance with errors, position the cursor in the code block where the error occurred and paste the error message
* Repeat until no errors occur