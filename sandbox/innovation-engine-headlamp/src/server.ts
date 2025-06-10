// Backend server to handle API requests, including Azure AI integration
import express from 'express';
import bodyParser from 'body-parser';
import cors from 'cors';
import dotenv from 'dotenv';
import { AzureAIService, type ChatMessage } from './services/azureAI';

// Load environment variables from .env file
dotenv.config();

// Create Express app
const app = express();
const port = process.env.PORT || 4001; // Changed to 4001 to avoid conflict with shell-exec-backend.js

// Middleware
app.use(cors());
app.use(bodyParser.json());

// Create Azure AI service
const azureAIService = new AzureAIService({
  apiKey: process.env.AZURE_OPENAI_API_KEY || '',
  endpoint: process.env.AZURE_OPENAI_ENDPOINT || '',
  deploymentId: process.env.AZURE_OPENAI_DEPLOYMENT_ID || '',
});

// Exec command endpoint (existing functionality)
app.post('/api/exec', async (req, res) => {
  try {
    const { command } = req.body;
    
    // This is a placeholder for the existing exec functionality
    // We're keeping this to maintain compatibility with the existing code
    
    res.json({
      stdout: `Executed command: ${command}`,
      stderr: '',
      exitCode: 0,
    });
  } catch (error: any) {
    console.error('Error executing command:', error);
    res.status(500).json({ error: error.message || 'An unknown error occurred' });
  }
});

// Assistant endpoint - processes messages through Azure AI
app.post('/api/assistant', async (req, res) => {
  try {
    const { messages } = req.body;
    
    if (!messages || !Array.isArray(messages)) {
      return res.status(400).json({ error: 'Invalid request: missing or invalid messages array' });
    }
    
    // Check for environment variables based on NODE_ENV
    const isDevelopment = process.env.NODE_ENV === 'development' || process.env.NODE_ENV === 'test';
    
    if (!process.env.AZURE_OPENAI_API_KEY || 
        !process.env.AZURE_OPENAI_ENDPOINT || 
        !process.env.AZURE_OPENAI_DEPLOYMENT_ID) {
      
      // In production mode, return an error response
      if (!isDevelopment) {
        console.error('CRITICAL: Azure OpenAI credentials not configured in production mode.');
        return res.status(500).json({ 
          error: `Azure OpenAI is not properly configured. 
                  This application requires Azure OpenAI credentials in production mode.
                  Please set the environment variables and restart the application.`
        });
      }
      
      // In development mode, use a fallback response
      console.warn('Azure OpenAI credentials not configured in development mode. Using fallback response.');
      return res.json({ 
        response: `This is a fallback response because the Azure OpenAI API is not configured. 
                   Please set the AZURE_OPENAI_API_KEY, AZURE_OPENAI_ENDPOINT, and 
                   AZURE_OPENAI_DEPLOYMENT_ID environment variables.` 
      });
    }
    
    // Transform messages to the format expected by Azure OpenAI API if needed
    const formattedMessages = messages.map(msg => ({
      role: msg.role,
      content: msg.content
    })) as ChatMessage[];
    
    // Add a system message if none exists
    if (!formattedMessages.some(msg => msg.role === 'system')) {
      formattedMessages.unshift({
        role: 'system',
        content: 'You are the Innovation Engine Assistant, a helpful AI assistant focused on Kubernetes and Executable Documents. Provide clear, concise responses to user queries.'
      });
    }
    
    // Get completion from Azure AI
    const response = await azureAIService.getCompletion(formattedMessages);
    
    res.json({ response });
  } catch (error: any) {
    console.error('Error processing assistant request:', error);
    res.status(500).json({ error: error.message || 'An unknown error occurred' });
  }
});

// Overview generation endpoint
app.post('/api/overview', async (req, res) => {
  try {
    const { topic } = req.body;
    
    if (!topic || typeof topic !== 'string') {
      return res.status(400).json({ error: 'Invalid request: missing or invalid topic' });
    }
    
    // Check for environment variables based on NODE_ENV
    const isDevelopment = process.env.NODE_ENV === 'development' || process.env.NODE_ENV === 'test';
    
    if (!process.env.AZURE_OPENAI_API_KEY || 
        !process.env.AZURE_OPENAI_ENDPOINT || 
        !process.env.AZURE_OPENAI_DEPLOYMENT_ID) {
      
      // In production mode, return an error response
      if (!isDevelopment) {
        console.error('CRITICAL: Azure OpenAI credentials not configured in production mode.');
        return res.status(500).json({ 
          error: `Azure OpenAI is not properly configured. 
                  This application requires Azure OpenAI credentials in production mode.
                  Please set the environment variables and restart the application.`
        });
      }
      
      // In development mode, use a fallback response
      console.warn('Azure OpenAI credentials not configured in development mode. Using fallback response for overview.');
      return res.json({ 
        overview: `This is a fallback overview of "${topic}" because the Azure OpenAI API is not configured. 
                  Please set the AZURE_OPENAI_API_KEY, AZURE_OPENAI_ENDPOINT, and 
                  AZURE_OPENAI_DEPLOYMENT_ID environment variables.` 
      });
    }
    
    // Get overview from Azure AI
    const promptFile = 'overview.txt';
    
    // Generate the overview with the selected prompt
    const overview = await azureAIService.generateOverview(topic, {
      systemPromptFile: promptFile
    });
    
    res.json({ overview });
  } catch (error: any) {
    console.error('Error generating overview:', error);
    res.status(500).json({ error: error.message || 'An unknown error occurred' });
  }
});

// Health check endpoint
app.get('/api/health', (req, res) => {
  console.log('Health check endpoint called');
  res.status(200).json({ status: 'ok' });
});

// Simple test endpoint for debugging
app.get('/test', (req, res) => {
  console.log('Test endpoint called');
  res.status(200).send('Server is working!');
});

// Start server
app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});

export default app;
