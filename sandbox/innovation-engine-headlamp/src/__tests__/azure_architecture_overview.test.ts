// Azure Architecture Overview Generator Tests
import { beforeEach, describe, expect, test, vi } from 'vitest';
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

describe('Azure Architecture Overview Generator', () => {
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
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Mocked Azure architecture overview content' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    const topic = 'Web Application with SQL Database';
    const overview = await azureAIService.generateOverview(topic);
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked Azure architecture overview content');
    
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
            content: 'Create an overview of: Web Application with SQL Database'
          })
        ])
      }),
      expect.any(Object)
    );
  });

  test('should use custom options when generating Azure architecture overview', async () => {
    // Mock successful response
    const mockResponse = {
      choices: [{
        message: { content: 'Mocked Azure architecture overview with custom options' }
      }]
    };
    
    mockCreate.mockResolvedValueOnce(mockResponse);
    
    const topic = 'Microservices Architecture on AKS';
    const overview = await azureAIService.generateOverview(topic, {
      maxTokens: 2000,
      temperature: 0.4,
    });
    
    expect(overview).toBeTruthy();
    expect(typeof overview).toBe('string');
    expect(overview).toBe('Mocked Azure architecture overview with custom options');
    
    // Verify that custom options were applied
    expect(mockCreate).toHaveBeenCalledWith(
      expect.objectContaining({
        model: 'test-deployment',
        max_tokens: 2000,
        temperature: 0.4
      }),
      expect.any(Object)
    );
  });
});