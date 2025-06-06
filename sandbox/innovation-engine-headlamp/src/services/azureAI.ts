// Azure AI service for the Innovation Engine Assistant
// This service handles communication with the Azure AI API

interface AzureAIConfig {
  apiKey: string;
  endpoint: string;
  deploymentId: string;
}

interface ChatMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

interface ChatCompletionRequest {
  messages: ChatMessage[];
  max_tokens?: number;
  temperature?: number;
  top_p?: number;
  frequency_penalty?: number;
  presence_penalty?: number;
  stop?: string[];
}

interface ChatCompletionResponse {
  id: string;
  choices: {
    message: ChatMessage;
    finish_reason: string;
    index: number;
  }[];
  usage: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

/**
 * A service class to interact with Azure OpenAI API
 */
export class AzureAIService {
  private apiKey: string;
  private endpoint: string;
  private deploymentId: string;
  private apiVersion = '2023-05-15'; // Update this to the latest API version as needed

  /**
   * Create a new Azure AI Service client
   * @param config Configuration for the Azure AI service
   */
  constructor(config: AzureAIConfig) {
    this.apiKey = config.apiKey;
    this.endpoint = config.endpoint;
    this.deploymentId = config.deploymentId;
  }

  /**
   * Get a completion from Azure OpenAI API with advanced retry logic
   * @param messages Array of messages representing the conversation
   * @param options Optional parameters for the completion
   * @returns A promise with the completion response
   */
  async getCompletion(
    messages: ChatMessage[],
    options: {
      maxTokens?: number;
      temperature?: number;
      topP?: number;
      frequencyPenalty?: number;
      presencePenalty?: number;
      stop?: string[];
      maxRetries?: number;
      retryDelay?: number;
      maxRateLimitRetries?: number;
    } = {}
  ): Promise<string> {
    const requestBody: ChatCompletionRequest = {
      messages,
      max_tokens: options.maxTokens || 1000,
      temperature: options.temperature || 0.7,
      top_p: options.topP || 0.95,
      frequency_penalty: options.frequencyPenalty || 0,
      presence_penalty: options.presencePenalty || 0,
      stop: options.stop,
    };

    const maxRetries = options.maxRetries || 3;
    const retryDelay = options.retryDelay || 1000;
    const maxRateLimitRetries = options.maxRateLimitRetries || 2;
    let lastError: Error | null = null;
    let rateLimitRetryCount = 0;

    // Implement retry logic for resilience
    for (let attempt = 0; attempt < maxRetries; attempt++) {
      try {
        // If this isn't the first attempt, add a delay before retrying
        if (attempt > 0) {
          console.log(`Retry attempt ${attempt} for Azure AI API call after ${retryDelay}ms delay...`);
          await new Promise(resolve => setTimeout(resolve, retryDelay * attempt));
        }

        const response = await fetch(
          `${this.endpoint}/openai/deployments/${this.deploymentId}/chat/completions?api-version=${this.apiVersion}`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'api-key': this.apiKey,
            },
            body: JSON.stringify(requestBody),
          }
        );

        // Handle rate limiting (429) with exponential backoff and maximum retries
        if (response.status === 429) {
          rateLimitRetryCount++;
          
          // Check if we've exceeded the maximum rate limit retries
          if (rateLimitRetryCount > maxRateLimitRetries) {
            throw new Error(`Azure AI API rate limit exceeded after ${maxRateLimitRetries} retries. Please try again later.`);
          }
          
          const retryAfter = response.headers.get('Retry-After');
          const retryMs = retryAfter ? parseInt(retryAfter, 10) * 1000 : retryDelay * Math.pow(2, attempt);
          console.log(`Rate limited by Azure AI API. Retrying after ${retryMs}ms... (${rateLimitRetryCount}/${maxRateLimitRetries})`);
          await new Promise(resolve => setTimeout(resolve, retryMs));
          continue;
        }
        
        if (!response.ok) {
          const errorText = await response.text();
          const statusCode = response.status;
          
          // Don't retry on client errors (except rate limiting which is handled above)
          if (statusCode >= 400 && statusCode < 500) {
            throw new Error(`Azure AI API error: ${statusCode} ${errorText}`);
          }
          
          // For server errors, continue with retry logic
          lastError = new Error(`Azure AI API error: ${statusCode} ${errorText}`);
          continue;
        }

        const data = (await response.json()) as ChatCompletionResponse;
        return data.choices[0].message.content;
      } catch (error: any) {
        console.error('Error getting completion from Azure AI:', error);
        
        // If it's a rate limiting error, propagate it immediately
        if (error.message && error.message.includes('rate limit exceeded')) {
          throw error;
        }
        
        // For other errors on the last retry attempt, throw the error
        if (attempt === maxRetries - 1) {
          throw error;
        }
        
        // Otherwise, store the error and continue retrying
        lastError = error;
      }
    }
    
    // If all retries failed, throw the last error
    if (lastError) {
      throw lastError;
    }
    
    // This should never happen but TypeScript needs it
    throw new Error('Failed to get completion from Azure AI after retries');
  }

  /**
   * Load configuration from environment variables
   * @returns AzureAIConfig object with API key and endpoint
   */
  static loadConfig(): AzureAIConfig {
    return {
      apiKey: process.env.AZURE_OPENAI_API_KEY || '',
      endpoint: process.env.AZURE_OPENAI_ENDPOINT || '',
      deploymentId: process.env.AZURE_OPENAI_DEPLOYMENT_ID || '',
    };
  }
}

// Export types for use in other modules
export type { AzureAIConfig, ChatMessage, ChatCompletionResponse };
