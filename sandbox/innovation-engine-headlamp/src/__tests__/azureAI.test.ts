// Unit tests for AzureAI service
import { describe, test, expect, beforeEach, vi } from 'vitest';
import { AzureAIService } from '../services/azureAI';

// Mock fetch for testing
global.fetch = vi.fn();

describe('AzureAIService', () => {
  let service: AzureAIService;
  
  beforeEach(() => {
    // Reset mocks before each test
    vi.resetAllMocks();
    
    // Create service instance with test config
    service = new AzureAIService({
      apiKey: 'test-api-key',
      endpoint: 'https://test-endpoint.openai.azure.com',
      deploymentId: 'test-deployment'
    });
  });
  
  test('should create an instance with the provided configuration', () => {
    expect(service).toBeDefined();
  });
  
  test('getCompletion should make a request to Azure OpenAI API', async () => {
    // Mock successful response
    const mockResponse = {
      id: 'response-id',
      choices: [{
        message: { role: 'assistant', content: 'Test response' },
        finish_reason: 'stop',
        index: 0
      }],
      usage: {
        prompt_tokens: 10,
        completion_tokens: 20,
        total_tokens: 30
      }
    };
    
    // Setup the mock fetch
    (global.fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse
    } as Response);
    
    // Test messages with explicitly typed role
    const messages = [
      { role: 'system' as const, content: 'You are a helpful assistant.' },
      { role: 'user' as const, content: 'Hello, can you help me?' }
    ];
    
    // Call the getCompletion method
    const result = await service.getCompletion(messages);
    
    // Verify the result
    expect(result).toBe('Test response');
    
    // Verify that fetch was called with the right parameters
    expect(global.fetch).toHaveBeenCalledTimes(1);
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('https://test-endpoint.openai.azure.com/openai/deployments/test-deployment/chat/completions'),
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({
          'Content-Type': 'application/json',
          'api-key': 'test-api-key'
        })
      })
    );
  });
  
  test('getCompletion should handle API errors', async () => {
    // Mock error response
    (global.fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 401,
      text: async () => 'Unauthorized access'
    } as Response);
    
    // Test messages with explicitly typed role
    const messages = [
      { role: 'user' as const, content: 'Hello, can you help me?' }
    ];
    
    // Expect the getCompletion method to throw an error
    await expect(service.getCompletion(messages)).rejects.toThrow('Azure AI API error: 401');
  });
  
  test('getCompletion should handle network errors', async () => {
    // Mock network error
    (global.fetch as unknown as ReturnType<typeof vi.fn>).mockRejectedValueOnce(new Error('Network error'));
    
    // Test messages with explicitly typed role
    const messages = [
      { role: 'user' as const, content: 'Hello?' }
    ];
    
    // Expect the getCompletion method to throw an error
    await expect(service.getCompletion(messages)).rejects.toThrow('Network error');
  });
  
  test('should pass options correctly to the API call', async () => {
    // Mock successful response
    const mockResponse = {
      id: 'response-id',
      choices: [{
        message: { role: 'assistant', content: 'Custom options response' },
        finish_reason: 'stop',
        index: 0
      }],
      usage: {
        prompt_tokens: 10,
        completion_tokens: 20,
        total_tokens: 30
      }
    };
    
    // Setup the mock fetch
    (global.fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse
    } as Response);
    
    // Test messages and custom options with explicitly typed role
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
    
    // Verify that fetch was called with the right parameters
    expect(global.fetch).toHaveBeenCalledTimes(1);
    
    // Check that the request body includes the custom options
    const fetchCallArgs = (global.fetch as unknown as ReturnType<typeof vi.fn>).mock.calls[0][1];
    const requestBody = JSON.parse(fetchCallArgs.body as string);
    
    expect(requestBody).toEqual({
      messages: messages,
      max_tokens: 2000,
      temperature: 0.5,
      top_p: 0.8,
      frequency_penalty: 0.2,
      presence_penalty: 0.2,
      stop: ['###']
    });
  });
  
  test('should retry on rate limiting (429) responses', async () => {
    // Mock rate limit response followed by success
    (global.fetch as unknown as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'response-after-retry',
          choices: [{
            message: { role: 'assistant', content: 'Response after retry' },
            finish_reason: 'stop',
            index: 0
          }],
          usage: { prompt_tokens: 10, completion_tokens: 20, total_tokens: 30 }
        })
      } as Response);
    
    // Spy on setTimeout to make the test run faster
    vi.spyOn(global, 'setTimeout').mockImplementation((fn: any) => {
      fn();
      return 0 as any;
    });
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Hello with rate limit' }];
    
    // Call the getCompletion method
    const result = await service.getCompletion(messages);
    
    // Verify the result
    expect(result).toBe('Response after retry');
    
    // Verify that fetch was called twice (initial + retry)
    expect(global.fetch).toHaveBeenCalledTimes(2);
  });

  test('should fail after exceeding maximum rate limit retries', async () => {
    // Mock multiple rate limit responses to exceed the maxRateLimitRetries
    (global.fetch as unknown as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response);
    
    // Spy on setTimeout to make the test run faster
    vi.spyOn(global, 'setTimeout').mockImplementation((fn: any) => {
      fn();
      return 0 as any;
    });
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Hello with excessive rate limit' }];
    
    // Set a low maxRateLimitRetries value to make the test fail faster
    const options = {
      maxRateLimitRetries: 2
    };
    
    // Expect the getCompletion method to throw a rate limit error after exceeding retries
    await expect(service.getCompletion(messages, options)).rejects.toThrow('rate limit exceeded');
    
    // Verify that fetch was called the expected number of times (initial + maxRateLimitRetries)
    expect(global.fetch).toHaveBeenCalledTimes(options.maxRateLimitRetries + 1);
  });

  test('should respect configurable retry settings', async () => {
    // Mock a successful response after several failures
    (global.fetch as unknown as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce({
        ok: false,
        status: 500,
        text: async () => 'Internal server error'
      } as Response)
      .mockResolvedValueOnce({
        ok: false,
        status: 500,
        text: async () => 'Internal server error'
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'response-after-server-errors',
          choices: [{
            message: { role: 'assistant', content: 'Response after multiple server errors' },
            finish_reason: 'stop',
            index: 0
          }],
          usage: { prompt_tokens: 10, completion_tokens: 20, total_tokens: 30 }
        })
      } as Response);
      
    // Spy on setTimeout to verify it's called with the correct delay
    const setTimeoutSpy = vi.spyOn(global, 'setTimeout').mockImplementation((fn: any) => {
      fn();
      return 0 as any;
    });
    
    // Test messages and custom retry settings
    const messages = [{ role: 'user' as const, content: 'Test with custom retry settings' }];
    const options = {
      maxRetries: 5,
      retryDelay: 2000
    };
    
    // Call the getCompletion method with custom options
    const result = await service.getCompletion(messages, options);
    
    // Verify the result
    expect(result).toBe('Response after multiple server errors');
    
    // Verify that fetch was called the expected number of times
    expect(global.fetch).toHaveBeenCalledTimes(3);
    
    // Verify that setTimeout was called with the expected delays
    // First retry: 2000ms * 1, Second retry: 2000ms * 2
    expect(setTimeoutSpy).toHaveBeenCalledTimes(2);
    expect(setTimeoutSpy).toHaveBeenNthCalledWith(1, expect.any(Function), 2000);
    expect(setTimeoutSpy).toHaveBeenNthCalledWith(2, expect.any(Function), 4000);
  });
  
  test('should fail after exceeding maximum rate limit retries', async () => {
    // Mock multiple rate limit responses to exceed the maxRateLimitRetries
    (global.fetch as unknown as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response);
    
    // Spy on setTimeout to make the test run faster
    vi.spyOn(global, 'setTimeout').mockImplementation((fn: any) => {
      fn();
      return 0 as any;
    });
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Hello with excessive rate limit' }];
    
    // Set a low maxRateLimitRetries value to make the test fail faster
    const options = {
      maxRateLimitRetries: 2
    };
    
    // Expect the getCompletion method to throw a rate limit error after exceeding retries
    await expect(service.getCompletion(messages, options)).rejects.toThrow('rate limit exceeded');
    
    // Verify that fetch was called the expected number of times (initial + maxRateLimitRetries)
    expect(global.fetch).toHaveBeenCalledTimes(options.maxRateLimitRetries + 1);
  });
  
  // This test is marked to skip in CI environments where rate limits may be a concern
  test.runIf(!process.env.CI)('should respect Retry-After header for rate limiting', async () => {
    // Mock rate limit response with a specific retry-after header
    (global.fetch as unknown as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '2' }), // 2 seconds
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'response-after-specific-delay',
          choices: [{
            message: { role: 'assistant', content: 'Response after specific delay' },
            finish_reason: 'stop',
            index: 0
          }],
          usage: { prompt_tokens: 10, completion_tokens: 20, total_tokens: 30 }
        })
      } as Response);
    
    // Spy on setTimeout to verify it's called with the correct delay from Retry-After
    const setTimeoutSpy = vi.spyOn(global, 'setTimeout').mockImplementation((fn: any) => {
      fn();
      return 0 as any;
    });
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test Retry-After header' }];
    
    // Call the getCompletion method
    const result = await service.getCompletion(messages);
    
    // Verify the result
    expect(result).toBe('Response after specific delay');
    
    // Verify that setTimeout was called with the expected delay (2000ms from Retry-After header)
    expect(setTimeoutSpy).toHaveBeenCalledWith(expect.any(Function), 2000);
  });
});
