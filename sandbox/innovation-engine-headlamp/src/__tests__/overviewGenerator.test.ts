import { beforeAll, describe, expect, test } from 'vitest';
import dotenv from 'dotenv';
import fs from 'fs';
import path from 'path';
import { AzureAIService } from '../services/azureAI';
import { formatSystemPrompt, loadSystemPrompt } from '../utils/promptUtils';
import { loadSystemPrompt, formatSystemPrompt } from '../utils/promptUtils';

// Load environment variables
const envPath = path.join(process.cwd(), '.env');
if (fs.existsSync(envPath)) {
  console.log(`Loading environment variables from ${envPath}`);
  dotenv.config({ path: envPath });
}

// Skip tests if environment variables are not set
const runTests = process.env.AZURE_OPENAI_API_KEY &&
                process.env.AZURE_OPENAI_ENDPOINT &&
                process.env.AZURE_OPENAI_DEPLOYMENT_ID;

describe('Overview Generator Tests', () => {
  let azureAIService: AzureAIService;
  
  beforeAll(() => {
    if (runTests) {
      azureAIService = new AzureAIService({
        apiKey: process.env.AZURE_OPENAI_API_KEY!,
        endpoint: process.env.AZURE_OPENAI_ENDPOINT!,
        deploymentId: process.env.AZURE_OPENAI_DEPLOYMENT_ID!,
      });
    }
  });

  test('should load system prompt from file', () => {
    const prompt = loadSystemPrompt('overview.txt');
    expect(prompt).toContain('You are an Azure cloud architect');
  });

  test('should format system prompt with variables', () => {
    const template = 'Hello {name}, welcome to {service}!';
    const formatted = formatSystemPrompt(template, { name: 'User', service: 'Kubernetes' });
    expect(formatted).toBe('Hello User, welcome to Kubernetes!');
  });

  // Only run API tests if environment variables are configured
  (runTests ? test : test.skip)('should generate overview for a topic', async () => {
    // This test will be skipped if environment variables are not set
    const topic = 'Kubernetes Deployments';
    const overview = await azureAIService.generateOverview(topic);
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview.length).toBeGreaterThan(100);
    expect(overview).toMatch(/deployment|Deployment|DEPLOYMENT/i);
  }, 30000); // 30 second timeout for AI API call

  (runTests ? test : test.skip)('should use custom options when generating overview', async () => {
    const topic = 'Kubernetes ConfigMaps';
    const overview = await azureAIService.generateOverview(topic, {
      maxTokens: 300,
      temperature: 0.3,
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    // With lower max tokens, we expect a shorter response
    expect(overview.length).toBeLessThan(2500);
  }, 30000);
});
