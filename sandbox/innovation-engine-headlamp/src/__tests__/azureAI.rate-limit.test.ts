// Specific tests for Azure OpenAI rate limiting handling
import { describe, test, expect, beforeEach, vi } from 'vitest';
import { AzureAIService } from '../services/azureAI';

// Mock fetch for testing
global.fetch = vi.fn();

// Extended timeout for these tests
vi.setConfig({ testTimeout: 10000 });

describe('AzureAIService Rate Limiting', () => {
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
  
  test('should handle rate limiting with exponential backoff', async () => {
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
            message: { role: 'assistant', content: 'Response after backoff' },
            finish_reason: 'stop',
            index: 0
          }],
          usage: { prompt_tokens: 10, completion_tokens: 20, total_tokens: 30 }
        })
      } as Response);
    
    // Use real setTimeout for this test to verify backoff
    const realSetTimeout = setTimeout;
    
    // Create a mock implementation that resolves faster but still respects relative timing
    vi.spyOn(global, 'setTimeout').mockImplementation((fn, delay) => {
      // Still use delay but reduce it by a factor of 10 for testing
      return realSetTimeout(fn as TimerHandler, delay ? Math.min(delay / 10, 100) : 0) as unknown as NodeJS.Timeout;
    });
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test exponential backoff' }];
    
    const startTime = Date.now();
    const result = await service.getCompletion(messages);
    const duration = Date.now() - startTime;
    
    // Verify that it took some time for backoff (but not the full amount due to our mock)
    expect(duration).toBeGreaterThan(50);
    expect(result).toBe('Response after backoff');
  });

  test('should respect Retry-After header for rate limiting', async () => {
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

  test('should fail cleanly with appropriate error when rate limit consistently exceeded', async () => {
    // Mock consistent rate limit responses
    (global.fetch as unknown as ReturnType<typeof vi.fn>)
      .mockResolvedValue({
        ok: false,
        status: 429,
        headers: new Headers({ 'Retry-After': '1' }),
        text: async () => 'Too many requests'
      } as Response);
    
    // Make setTimeout instant for this test
    vi.spyOn(global, 'setTimeout').mockImplementation((fn: any) => {
      fn();
      return 0 as any;
    });
    
    // Test messages with low maxRateLimitRetries to fail quickly
    const messages = [{ role: 'user' as const, content: 'Test rate limit failure' }];
    const options = {
      maxRateLimitRetries: 2
    };
    
    // Expect the getCompletion method to fail with a specific error message
    await expect(service.getCompletion(messages, options)).rejects.toThrow('Azure AI API rate limit exceeded');
    
    // Verify that fetch was called the correct number of times (initial + maxRateLimitRetries)
    expect(global.fetch).toHaveBeenCalledTimes(options.maxRateLimitRetries + 1);
  });

  // This test is marked as 'skipIfCI' to skip in CI environments where rate limits may be a concern
  test.skipIf(process.env.CI === 'true')('should handle real rate limiting with actual delays', async () => {
    // This test will use real timeouts to simulate actual rate limiting behavior
    // Mock multiple rate limit responses followed by success
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
        headers: new Headers({ 'Retry-After': '2' }),
        text: async () => 'Too many requests'
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'response-after-real-delays',
          choices: [{
            message: { role: 'assistant', content: 'Response after real delays' },
            finish_reason: 'stop',
            index: 0
          }],
          usage: { prompt_tokens: 10, completion_tokens: 20, total_tokens: 30 }
        })
      } as Response);
    
    // Use actual setTimeout (not mocked) for this test
    vi.restoreAllMocks();
    
    // Test messages
    const messages = [{ role: 'user' as const, content: 'Test real rate limiting delays' }];
    
    // Measure the time it takes to complete
    const startTime = Date.now();
    const result = await service.getCompletion(messages, { 
      retryDelay: 100 // Use smaller delay for test, but still measurable 
    });
    const duration = Date.now() - startTime;
    
    // Verify the result
    expect(result).toBe('Response after real delays');
    
    // Verify that it took a significant amount of time due to the delays
    // First retry: 1s, Second retry: 2s, so at least 3s total
    // But we've reduced the retryDelay for testing, so adjust expectations
    expect(duration).toBeGreaterThan(300); // Should be at least 300ms with our test settings
  });
});
