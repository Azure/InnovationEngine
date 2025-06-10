// Unit tests for Azure AI exec doc generation
import { beforeEach, describe, expect, test, vi } from 'vitest';
import { AzureAIService } from '../services/azureAI';
import { loadSystemPrompt } from '../utils/promptUtils';

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

describe('ExecDoc Generator Tests', () => {
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
  
  test('should load the execDoc prompt file', () => {
    const prompt = loadSystemPrompt('execDoc.txt');
    expect(prompt).toContain('You are an expert in Kubernetes and executable documentation');
    expect(prompt).toContain('## Major Steps');
  });
  
  test('should generate overview for Kubernetes deployment', async () => {
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Mocked executable document overview' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    const topic = 'Kubernetes executable document for: Deploy a stateful application';
    const overview = await azureAIService.generateOverview(topic, {
      systemPromptFile: 'execDoc.txt'
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked executable document overview');
    
    // Verify that the SDK was called with the right arguments
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        messages: expect.arrayContaining([
          expect.objectContaining({
            role: 'system',
            content: expect.stringContaining('You are an expert in Kubernetes and executable documentation')
          }),
          expect.objectContaining({
            role: 'user',
            content: 'Create an overview of: Kubernetes executable document for: Deploy a stateful application'
          })
        ])
      }),
      expect.any(Object)
    );
  });
  
  test('should handle error when generating overview', async () => {
    // Mock an error response from the SDK
    const mockError = new Error('Azure AI API server error (500): Internal server error');
    (mockError as any).status = 500;
    
    mockCreate.mockRejectedValueOnce(mockError);
    
    // Test should throw an error
    await expect(azureAIService.generateOverview('Test topic', {
      systemPromptFile: 'execDoc.txt'
    })).rejects.toThrow('Azure AI API server error (500): Internal server error');
    
    // Verify that the SDK was called
    expect(mockCreate).toHaveBeenCalled();
  });
});
