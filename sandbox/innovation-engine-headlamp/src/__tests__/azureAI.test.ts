// Unit tests for AzureAI service
import { describe, test, expect, beforeEach, vi } from 'vitest';
import { AzureAIService } from '../services/azureAI';

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

// Import the mocked module
import { AzureOpenAI } from 'openai';

describe('AzureAIService', () => {
  let service: AzureAIService;
  
  beforeEach(() => {
    // Clear all mocks before each test
    vi.clearAllMocks();
    
    // Create service instance with test config
    service = new AzureAIService({
      apiKey: 'test-api-key',
      endpoint: 'https://test-endpoint.openai.azure.com',
      deploymentId: 'test-deployment'
    });
  });

  test('should create an instance with the provided configuration', () => {
    expect(service).toBeDefined();
    expect(AzureOpenAI).toHaveBeenCalledWith({
      apiKey: 'test-api-key',
      endpoint: 'https://test-endpoint.openai.azure.com',
      apiVersion: '2023-12-01-preview',
      dangerouslyAllowBrowser: true
    });
  });

  test('getCompletion should make a request to Azure OpenAI API', async () => {
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Test response' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages
    const messages = [
      { role: 'system' as const, content: 'You are a helpful assistant.' },
      { role: 'user' as const, content: 'Hello, can you help me?' }
    ];
    
    // Call the getCompletion method
    const result = await service.getCompletion(messages);
    
    // Verify the result
    expect(result).toBe('Test response');
    
    // Verify that the SDK create method was called with the right parameters
    expect(mockCreate).toHaveBeenCalledTimes(1);
    expect(mockCreate).toHaveBeenCalledWith({
      model: 'test-deployment',
      messages: [
        { role: 'system', content: 'You are a helpful assistant.' },
        { role: 'user', content: 'Hello, can you help me?' }
      ],
      max_tokens: 1000,
      temperature: 0.7,
      top_p: 0.95,
      frequency_penalty: 0,
      presence_penalty: 0,
      stop: undefined
    }, {
      maxRetries: 3
    });
  });

  test('getCompletion should handle API errors', async () => {
    // Mock API error from the SDK
    const error = new Error('Azure AI API client error (401): Unauthorized access') as any;
    error.status = 401;
    mockCreate.mockRejectedValueOnce(error);
    
    // Test messages
    const messages = [
      { role: 'user' as const, content: 'Hello, can you help me?' }
    ];
    
    // Expect the getCompletion method to throw an error
    await expect(service.getCompletion(messages)).rejects.toThrow('Azure AI API client error (401)');
  });

  test('getCompletion should handle network errors', async () => {
    // Mock network error from the SDK
    const error = new Error('Network error');
    mockCreate.mockRejectedValueOnce(error);
    
    // Test messages
    const messages = [
      { role: 'user' as const, content: 'Hello?' }
    ];
    
    // Expect the getCompletion method to throw an error
    await expect(service.getCompletion(messages)).rejects.toThrow('Network error');
  });

  test('should pass options correctly to the API call', async () => {
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Custom options response' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages and custom options
    const messages = [{ role: 'user' as const, content: 'Hello with options' }];
    const options = {
      maxTokens: 2000,
      temperature: 0.5,
      topP: 0.8,
      frequencyPenalty: 0.2,
      presencePenalty: 0.2,
      stop: ['###']
    };
    
    // Call the getCompletion method with custom options
    const result = await service.getCompletion(messages, options);
    
    // Verify the result
    expect(result).toBe('Custom options response');
    
    // Verify that the SDK create method was called with custom options
    expect(mockCreate).toHaveBeenCalledTimes(1);
    expect(mockCreate).toHaveBeenCalledWith({
      model: 'test-deployment',
      messages: messages,
      max_tokens: 2000,
      temperature: 0.5,
      top_p: 0.8,
      frequency_penalty: 0.2,
      presence_penalty: 0.2,
      stop: ['###']
    }, {
      maxRetries: 3
    });
  });

  test('should configure retry settings for the SDK', async () => {
    // Mock successful response to test that retry settings are passed correctly
    const mockResponse = {
      choices: [{
        message: { content: 'Response with retry config' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Hello with retry config' }];
    
    // Call the getCompletion method with custom retry settings
    const result = await service.getCompletion(messages, { maxRetries: 5 });
    
    // Verify the result
    expect(result).toBe('Response with retry config');
    
    // Verify that the SDK create method was called with correct retry configuration
    expect(mockCreate).toHaveBeenCalledTimes(1);
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        messages: messages
      }),
      expect.objectContaining({
        maxRetries: 5
      })
    );
  });

  test('should handle SDK retry failures gracefully', async () => {
    // Mock a rate limit error that the SDK will eventually fail with after retries
    const rateLimitError = new Error('Azure AI API rate limit exceeded: Too many requests') as any;
    rateLimitError.status = 429;
    
    mockCreate.mockRejectedValue(rateLimitError);
    
    // Test messages and low retry settings to fail quickly
    const messages = [{ role: 'user' as const, content: 'Test SDK retry handling' }];
    const options = {
      maxRetries: 1
    };
    
    // Expect the getCompletion method to fail with rate limit error after SDK retries
    await expect(service.getCompletion(messages, options)).rejects.toThrow('rate limit exceeded');
    
    // Verify that the SDK create method was called with correct retry configuration
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        messages: messages
      }),
      expect.objectContaining({
        maxRetries: 1
      })
    );
  });
  
  test('should handle server errors from SDK', async () => {
    // Mock server error that the SDK will fail with
    const serverError = new Error('Azure AI API server error (500): Internal server error') as any;
    serverError.status = 500;
    
    mockCreate.mockRejectedValue(serverError);
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Hello with server error' }];
    
    // Set a low maxRetries value to make the test fail faster
    const options = {
      maxRetries: 2
    };
    
    // Expect the getCompletion method to throw a server error after SDK handles retries
    await expect(service.getCompletion(messages, options)).rejects.toThrow('server error (500)');
  });
  
  // This test verifies that retry configuration is correctly passed to the SDK
  test.runIf(!process.env.CI)('should pass retry configuration to SDK', async () => {
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Response with retry configuration' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test retry configuration' }];
    
    // Call the getCompletion method with retry config
    const result = await service.getCompletion(messages, { maxRetries: 3 });
    
    // Verify the result
    expect(result).toBe('Response with retry configuration');
    
    // Verify that the SDK properly received the retry configuration
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        messages: messages
      }),
      expect.objectContaining({
        maxRetries: 3
      })
    );
  });
});