import React from 'react';
import {
  registerRoute,
  registerSidebarEntry,
} from '@kinvolk/headlamp-plugin/lib';
import { SectionBox } from '@kinvolk/headlamp-plugin/lib/CommonComponents';
import Typography from '@mui/material/Typography';
import { ExecDocEditor } from './components/ExecDocEditor';

// Adding a route for the Exec Doc Editor page
registerRoute({
  path: '/exec-doc-editor',
  sidebar: {
    item: 'exec-doc-editor',
    sidebar: 'Innovation-engine',
  },
  useClusterURL: false,
  noAuthRequired: true,
  name: 'exec-doc-editor',
  exact: true,
  component: () => {
    return (
      <SectionBox title="Executable Document Editor" textAlign="left" paddingTop={2}>
        <ExecDocEditor />
      </SectionBox>
    );
  },
});

// Adding the Exec Doc Editor sidebar entry
registerSidebarEntry({
  name: 'exec-doc-editor',
  label: 'Exec Doc Editor',
  url: '/exec-doc-editor',
  icon: 'mdi:file-document-edit',
  sidebar: 'Innovation-engine',
});

// Adding the Architect sidebar entry that also links to Exec Doc Editor
registerSidebarEntry({
  name: 'architect',
  label: 'Architect',
  url: '/exec-doc-editor',
  icon: 'mdi:code-braces-box',
  sidebar: 'Innovation-engine',
});
