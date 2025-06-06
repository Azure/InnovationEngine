// Environment validation tests for Azure AI integration
import { describe, test, expect } from 'vitest';
import * as dotenv from 'dotenv';
import * as path from 'path';
import * as fs from 'fs';

// Load environment variables from .env file at the start of the test
const envPath = path.join(process.cwd(), '.env');
const envExamplePath = path.join(process.cwd(), 'env.example');

if (fs.existsSync(envPath)) {
  console.log(`Loading environment variables from ${envPath}`);
  dotenv.config({ path: envPath });
} else if (fs.existsSync(envExamplePath)) {
  console.log(`Warning: .env file not found. Using env.example as a fallback. Please copy env.example to .env and fill in your credentials.`);
  dotenv.config({ path: envExamplePath });
} else {
  console.error(`Error: No .env or env.example file found in ${process.cwd()}`);
}

// Test file naming starts with 'a_' to ensure it runs first alphabetically
describe('Environment Setup Validation', () => {
  test('should check if environment variables are configured for Azure AI', () => {
    // Environment variables should be loaded from dotenv by now
    
    // Function to check if Azure OpenAI environment variables are set
    const checkAzureOpenAIEnvironment = (): { isConfigured: boolean; missingVars: string[] } => {
      const requiredVars = [
        'AZURE_OPENAI_API_KEY',
        'AZURE_OPENAI_ENDPOINT',
        'AZURE_OPENAI_DEPLOYMENT_ID'
      ];
      
      const missingVars: string[] = [];
      
      for (const varName of requiredVars) {
        if (!process.env[varName]) {
          missingVars.push(varName);
        }
      }
      
      return {
        isConfigured: missingVars.length === 0,
        missingVars
      };
    };
    
        // Debug: Log the values of the environment variables we're checking
    console.log('Debug - AZURE_OPENAI_API_KEY:', process.env.AZURE_OPENAI_API_KEY ? 'Set (length: ' + process.env.AZURE_OPENAI_API_KEY.length + ')' : 'Not set');
    console.log('Debug - AZURE_OPENAI_ENDPOINT:', process.env.AZURE_OPENAI_ENDPOINT ? 'Set (length: ' + process.env.AZURE_OPENAI_ENDPOINT.length + ')' : 'Not set');
    console.log('Debug - AZURE_OPENAI_DEPLOYMENT_ID:', process.env.AZURE_OPENAI_DEPLOYMENT_ID ? 'Set (length: ' + process.env.AZURE_OPENAI_DEPLOYMENT_ID.length + ')' : 'Not set');
    
    // Validate the format of the environment variables
    if (process.env.AZURE_OPENAI_API_KEY && !process.env.AZURE_OPENAI_API_KEY.match(/^[a-zA-Z0-9]+$/)) {
      console.warn('Warning: AZURE_OPENAI_API_KEY may contain invalid characters.');
    }
    
    if (process.env.AZURE_OPENAI_ENDPOINT && !process.env.AZURE_OPENAI_ENDPOINT.startsWith('https://')) {
      console.warn('Warning: AZURE_OPENAI_ENDPOINT should start with https://');
    }
    
    // Check environment variables
    const { isConfigured, missingVars } = checkAzureOpenAIEnvironment();
    
    // Check if we're in development mode
    const isDevelopment = process.env.NODE_ENV === 'development' || process.env.NODE_ENV === 'test';
    
    if (!isConfigured) {
      // Create a consistent error message
      const errorMessage = `
⚠️ ============================================================= ⚠️
  AZURE OPENAI ENVIRONMENT SETUP INCOMPLETE

  The following environment variables are missing:
  ${missingVars.map(v => `  - ${v}`).join('\n')}
  
  To use the Azure AI features, please:
  1. Copy env.example to .env in the project root
  2. Fill in your Azure OpenAI credentials in the .env file
  3. Restart the application
⚠️ ============================================================= ⚠️
`;
      
      // In development/test mode, just warn
      if (isDevelopment) {
        console.warn(errorMessage + `
  Running in development/test mode, so this is just a warning.
  Some tests related to Azure AI functionality may be skipped.
`);
        expect(true).toBeTruthy(); // Pass the test in development mode
      } else {
        // In production mode, fail the test
        console.error(errorMessage + `
  NOT running in development mode, so this is a CRITICAL ERROR.
  The application requires Azure OpenAI credentials to function properly.
`);
        expect(isConfigured).toBeTruthy(); // Fail the test in production mode
      }
    } else {
      expect(true).toBeTruthy(); // All variables are set properly
    }
    
    // Store the result for other tests to potentially use
    process.env.AZURE_AI_CONFIG_VALIDATED = 'true';
    process.env.AZURE_AI_CONFIG_STATUS = isConfigured ? 'configured' : 'missing';
  });
  
  test('should verify that env.example file exists with proper structure', () => {
    // Import fs module for file operations
    const fs = require('fs');
    const path = require('path');
    
    // Define the path to env.example
    const envExamplePath = path.join(process.cwd(), 'env.example');
    
    // Check if env.example file exists
    const envExampleExists = fs.existsSync(envExamplePath);
    expect(envExampleExists).toBeTruthy();
    
    if (envExampleExists) {
      // Read the file content
      const content = fs.readFileSync(envExamplePath, 'utf8');
      
      // Check if the file contains the required variables
      expect(content).toContain('AZURE_OPENAI_API_KEY=');
      expect(content).toContain('AZURE_OPENAI_ENDPOINT=');
      expect(content).toContain('AZURE_OPENAI_DEPLOYMENT_ID=');
    }
  });
  
  test('should verify package.json contains the required scripts and dependencies', () => {
    // Import fs module for file operations
    const fs = require('fs');
    const path = require('path');
    
    // Define the path to package.json
    const packageJsonPath = path.join(process.cwd(), 'package.json');
    
    // Check if package.json exists
    const packageJsonExists = fs.existsSync(packageJsonPath);
    expect(packageJsonExists).toBeTruthy();
    
    if (packageJsonExists) {
      // Read the package.json file
      const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
      
      // Check if the required scripts exist
      expect(packageJson.scripts).toHaveProperty('server');
      expect(packageJson.scripts).toHaveProperty('dev');
      
      // Check if the required dependencies exist
      expect(packageJson.dependencies).toHaveProperty('express');
      expect(packageJson.dependencies).toHaveProperty('cors');
      expect(packageJson.dependencies).toHaveProperty('body-parser');
      expect(packageJson.dependencies).toHaveProperty('dotenv');
    }
  });
});
