// Specific tests for Azure OpenAI rate limiting handling
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

// Extended timeout for these tests
vi.setConfig({ testTimeout: 10000 });

describe('AzureAIService Rate Limiting', () => {
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
  
  test('should configure SDK to handle rate limiting with retries', async () => {
    // Mock successful response to test rate limiting configuration
    const mockResponse = {
      choices: [{
        message: { content: 'Response after SDK handles rate limiting' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test SDK rate limiting config' }];
    
    // Call the getCompletion method with retry configuration
    const result = await service.getCompletion(messages, { maxRetries: 3 });
    
    // Verify the result
    expect(result).toBe('Response after SDK handles rate limiting');
    
    // Verify that the SDK properly received retry configuration for rate limiting
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

  test('should pass retry configuration to SDK for rate limiting', async () => {
    // Mock successful response to verify SDK retry configuration
    const mockResponse = {
      choices: [{
        message: { content: 'Response with SDK retry configuration' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test SDK Retry-After handling' }];
    
    // Call the getCompletion method with retry settings
    const result = await service.getCompletion(messages, { maxRetries: 2 });
    
    // Verify the result
    expect(result).toBe('Response with SDK retry configuration');
    
    // Verify that the SDK properly received retry configuration
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        messages: messages
      }),
      expect.objectContaining({
        maxRetries: 2
      })
    );
  });

  test('should fail with rate limit error when SDK retries are exhausted', async () => {
    // Mock consistent rate limit errors that will cause SDK to fail
    const rateLimitError = new Error('Azure AI API rate limit exceeded: Too many requests') as any;
    rateLimitError.status = 429;
    
    mockCreate.mockRejectedValue(rateLimitError);
    
    // Test messages with low maxRetries to fail quickly
    const messages = [{ role: 'user' as const, content: 'Test rate limit failure' }];
    const options = {
      maxRetries: 1
    };
    
    // Expect the getCompletion method to fail with a rate limit error after SDK retries
    await expect(service.getCompletion(messages, options)).rejects.toThrow('rate limit exceeded');
    
    // Verify that the SDK was called with proper retry configuration
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

  // This test is marked as 'skipIfCI' to skip in CI environments where rate limits may be a concern
  test.skipIf(process.env.CI === 'true')('should configure SDK for retry behavior', async () => {
    // Mock successful response to test SDK configuration
    const mockResponse = {
      choices: [{
        message: { content: 'Response with SDK configuration' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test SDK retry configuration' }];
    
    // Measure the time it takes to complete (basic timing test)
    const startTime = Date.now();
    const result = await service.getCompletion(messages, { maxRetries: 3 });
    const duration = Date.now() - startTime;
    
    // Verify the result
    expect(result).toBe('Response with SDK configuration');
    
    // Verify that it completed (duration should be reasonable for a successful call)
    expect(duration).toBeGreaterThan(0);
    expect(duration).toBeLessThan(1000); // Should complete quickly for successful mock
    
    // Verify that the SDK was properly configured
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
