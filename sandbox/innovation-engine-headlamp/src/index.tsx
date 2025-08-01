import {
  registerRoute,
  registerSidebarEntry,
} from '@kinvolk/headlamp-plugin/lib';
import { ExecutableDocsContextProvider } from './Context/ExecutableDocsContext';
import LeftPane from './LeftPane/LeftPane';
import XTermTerminal from './Terminal/Terminal';

// Add an entry to the home sidebar (not in cluster).
registerSidebarEntry({
  name: 'mypluginsidebar',
  label: 'Executable Docs',
  url: '/Executable-Docs',
  icon: 'mdi:comment-quote',
  sidebar: 'HOME',
});

registerRoute({
  path: '/Executable-Docs',
  sidebar: {
    item: 'Executable-Docs',
    sidebar: 'Innovation-engine',
  },
  useClusterURL: false,
  noAuthRequired: true, // No authentication is required to see the view
  name: 'Executable-Docs',
  exact: true,
  component: () => {
    return (
      <ExecutableDocsContextProvider>
        <div style={{ height: '100vh', display: 'flex', gap: '10px' }}>
          <div style={{ width: "calc(50% - 5px)", height: "100%" }}>
            <LeftPane />
          </div>
          <div style={{ width: "calc(50% - 5px)", height: "100%" }}>
            <XTermTerminal />
          </div>
        </div>
      </ExecutableDocsContextProvider>
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