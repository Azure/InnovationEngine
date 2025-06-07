// Utility functions for loading prompt templates
import * as fs from 'fs';
import * as path from 'path';

/**
 * Load a system prompt template from a file
 * @param fileName Name of the prompt file (without path)
 * @returns The content of the prompt file as a string
 */
export function loadSystemPrompt(fileName: string): string {
  try {
    const promptPath = path.join(__dirname, '..', 'prompts', fileName);
    return fs.readFileSync(promptPath, 'utf8');
  } catch (error) {
    console.error(`Error loading prompt file "${fileName}":`, error);
    // Return a fallback prompt if the file can't be loaded
    return 'You are a helpful assistant. Please provide information based on the user\'s request.';
  }
}

/**
 * Format a system prompt with variables
 * @param promptTemplate The prompt template with {placeholders}
 * @param variables Object containing variable values to insert
 * @returns The formatted prompt with variables replaced
 */
export function formatSystemPrompt(promptTemplate: string, variables: Record<string, string> = {}): string {
  let formattedPrompt = promptTemplate;
  
  // Replace each {variable} with its value
  Object.entries(variables).forEach(([key, value]) => {
    const placeholder = `{${key}}`;
    formattedPrompt = formattedPrompt.replace(new RegExp(placeholder, 'g'), value);
  });
  
  return formattedPrompt;
}
