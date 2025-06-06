// API server integration tests
import request from 'supertest';
import { describe, test, expect, beforeEach, vi } from 'vitest';

// Create a mock for the AzureAIService
vi.mock('../services/azureAI', () => {
  const AzureAIService = vi.fn().mockImplementation(() => {
    return {
      getCompletion: vi.fn().mockImplementation(async (messages) => {
        return `This is a mocked response for: ${messages[messages.length - 1].content}`;
      })
    };
  });
  
  AzureAIService.loadConfig = vi.fn().mockReturnValue({
    apiKey: 'test-api-key',
    endpoint: 'https://test-endpoint.openai.azure.com',
    deploymentId: 'test-deployment'
  });
  
  return { AzureAIService };
});

// Set NODE_ENV to test to bypass environment variable checks
process.env.NODE_ENV = 'test';

// Import our server app after mocking dependencies
import app from '../server';

describe('API Server Integration', () => {
  
  // Test middleware for handling requests
  beforeEach(() => {
    vi.clearAllMocks();
  });
  
  test('should respond to health check', async () => {
    const response = await request(app).get('/api/health');
    expect(response.status).toBe(200);
    expect(response.body).toHaveProperty('status', 'ok');
  });
  
  test('should process assistant requests properly', async () => {
    const testMessages = [
      { role: 'user', content: 'Hello, how can you help me?' }
    ];
    
    const response = await request(app)
      .post('/api/assistant')
      .send({ messages: testMessages })
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json');
    
    expect(response.status).toBe(200);
    expect(response.body).toHaveProperty('response');
    expect(typeof response.body.response).toBe('string');
    
    // If we made it here without errors, the test passes
  });
  
  test('should handle missing messages in request', async () => {
    const response = await request(app)
      .post('/api/assistant')
      .send({}) // Empty request body
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json');
    
    expect(response.status).toBe(400);
    expect(response.body).toHaveProperty('error');
    expect(response.body.error).toContain('Invalid request');
  });
  
  test('should handle invalid messages format', async () => {
    const response = await request(app)
      .post('/api/assistant')
      .send({ messages: 'not an array' }) // Invalid messages format
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json');
    
    expect(response.status).toBe(400);
    expect(response.body).toHaveProperty('error');
    expect(response.body.error).toContain('Invalid request');
  });
});
