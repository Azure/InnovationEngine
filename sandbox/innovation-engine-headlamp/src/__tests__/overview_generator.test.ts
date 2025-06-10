// Unit tests for Overview Generator functionality
import { beforeEach, describe, expect, test, vi } from 'vitest';
import dotenv from 'dotenv';
import fs from 'fs';
import path from 'path';
import { AzureAIService } from '../services/azureAI';
import { formatSystemPrompt, loadSystemPrompt } from '../utils/promptUtils';

// Create a mock function that will be shared across all tests
const mockCreate = vi.fn();

// Mock the Azure OpenAI client
vi.mock('openai', () => {
  return {
    AzureOpenAI: vi.fn().mockImplementation(() => ({
      chat: {
        completions: {
          create: mockCreate
        }
      }
    }))
  };
});

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
    vi.clearAllMocks();
    
    // Create service instance with test config
    azureAIService = new AzureAIService({
      apiKey: 'test-api-key',
      endpoint: 'https://test-endpoint.openai.azure.com',
      deploymentId: 'test-deployment'
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
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Mocked overview content about Kubernetes resources' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    const topic = 'Kubernetes Deployments';
    const overview = await azureAIService.generateOverview(topic);
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked overview content about Kubernetes resources');
    
    // Verify that the SDK was called with the right arguments
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        messages: expect.arrayContaining([
          expect.objectContaining({
            role: 'system',
            content: expect.stringContaining('You are an Azure cloud architect')
          }),
          expect.objectContaining({
            role: 'user',
            content: 'Create an overview of: Kubernetes Deployments'
          })
        ])
      }),
      expect.any(Object)
    );
  });

  test('should use custom options when generating overview', async () => {
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Mocked overview with custom options' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    const topic = 'Kubernetes ConfigMaps';
    const overview = await azureAIService.generateOverview(topic, {
      maxTokens: 300,
      temperature: 0.3,
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked overview with custom options');
    
    // Verify that custom options were applied
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        max_tokens: 300,
        temperature: 0.3
      }),
      expect.any(Object)
    );
  });
});
