// Unit tests for Azure AI exec doc generation
import { describe, test, expect, beforeEach, vi } from 'vitest';
import { AzureAIService } from '../services/azureAI';
import { loadSystemPrompt } from '../utils/promptUtils';

// Mock fetch for testing
global.fetch = vi.fn();

describe('ExecDoc Generator Tests', () => {
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
          message: { role: 'assistant', content: 'Mocked executable document overview' },
          finish_reason: 'stop',
          index: 0
        }],
        usage: {
          prompt_tokens: 100,
          completion_tokens: 150,
          total_tokens: 250
        }
      })
    });
  });
  
  test('should load the execDoc prompt file', () => {
    const prompt = loadSystemPrompt('execDoc.txt');
    expect(prompt).toContain('You are an expert in Kubernetes and executable documentation');
    expect(prompt).toContain('## Major Steps');
  });
  
  test('should generate overview for Kubernetes deployment', async () => {
    const topic = 'Kubernetes executable document for: Deploy a stateful application';
    const overview = await azureAIService.generateOverview(topic, {
      systemPromptFile: 'execDoc.txt'
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked executable document overview');
    
    // Verify that fetch was called with the right arguments
    expect(global.fetch).toHaveBeenCalled();
    const fetchCall = (global.fetch as any).mock.calls[0];
    const requestBody = JSON.parse(fetchCall[1].body);
    
    // Verify system prompt is for Kubernetes executable documentation
    expect(requestBody.messages[0].role).toBe('system');
    expect(requestBody.messages[0].content).toContain('You are an expert in Kubernetes and executable documentation');
    
    // Verify user message contains the topic
    expect(requestBody.messages[1].role).toBe('user');
    expect(requestBody.messages[1].content).toContain('Create an overview of: Kubernetes executable document for: Deploy a stateful application');
  });
  
  test('should handle error when generating overview', async () => {
    // Mock a failed response
    (global.fetch as any).mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => 'Internal Server Error'
    });
    
    // Test should throw an error
    await expect(azureAIService.generateOverview('Test topic', {
      systemPromptFile: 'execDoc.txt'
    })).rejects.toThrow();
    
    // Verify that fetch was called
    expect(global.fetch).toHaveBeenCalled();
  });
});
