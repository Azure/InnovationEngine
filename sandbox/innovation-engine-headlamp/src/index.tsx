import {
  registerRoute,
  registerRouteFilter,
  registerSidebarEntry,
  registerSidebarEntryFilter,
} from '@kinvolk/headlamp-plugin/lib';
import { SectionBox } from '@kinvolk/headlamp-plugin/lib/CommonComponents';
import Typography from '@mui/material/Typography';
import React from 'react';

// Add an entry to the home sidebar (not in cluster).
registerSidebarEntry({
  name: 'mypluginsidebar',
  label: 'Getting Started',
  url: '/getting-started',
  icon: 'mdi:comment-quote',
  sidebar: 'HOME',
});

registerRoute({
  path: '/getting-started',
  sidebar: {
    item: 'getting-started',
    sidebar: 'Innovation-engine',
  },
  useClusterURL: false,
  noAuthRequired: true, // No authentication is required to see the view
  name: 'getting-started',
  exact: true,
  component: () => {
    const navigateToShellExec = () => {
      window.location.hash = '#/shell-exec';
      setTimeout(() => {
        window.dispatchEvent(
          new CustomEvent('prefill-innovation-engine-command', {
            detail: 'ie execute ../../../scenarios/testing/variableHierarchy.md',
          })
        );
      }, 100);
    };
    return (
      <SectionBox title="Getting Started" textAlign="center" paddingTop={2}>
        <Typography>This is where Innovation Engine lives</Typography>
        <br />
        <a
          href="#"
          onClick={e => {
            e.preventDefault();
            navigateToShellExec();
          }}
          style={{
            fontWeight: 'bold',
            color: '#1976d2',
            textDecoration: 'underline',
            cursor: 'pointer',
          }}
        >
          Test Innovation Engine
        </a>
      </SectionBox>
    );
  },
});

// Adds a completely new sidebar + entry because the sidebar "innovation-engine" does not exist.
registerSidebarEntry({
  name: 'backtohome',
  label: 'Back to Home',
  url: '/',
  icon: 'mdi:hexagon',
  sidebar: 'Innovation-engine',
});

// Adds a entry to the recently created sidebar "innovation-engine".
registerSidebarEntry({
  name: 'getting-started',
  label: 'Getting Started',
  url: '/getting-started',
  icon: 'mdi:comment-quote',
  sidebar: 'Innovation-engine',
});

// Please see https://icon-sets.iconify.design/mdi/ for icons.
registerRoute({
  path: '/shell-exec',
  sidebar: {
    item: 'shell-exec',
    sidebar: 'Innovation-engine',
  },
  useClusterURL: false,
  noAuthRequired: true,
  name: 'shell-exec',
  exact: true,
  component: () => {
    const [command, setCommand] = React.useState('ie execute ');
    const [output, setOutput] = React.useState('');
    const [error, setError] = React.useState('');
    const [loading, setLoading] = React.useState(false);

    // Listen for prefill event
    React.useEffect(() => {
      const handler = e => {
        if (e.detail) setCommand(e.detail);
      };
      window.addEventListener('prefill-innovation-engine-command', handler);
      return () => window.removeEventListener('prefill-innovation-engine-command', handler);
    }, []);

    const handleExec = async () => {
      setLoading(true);
      setOutput('');
      setError('');
      try {
        const res = await fetch('http://localhost:4000/api/exec', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ command }),
        });
        const data = await res.json();
        if (res.ok) {
          setOutput(`stdout:\n${data.stdout}\nstderr:\n${data.stderr}\nexitCode: ${data.exitCode}`);
        } else {
          setError(data.error || 'Unknown error');
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    return (
      <SectionBox title="Innovation Engine" textAlign="center" paddingTop={2}>
        <Typography>Enter an allowlisted shell command (e.g., ie):</Typography>
        <input
          type="text"
          value={command}
          onChange={e => setCommand(e.target.value)}
          onKeyDown={e => {
            if (e.key === 'Enter' && !loading) {
              handleExec();
            }
          }}
          style={{ width: '60%', margin: '1em 0', padding: '0.5em' }}
        />
        <br />
        <button onClick={handleExec} disabled={loading} style={{ padding: '0.5em 1em' }}>
          {loading ? 'Running...' : 'Execute'}
        </button>
        <pre style={{ textAlign: 'left', marginTop: '1em', background: '#f5f5f5', padding: '1em' }}>
          {output}
        </pre>
        {error && <Typography color="error">Error: {error}</Typography>}
      </SectionBox>
    );
  },
});

registerSidebarEntry({
  name: 'shell-exec',
  label: 'Innovation Engine',
  url: '/shell-exec',
  icon: 'mdi:console',
  sidebar: 'Innovation-engine',
});

// Adding the Assistant sidebar entry
registerSidebarEntry({
  name: 'assistant',
  label: 'Assistant',
  url: '/assistant',
  icon: 'mdi:robot',
  sidebar: 'Innovation-engine',
});

// Adding a route for the Assistant page
registerRoute({
  path: '/assistant',
  sidebar: {
    item: 'assistant',
    sidebar: 'Innovation-engine',
  },
  useClusterURL: false,
  noAuthRequired: true,
  name: 'assistant',
  exact: true,
  component: () => {
    const [userQuery, setUserQuery] = React.useState('');
    const [chatHistory, setChatHistory] = React.useState([
      { role: 'assistant', content: 'Hello! I\'m the Innovation Engine Assistant. How can I help you with your Kubernetes or Executable Document needs?' }
    ]);
    const [isProcessing, setIsProcessing] = React.useState(false);
    const chatContainerRef = React.useRef(null);

    // Auto-scroll to bottom of chat when history changes
    React.useEffect(() => {
      if (chatContainerRef.current) {
        chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
      }
    }, [chatHistory]);

    // Handle sending a new query to the assistant
    const handleSendQuery = () => {
      if (!userQuery.trim()) return;
      
      // Add user query to chat history
      setChatHistory(prev => [...prev, { role: 'user', content: userQuery }]);
      
      // Set processing state to show loading
      setIsProcessing(true);
      
      // Simulate a response (this would be replaced with actual GitHub Copilot API call)
      setTimeout(() => {
        // Add simulated response
        setChatHistory(prev => [...prev, { 
          role: 'assistant', 
          content: `I received your request: "${userQuery}". This is a placeholder response. In the full implementation, this would be processed by GitHub Copilot.`
        }]);
        setUserQuery('');
        setIsProcessing(false);
      }, 1000);
    };

    return (
      <SectionBox title="Innovation Engine Assistant" textAlign="left" paddingTop={2}>
        {/* Chat history display area */}
        <div 
          ref={chatContainerRef}
          style={{ 
            height: '400px', 
            overflowY: 'auto', 
            marginBottom: '20px',
            padding: '10px',
            border: '1px solid #e0e0e0',
            borderRadius: '4px',
            backgroundColor: '#f9f9f9'
          }}
        >
          {chatHistory.map((message, index) => (
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
          {isProcessing && (
            <div style={{ textAlign: 'center', padding: '10px' }}>
              <Typography color="textSecondary">Processing your request...</Typography>
            </div>
          )}
        </div>
        
        {/* Input area */}
        <div style={{ display: 'flex', alignItems: 'flex-start', gap: '10px' }}>
          <textarea
            value={userQuery}
            onChange={(e) => setUserQuery(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                handleSendQuery();
              }
            }}
            placeholder="Enter your question or request here (e.g., 'Create a deployment for my app')"
            style={{
              flexGrow: 1,
              padding: '10px',
              borderRadius: '4px',
              border: '1px solid #ccc',
              minHeight: '100px',
              resize: 'vertical',
              fontFamily: 'inherit',
              fontSize: '14px'
            }}
          />
          <button
            onClick={handleSendQuery}
            disabled={isProcessing || !userQuery.trim()}
            style={{
              padding: '10px 20px',
              backgroundColor: '#1976d2',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: isProcessing || !userQuery.trim() ? 'not-allowed' : 'pointer',
              opacity: isProcessing || !userQuery.trim() ? 0.7 : 1
            }}
          >
            Send
          </button>
        </div>

        {/* Quick start section */}
        <div style={{ marginTop: '20px', borderTop: '1px solid #e0e0e0', paddingTop: '20px' }}>
          <Typography variant="h6">Quick Start:</Typography>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px', marginTop: '10px' }}>
            {['Create a simple deployment', 'Expose a service', 'Author an Executable Document', 'Learn about Kubernetes basics'].map((suggestion, index) => (
              <button
                key={index}
                onClick={() => setUserQuery(suggestion)}
                style={{
                  padding: '8px 12px',
                  backgroundColor: '#f1f1f1',
                  border: '1px solid #ddd',
                  borderRadius: '16px',
                  cursor: 'pointer'
                }}
              >
                {suggestion}
              </button>
            ))}
          </div>
        </div>
      </SectionBox>
    );
  },
});
