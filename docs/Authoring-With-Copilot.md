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

<!-- TODO: implement the workflow as documented -->

* Create a new document
* `CTRL+I`
* Type "Outline an executable document which [Objective]"
* Copilot will attempt to outline the document for your, providing heading titles and intro paragraphs
* Review the document, if any section is missing or needs adjustment position the cursor at that point, hit `CTRL-I`, give the instruction
* Work through the document creating the code blocks 