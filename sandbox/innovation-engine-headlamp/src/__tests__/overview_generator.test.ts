// Unit tests for Overview Generator functionality
import { beforeEach, describe, expect, test, vi } from 'vitest';
import dotenv from 'dotenv';
import fs from 'fs';
import path from 'path';
import { AzureAIService } from '../services/azureAI';
import { formatSystemPrompt, loadSystemPrompt } from '../utils/promptUtils';

// Mock fetch for testing
global.fetch = vi.fn();

// Load environment variables
const envPath = path.join(process.cwd(), '.env');
if (fs.existsSync(envPath)) {
  console.log(`Loading environment variables from ${envPath}`);
  dotenv.config({ path: envPath });
}

describe('Overview Generator Tests', () => {
  let azureAIService: AzureAIService;
  
  beforeEach(() => {
    // Reset mocks before each test
    vi.resetAllMocks();
    
    // Create service instance with test config
    azureAIService = new AzureAIService({
      apiKey: 'test-api-key',
      endpoint: 'https://test-endpoint.openai.azure.com',
      deploymentId: 'test-deployment'
    });
    
    // Mock successful response for tests
    (global.fetch as any).mockResolvedValue({
      ok: true,
      status: 200,
      headers: {
        get: (name: string) => name === 'Retry-After' ? '1' : null
      },
      json: async () => ({
        id: 'response-id',
        choices: [{
          message: {
            content: 'Mocked overview content about Kubernetes resources'
          },
          finish_reason: 'stop',
          index: 0
        }],
        usage: {
          prompt_tokens: 100,
          completion_tokens: 150,
          total_tokens: 250
        }
      }),
      text: async () => 'text response'
    });
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

  test('should generate overview for a topic', async () => {
    const topic = 'Kubernetes Deployments';
    const overview = await azureAIService.generateOverview(topic);
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked overview content about Kubernetes resources');
    
    // Verify that fetch was called with the right arguments
    expect(global.fetch).toHaveBeenCalled();
    const fetchCall = (global.fetch as any).mock.calls[0];
    expect(fetchCall[1].body).toContain('Create an overview of: Kubernetes Deployments');
  });

  test('should use custom options when generating overview', async () => {
    const topic = 'Kubernetes ConfigMaps';
    const overview = await azureAIService.generateOverview(topic, {
      maxTokens: 300,
      temperature: 0.3,
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked overview content about Kubernetes resources');
    
    // Verify that fetch was called with the right arguments and custom options
    expect(global.fetch).toHaveBeenCalled();
    const fetchCall = (global.fetch as any).mock.calls[0];
    const body = JSON.parse(fetchCall[1].body);
    expect(body.max_tokens).toBe(300);
    expect(body.temperature).toBe(0.3);
    expect(body.messages[1].content).toContain('Create an overview of: Kubernetes ConfigMaps');
  });
});
