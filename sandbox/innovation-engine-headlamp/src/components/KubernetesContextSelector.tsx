import React from 'react';
import Typography from '@mui/material/Typography';

interface KubernetesContextSelectorProps {
  contexts: string[];
  currentContext: string;
  onChangeContext: (context: string) => void;
  namespaces: string[];
  currentNamespace: string;
  onChangeNamespace: (namespace: string) => void;
  isAdmin: boolean;
}

export const KubernetesContextSelector: React.FC<KubernetesContextSelectorProps> = ({
  contexts,
  currentContext,
  onChangeContext,
  namespaces,
  currentNamespace,
  onChangeNamespace,
  isAdmin
}) => {
  const handleContextChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    onChangeContext(e.target.value);
  };

  const handleNamespaceChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    onChangeNamespace(e.target.value);
  };

  return (
    <div style={{ borderBottom: '1px solid #e0e0e0', paddingBottom: '16px', marginBottom: '16px' }}>
      <Typography variant="h6" style={{ marginBottom: '12px' }}>Kubernetes Context</Typography>
      
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px' }}>
        <div>
          <label htmlFor="context-selector" style={{ display: 'block', marginBottom: '6px' }}>
            Context:
          </label>
          <select
            id="context-selector"
            value={currentContext}
            onChange={handleContextChange}
            style={{
              padding: '8px',
              minWidth: '200px',
              borderRadius: '4px',
              border: '1px solid #ddd'
            }}
          >
            {contexts.map((context) => (
              <option key={context} value={context}>
                {context} {context === currentContext ? '(current)' : ''}
              </option>
            ))}
          </select>
        </div>
        
        <div>
          <label htmlFor="namespace-selector" style={{ display: 'block', marginBottom: '6px' }}>
            Namespace:
          </label>
          <select
            id="namespace-selector"
            value={currentNamespace}
            onChange={handleNamespaceChange}
            style={{
              padding: '8px',
              minWidth: '200px',
              borderRadius: '4px',
              border: '1px solid #ddd'
            }}
          >
            {namespaces.map((namespace) => (
              <option key={namespace} value={namespace}>
                {namespace} {namespace === currentNamespace ? '(current)' : ''}
              </option>
            ))}
          </select>
        </div>
        
        <div style={{ display: 'flex', alignItems: 'center', marginLeft: 'auto' }}>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            padding: '8px 12px',
            backgroundColor: '#f5f5f5',
            borderRadius: '4px',
            border: '1px solid #e0e0e0'
          }}>
            <span style={{ 
              width: '10px', 
              height: '10px', 
              borderRadius: '50%', 
              backgroundColor: isAdmin ? '#4caf50' : '#ff9800',
              display: 'inline-block'
            }}></span>
            <Typography variant="body2">
              {isAdmin ? 'Admin Access' : 'Standard Access'}
            </Typography>
          </div>
        </div>
      </div>
    </div>
  );
};
