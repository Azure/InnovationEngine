/**
 * API utilities for the frontend
 */

export const API = {
  /**
   * Base URL for server API calls
   */
  serverBaseUrl: 'http://localhost:4001', // Changed to 4001 to avoid conflict with shell-exec-backend.js
  
  /**
   * Generate an architectural overview for an Azure workload or solution
   * @param topic The workload or solution to generate an architecture overview for
   * @returns Promise with the overview content
   */
  async generateOverview(topic: string): Promise<string> {
    console.log(`Generating overview for topic: "${topic}" using server at ${this.serverBaseUrl}`);
    try {
      const response = await fetch(`${this.serverBaseUrl}/api/overview`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ topic }),
      });
      
      // First check if response is ok
    if (!response.ok) {
      try {
        // Try to parse as JSON first
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Failed to generate overview');
        } else {
          // If not JSON, get as text
          const errorText = await response.text();
          throw new Error(`Server error: ${response.status} - ${errorText.substring(0, 100)}...`);
        }
      } catch (parseError) {
        throw new Error(`Failed to generate overview: ${parseError.message}`);
      }
    }
    
    try {
      const data = await response.json();
      return data.overview || 'No overview generated';
    } catch (jsonError) {
      const textResponse = await response.clone().text();
      throw new Error(`Invalid JSON response: ${textResponse.substring(0, 100)}...`);
    }
    } catch (error) {
      console.error('Error in generateOverview:', error);
      throw error;
    }
  },
  
  /**
   * Send a query to the Azure AI Assistant
   * @param messages Array of message objects with role and content
   * @returns Promise with the assistant's response
   */
  async sendAssistantQuery(messages: Array<{role: string, content: string}>): Promise<string> {
    const response = await fetch(`${this.serverBaseUrl}/api/assistant`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ messages }),
    });
    
    if (!response.ok) {
      try {
        // Try to parse as JSON first
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Failed to get response from assistant');
        } else {
          // If not JSON, get as text
          const errorText = await response.text();
          throw new Error(`Server error: ${response.status} - ${errorText.substring(0, 100)}...`);
        }
      } catch (parseError) {
        throw new Error(`Failed to get response from assistant: ${parseError.message}`);
      }
    }
    
    const data = await response.json();
    return data.response || 'No response generated';
  },
  
  /**
   * Execute a command via the Innovation Engine CLI
   * @param command The command to execute
   * @returns Promise with the execution result
   */
  async executeCommand(command: string): Promise<{stdout: string, stderr: string, exitCode: number}> {
    const response = await fetch(`${this.serverBaseUrl}/api/exec`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ command }),
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to execute command');
    }
    
    return await response.json();
  }
};
