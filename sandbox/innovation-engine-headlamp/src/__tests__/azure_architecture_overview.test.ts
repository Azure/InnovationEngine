// Azure Architecture Overview Generator Tests
import { describe, test, expect, beforeEach, vi } from 'vitest';
import { AzureAIService } from '../services/azureAI';
import { loadSystemPrompt, formatSystemPrompt } from '../utils/promptUtils';

// Mock fetch for testing
global.fetch = vi.fn();

describe('Azure Architecture Overview Generator', () => {
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
            content: 'Mocked Azure architecture overview content'
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

  test('should load Azure architecture system prompt from file', () => {
    const prompt = loadSystemPrompt('overview.txt');
    expect(prompt).toContain('You are an Azure cloud architect');
    expect(prompt).toContain('### Major Components');
  });

  test('should format Azure system prompt with variables', () => {
    const template = 'Architecture for {workload} on {platform}';
    const formatted = formatSystemPrompt(template, { 
      workload: 'Web Application',
      platform: 'Azure'
    });
    expect(formatted).toBe('Architecture for Web Application on Azure');
  });

  test('should generate Azure architecture overview for a workload', async () => {
    const topic = 'Web Application with SQL Database';
    const overview = await azureAIService.generateOverview(topic);
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked Azure architecture overview content');
    
    // Verify that fetch was called with the right arguments
    expect(global.fetch).toHaveBeenCalled();
    const fetchCall = (global.fetch as any).mock.calls[0];
    const requestBody = JSON.parse(fetchCall[1].body);
    
    // Verify system prompt is for Azure cloud architecture
    expect(requestBody.messages[0].role).toBe('system');
    expect(requestBody.messages[0].content).toContain('You are an Azure cloud architect');
    
    // Verify user message contains the topic
    expect(requestBody.messages[1].role).toBe('user');
    expect(requestBody.messages[1].content).toContain('Create an overview of: Web Application with SQL Database');
  });

  test('should use custom options when generating Azure architecture overview', async () => {
    const topic = 'Microservices Architecture on AKS';
    const overview = await azureAIService.generateOverview(topic, {
      maxTokens: 2000,
      temperature: 0.4,
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked Azure architecture overview content');
    
    // Verify that fetch was called with the right arguments and custom options
    expect(global.fetch).toHaveBeenCalled();
    const fetchCall = (global.fetch as any).mock.calls[0];
    const requestBody = JSON.parse(fetchCall[1].body);
    
    // Verify completion options were passed correctly
    expect(requestBody.max_tokens).toBe(2000);
    expect(requestBody.temperature).toBe(0.4);
    
    // Verify user message contains the topic
    expect(requestBody.messages[1].content).toContain('Create an overview of: Microservices Architecture on AKS');
  });
});