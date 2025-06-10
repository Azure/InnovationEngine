// Azure AI service for the Innovation Engine Assistant
// This service handles communication with the Azure AI API using the official OpenAI SDK

import { AzureOpenAI } from 'openai';
import { loadSystemPrompt } from '../utils/promptUtils';

interface AzureAIConfig {
  apiKey: string;
  endpoint: string;
  deploymentId: string;
}

interface ChatMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

/**
 * A service class to interact with Azure OpenAI API using the official SDK
 */
export class AzureAIService {
  private client: AzureOpenAI;
  private deploymentId: string;

  /**
   * Create a new Azure AI Service client
   * @param config Configuration for the Azure AI service
   */
  constructor(config: AzureAIConfig) {
    this.deploymentId = config.deploymentId;
    
    this.client = new AzureOpenAI({
      apiKey: config.apiKey,
      endpoint: config.endpoint,
      apiVersion: '2023-12-01-preview', // Updated to latest stable API version
      dangerouslyAllowBrowser: true, // Allow in test environments
    });
  }

  /**
   * Get a completion from Azure OpenAI API
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
    } = {}
  ): Promise<string> {
    try {
      const completion = await this.client.chat.completions.create({
        model: this.deploymentId, // In Azure OpenAI, model refers to the deployment name
        messages: messages.map(msg => ({
          role: msg.role,
          content: msg.content
        })),
        max_tokens: options.maxTokens || 1000,
        temperature: options.temperature || 0.7,
        top_p: options.topP || 0.95,
        frequency_penalty: options.frequencyPenalty || 0,
        presence_penalty: options.presencePenalty || 0,
        stop: options.stop,
      }, {
        // SDK handles retries automatically, but we can configure them
        maxRetries: options.maxRetries || 3,
      });

      const content = completion.choices[0]?.message?.content;
      if (!content) {
        throw new Error('No content received from Azure OpenAI API');
      }

      return content;
    } catch (error: any) {
      console.error('Error getting completion from Azure AI:', error);
      
      // The SDK provides better error handling with specific error types
      if (error.status === 429) {
        throw new Error(`Azure AI API rate limit exceeded: ${error.message}`);
      } else if (error.status >= 400 && error.status < 500) {
        throw new Error(`Azure AI API client error (${error.status}): ${error.message}`);
      } else if (error.status >= 500) {
        throw new Error(`Azure AI API server error (${error.status}): ${error.message}`);
      }
      
      // For other errors, re-throw as-is
      throw error;
    }
  }

  /**
   * Generate an overview for a specific topic using a system prompt loaded from a file
   * @param topic The topic to generate an overview for
   * @param options Optional parameters for the completion
   * @returns A promise with the overview content
   */
  async generateOverview(
    topic: string,
    options: {
      maxTokens?: number;
      temperature?: number;
      topP?: number;
      systemPromptFile?: string;
    } = {}
  ): Promise<string> {
    // Load system prompt from file (default to overview.txt)
    const systemPromptFile = options.systemPromptFile || 'overview.txt';
    const systemPrompt = loadSystemPrompt(systemPromptFile);
    
    // Create the messages array with system and user prompts
    const messages: ChatMessage[] = [
      {
        role: 'system',
        content: systemPrompt
      },
      {
        role: 'user',
        content: `Create an overview of: ${topic}`
      }
    ];
    
    // Set completion options
    const completionOptions = {
      maxTokens: options.maxTokens || 1500,
      temperature: options.temperature || 0.5,
      topP: options.topP || 0.9,
    };
    
    // Get completion from OpenAI
    return this.getCompletion(messages, completionOptions);
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
export type { AzureAIConfig, ChatMessage };
