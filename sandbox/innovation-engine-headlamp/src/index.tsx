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
        <Typography>Enter a whitelisted shell command (e.g., ie):</Typography>
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
