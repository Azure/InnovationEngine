import Typography from '@mui/material/Typography';
import React from 'react';
import { Message } from './ExecDocTypes';

interface OverviewAuthoringProps {
  initialOverview?: string;
  onSaveOverview: (overview: string) => void;
  onGenerateSteps: () => void;
  authoringPhase?: 'create-overview' | 'implement-content' | 'refine-content';
}

export const OverviewAuthoring: React.FC<OverviewAuthoringProps> = ({
  initialOverview = '',
  onSaveOverview,
  onGenerateSteps,
  authoringPhase = 'create-overview'
}) => {
  const [promptInput, setPromptInput] = React.useState('');
  const [generatedOverview, setGeneratedOverview] = React.useState(initialOverview);
  const [isEditingOverview, setIsEditingOverview] = React.useState(false);
  const [messages, setMessages] = React.useState<Message[]>([]);
  const [isGenerating, setIsGenerating] = React.useState(false);
  const [activePanel, setActivePanel] = React.useState<'prompt' | 'preview'>('prompt');
  
  // Handle sending prompt to Copilot
  const handleSendPrompt = () => {
    if (!promptInput.trim()) return;
    
    // Add user message to conversation
    setMessages(prev => [...prev, { role: 'user', content: promptInput }]);
    
    // Set generating state
    setIsGenerating(true);
    
    // Simulate Copilot response (this would be replaced with actual API call)
    setTimeout(() => {
      // Generate a simple overview based on the prompt
      const newOverview = `
# ${promptInput}

## Overview

This is an Executable Document that will guide you through the process of ${promptInput.toLowerCase()}. 
Follow the steps below to complete this task in your Kubernetes environment.

## Prerequisites

- Kubernetes cluster
- kubectl configured to access your cluster
- Necessary permissions to deploy resources

## Expected Outcome

Successfully ${promptInput.toLowerCase()} in your Kubernetes cluster.
      `.trim();
      
      // Update the overview
      setGeneratedOverview(newOverview);
      
      // Add assistant message
      setMessages(prev => [...prev, { 
        role: 'assistant', 
        content: `I've generated an overview for "${promptInput}". You can edit it directly or ask me to make changes.`
      }]);
      
      // Reset states
      setPromptInput('');
      setIsGenerating(false);
      
      // Switch to preview panel after generating
      setActivePanel('preview');
    }, 1000);
  };
  
  // Handle proceeding to step generation
  const handleProceedToSteps = () => {
    onSaveOverview(generatedOverview);
    onGenerateSteps();
  };
  
  // Get phase-appropriate button text
  const getActionButtonText = () => {
    switch (authoringPhase) {
      case 'create-overview':
        return 'Generate Steps';
      default:
        return 'Generate Steps';
    }
  };
  
  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Typography variant="h5" style={{ marginBottom: '16px' }}>
        Create & Edit Document Overview
      </Typography>
      
      {/* Phase guidance */}
      <div style={{ 
        padding: '8px 12px', 
        backgroundColor: '#f0f9ff', 
        borderLeft: '4px solid #1976d2',
        marginBottom: '16px'
      }}>
        <Typography variant="body2">
          {authoringPhase === 'create-overview' 
            ? 'Step 1: Create an overview that describes what this document will accomplish.' 
            : 'Step 2: Refine your overview to ensure it accurately describes the intended workflow.'}
        </Typography>
      </div>
      
      {/* Panel Toggle Buttons */}
      <div style={{ display: 'flex', marginBottom: '16px' }}>
        <button
          onClick={() => setActivePanel('prompt')}
          style={{
            flex: 1,
            padding: '8px',
            backgroundColor: activePanel === 'prompt' ? '#1976d2' : '#f1f1f1',
            color: activePanel === 'prompt' ? 'white' : 'black',
            border: 'none',
            borderRadius: '4px 0 0 4px',
            cursor: 'pointer',
            fontWeight: activePanel === 'prompt' ? 'bold' : 'normal'
          }}
        >
          Prompt & Conversation
        </button>
        <button
          onClick={() => setActivePanel('preview')}
          style={{
            flex: 1,
            padding: '8px',
            backgroundColor: activePanel === 'preview' ? '#1976d2' : '#f1f1f1',
            color: activePanel === 'preview' ? 'white' : 'black',
            border: 'none',
            borderRadius: '0 4px 4px 0',
            cursor: 'pointer',
            fontWeight: activePanel === 'preview' ? 'bold' : 'normal'
          }}
        >
          Preview & Edit
        </button>
      </div>
      
      {/* Main Content Area - Conditionally render based on active panel */}
      <div style={{ flex: 1, display: 'flex', overflow: 'hidden' }}>
        {/* Prompt Panel */}
        <div 
          style={{ 
            flex: activePanel === 'prompt' ? 1 : 0,
            display: activePanel === 'prompt' ? 'flex' : 'none',
            flexDirection: 'column',
            overflow: 'hidden'
          }}
        >
          <div style={{ marginBottom: '16px' }}>
            <Typography variant="subtitle1">
              What kind of Executable Document do you want to create?
            </Typography>
            <div style={{ display: 'flex', marginTop: '8px', gap: '8px' }}>
              <input
                type="text"
                value={promptInput}
                onChange={(e) => setPromptInput(e.target.value)}
                onKeyDown={(e) => {
                  // Handle CTRL+ENTER to submit prompt
                  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey) && promptInput.trim()) {
                    handleSendPrompt();
                  }
                }}
                placeholder="E.g., Create a deployment for a Node.js application"
                style={{
                  flex: 1,
                  padding: '10px',
                  borderRadius: '4px',
                  border: '1px solid #ccc'
                }}
              />
              <button
                onClick={handleSendPrompt}
                disabled={isGenerating || !promptInput.trim()}
                style={{
                  padding: '10px 16px',
                  backgroundColor: '#1976d2',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: isGenerating || !promptInput.trim() ? 'not-allowed' : 'pointer',
                  opacity: isGenerating || !promptInput.trim() ? 0.7 : 1
                }}
              >
                {isGenerating ? 'Generating...' : 'Generate'}
              </button>
            </div>
          </div>
          
          {/* Conversation History */}
          <div 
            style={{ 
              flex: 1,
              overflowY: 'auto',
              border: '1px solid #e0e0e0',
              borderRadius: '4px',
              padding: '12px',
              backgroundColor: '#f9f9f9'
            }}
          >
            {messages.length === 0 && (
              <Typography color="textSecondary" align="center" style={{ marginTop: '20px' }}>
                Enter a prompt above to start generating your Executable Document
              </Typography>
            )}
            
            {messages.map((message, index) => (
              <div 
                key={index} 
                style={{
                  marginBottom: '10px',
                  textAlign: message.role === 'user' ? 'right' : 'left',
                }}
              >
                <div 
                  style={{
                    display: 'inline-block',
                    maxWidth: '80%',
                    padding: '10px',
                    borderRadius: '8px',
                    backgroundColor: message.role === 'user' ? '#1976d2' : '#ffffff',
                    color: message.role === 'user' ? 'white' : 'black',
                    boxShadow: '0 1px 2px rgba(0,0,0,0.1)',
                    border: message.role === 'assistant' ? '1px solid #e0e0e0' : 'none'
                  }}
                >
                  <Typography>{message.content}</Typography>
                </div>
              </div>
            ))}
            
            {isGenerating && (
              <div style={{ textAlign: 'center', padding: '10px' }}>
                <Typography color="textSecondary">Generating overview...</Typography>
              </div>
            )}
          </div>
        </div>
        
        {/* Preview Panel */}
        <div 
          style={{ 
            flex: activePanel === 'preview' ? 1 : 0,
            display: activePanel === 'preview' ? 'flex' : 'none',
            flexDirection: 'column',
            overflow: 'hidden'
          }}
        >
          <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="subtitle1">
              {isEditingOverview ? 'Edit Overview' : 'Preview Overview'}
            </Typography>
            <div>
              <button
                onClick={() => setIsEditingOverview(!isEditingOverview)}
                style={{
                  marginRight: '8px',
                  padding: '6px 12px',
                  backgroundColor: '#f0f0f0',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  cursor: 'pointer'
                }}
              >
                {isEditingOverview ? 'Preview' : 'Edit'}
              </button>
              
              <button
                onClick={handleProceedToSteps}
                disabled={!generatedOverview.trim()}
                style={{
                  padding: '6px 12px',
                  backgroundColor: '#4caf50',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: !generatedOverview.trim() ? 'not-allowed' : 'pointer',
                  opacity: !generatedOverview.trim() ? 0.7 : 1
                }}
              >
                {getActionButtonText()}
              </button>
            </div>
          </div>
          
          <div 
            style={{ 
              flex: 1,
              border: '1px solid #e0e0e0',
              borderRadius: '4px',
              padding: '12px',
              backgroundColor: '#ffffff',
              overflowY: 'auto'
            }}
          >
            {isEditingOverview ? (
              <textarea
                value={generatedOverview}
                onChange={(e) => setGeneratedOverview(e.target.value)}
                onKeyDown={(e) => {
                  // Handle CTRL+ENTER to save changes and exit edit mode
                  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey) && generatedOverview.trim()) {
                    e.preventDefault(); // Prevent default behavior (newline)
                    setIsEditingOverview(false); // Exit edit mode
                  }
                }}
                placeholder="Enter your document overview content. Press CTRL+ENTER to save changes."
                style={{
                  width: '100%',
                  height: '100%',
                  minHeight: '400px',
                  padding: '8px',
                  borderRadius: '4px',
                  border: '1px solid #ddd',
                  fontFamily: 'monospace',
                  fontSize: '14px',
                  lineHeight: '1.5',
                  resize: 'none'
                }}
              />
            ) : (
              <div 
                style={{ 
                  padding: '16px', 
                  maxWidth: '800px', 
                  margin: '0 auto',
                  fontFamily: 'system-ui, -apple-system, sans-serif',
                  lineHeight: 1.6
                }}
              >
                {generatedOverview.split('\n').map((line, i) => {
                  if (line.startsWith('# ')) {
                    return <Typography key={i} variant="h4" style={{ marginBottom: '16px' }}>{line.substring(2)}</Typography>;
                  } else if (line.startsWith('## ')) {
                    return <Typography key={i} variant="h5" style={{ marginTop: '24px', marginBottom: '12px' }}>{line.substring(3)}</Typography>;
                  } else if (line.startsWith('- ')) {
                    return <Typography key={i} component="li" style={{ marginLeft: '20px', marginBottom: '8px' }}>{line.substring(2)}</Typography>;
                  } else if (line === '') {
                    return <br key={i} />;
                  } else {
                    return <Typography key={i} paragraph>{line}</Typography>;
                  }
                })}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
