import React from 'react';
import Typography from '@mui/material/Typography';
import { ExecDoc, ExecDocStep } from './ExecDocTypes';
import { ExecDocStepEditor } from './ExecDocStepEditor';
import { OverviewAuthoring } from './OverviewAuthoring';
import { FileOperations } from './FileOperations';
import { KubernetesContextSelector } from './KubernetesContextSelector';

interface ExecDocEditorProps {
  initialDoc?: ExecDoc | null;
}

export const ExecDocEditor: React.FC<ExecDocEditorProps> = ({
  initialDoc = null
}) => {
  // State for document
  const [execDoc, setExecDoc] = React.useState<ExecDoc | null>(initialDoc);
  const [currentView, setCurrentView] = React.useState<'overview' | 'steps'>(initialDoc ? 'steps' : 'overview');
  
  // State for tracking current phase in the authoring process
  const [authoringPhase, setAuthoringPhase] = React.useState<'create-overview' | 'refine-overview' | 'implement-content' | 'refine-content'>(
    initialDoc ? 'refine-content' : 'create-overview'
  );
  
  // Functions to progress through authoring phases
  const moveToRefineOverview = () => {
    setAuthoringPhase('refine-overview');
    setCurrentView('overview');
  };
  
  const moveToImplementContent = () => {
    setAuthoringPhase('implement-content');
    setCurrentView('steps');
  };
  
  const moveToRefineContent = () => {
    setAuthoringPhase('refine-content');
    setCurrentView('steps');
  };
  
  // Function to get the current phase name for display
  const getPhaseDisplayName = (): string => {
    switch (authoringPhase) {
      case 'create-overview': return 'Phase 1: Create Overview';
      case 'refine-overview': return 'Phase 2: Refine Overview';
      case 'implement-content': return 'Phase 3: Implement Content';
      case 'refine-content': return 'Phase 4: Refine Content';
      default: return 'Document Authoring';
    }
  };
  
  // State for file operations
  const [recentFiles, setRecentFiles] = React.useState<string[]>([]);
  const [autoSaveEnabled, setAutoSaveEnabled] = React.useState(false);
  const [autoSaveInterval, setAutoSaveInterval] = React.useState(60); // seconds
  
  // State for Kubernetes context
  const [availableContexts, setAvailableContexts] = React.useState<string[]>(['default', 'minikube', 'docker-desktop']);
  const [currentContext, setCurrentContext] = React.useState('default');
  const [availableNamespaces, setAvailableNamespaces] = React.useState<string[]>(['default', 'kube-system', 'kube-public']);
  const [currentNamespace, setCurrentNamespace] = React.useState('default');
  const [isAdmin, setIsAdmin] = React.useState(false);

  // Create a new document with overview
  const handleSaveOverview = (overview: string) => {
    if (!execDoc) {
      // Create new doc - we're in the create-overview phase
      setExecDoc({
        id: `doc-${Date.now()}`,
        title: overview.split('\n')[0].replace(/^# /, '') || 'Untitled Document',
        overview,
        steps: [],
        createdAt: new Date(),
        updatedAt: new Date(),
        kubeContext: currentContext,
        kubeNamespace: currentNamespace
      });
      // Move to refine-overview phase
      moveToRefineOverview();
    } else {
      // Update existing doc
      setExecDoc({
        ...execDoc,
        overview,
        updatedAt: new Date()
      });
    }
  };

  // Generate steps from overview (simulated)
  const handleGenerateSteps = () => {
    if (!execDoc) return;
    
    // Create sample steps based on overview (in a real implementation, this would use Copilot API)
    const newSteps: ExecDocStep[] = [
      {
        id: `step-${Date.now()}-1`,
        title: 'Setup Environment',
        description: 'Ensure you have the required tools and permissions to proceed.',
        isExpanded: true,
        isCodeBlock: false
      },
      {
        id: `step-${Date.now()}-2`,
        title: 'Create Configuration',
        description: 'Create the necessary configuration files for your deployment.',
        isExpanded: true,
        isCodeBlock: true,
        code: `apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: example-config\ndata:\n  config.json: |\n    {\n      "key": "value"\n    }`
      },
      {
        id: `step-${Date.now()}-3`,
        title: 'Deploy Application',
        description: 'Deploy the application to your Kubernetes cluster.',
        isExpanded: true,
        isCodeBlock: true,
        code: `kubectl apply -f deployment.yaml`
      },
      {
        id: `step-${Date.now()}-4`,
        title: 'Verify Deployment',
        description: 'Verify that the deployment was successful.',
        isExpanded: true,
        isCodeBlock: true,
        code: `kubectl get pods\nkubectl get services`
      }
    ];
    
    setExecDoc({
      ...execDoc,
      steps: newSteps,
      updatedAt: new Date()
    });
    
    // Move to implement-content phase
    moveToImplementContent();
  };

  // Handle step changes
  const handleStepChange = (updatedStep: ExecDocStep) => {
    if (!execDoc) return;
    
    setExecDoc({
      ...execDoc,
      steps: execDoc.steps.map(step => 
        step.id === updatedStep.id ? updatedStep : step
      ),
      updatedAt: new Date()
    });
    
    // If we're still in implement-content phase and make a change, move to refine-content
    if (authoringPhase === 'implement-content') {
      moveToRefineContent();
    }
  };

  // Handle step execution (simulated)
  const handleRunStep = (stepId: string) => {
    if (!execDoc) return;
    
    // First, mark the step as running
    setExecDoc({
      ...execDoc,
      steps: execDoc.steps.map(step => 
        step.id === stepId ? { 
          ...step, 
          executed: true,
          executionStatus: 'running',
          executionOutput: 'Running command...'
        } : step
      )
    });
    
    // Simulate execution with a delay
    setTimeout(() => {
      const step = execDoc.steps.find(s => s.id === stepId);
      if (!step) return;
      
      const isSuccess = Math.random() > 0.2; // 80% chance of success for demo
      const output = isSuccess 
        ? `Command executed successfully.\n${step.code ? `> ${step.code}\n` : ''}Output: Operation completed.`
        : `Error executing command.\n${step.code ? `> ${step.code}\n` : ''}Error: Could not complete the operation in context "${currentContext}".`;
      
      setExecDoc({
        ...execDoc,
        steps: execDoc.steps.map(s => 
          s.id === stepId ? { 
            ...s, 
            executed: true,
            executionStatus: isSuccess ? 'success' : 'failure',
            executionOutput: output
          } : s
        )
      });
    }, 1500);
  };

  // File operations (simulated)
  const handleSaveExecDoc = (doc: ExecDoc) => {
    alert(`Doc would be saved as: ${doc.title}.md`);
    
    // Add to recent files (simplified simulation)
    if (!recentFiles.includes(`/home/user/documents/${doc.title}.md`)) {
      setRecentFiles([
        `/home/user/documents/${doc.title}.md`,
        ...recentFiles.slice(0, 4) // Keep only 5 recent files
      ]);
    }
  };

  const handleLoadExecDoc = () => {
    alert('In a full implementation, a file picker would open here.');
  };

  const handleExportExecDoc = (format: 'markdown' | 'html' | 'pdf') => {
    if (!execDoc) return;
    alert(`Doc would be exported as: ${execDoc.title}.${format}`);
  };

  // Render overview authoring view
  const renderOverviewPanel = () => {
    return (
      <OverviewAuthoring
        initialOverview={execDoc?.overview || ''}
        onSaveOverview={handleSaveOverview}
        onGenerateSteps={handleGenerateSteps}
        authoringPhase={authoringPhase}
      />
    );
  };

  // Render step editing view
  const renderStepsPanel = () => {
    if (!execDoc) {
      return (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Typography>No document loaded. Create a new one first.</Typography>
        </div>
      );
    }

    return (
      <div>
        {/* Document Title */}
        <Typography variant="h4" style={{ marginBottom: '24px' }}>
          {execDoc.title}
        </Typography>
        
        {/* Overview Section */}
        <div style={{ marginBottom: '24px' }}>
          <div style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            paddingBottom: '8px',
            borderBottom: '1px solid #e0e0e0',
            marginBottom: '12px'
          }}>
            <Typography variant="h6">Overview</Typography>
            <button
              onClick={() => setCurrentView('overview')}
              style={{
                padding: '6px 12px',
                backgroundColor: '#f0f0f0',
                border: '1px solid #ddd',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Edit Overview
            </button>
          </div>
          
          {/* Preview of overview */}
          <div style={{ 
            padding: '16px', 
            backgroundColor: '#f9f9f9', 
            borderRadius: '4px',
            marginBottom: '24px',
            maxHeight: '200px',
            overflowY: 'auto'
          }}>
            {execDoc.overview.split('\n').map((line, i) => {
              if (line.startsWith('# ')) {
                return null; // Skip the title as it's shown above
              } else if (line.startsWith('## ')) {
                return <Typography key={i} variant="h6" style={{ marginTop: '16px', marginBottom: '8px' }}>{line.substring(3)}</Typography>;
              } else if (line.startsWith('- ')) {
                return <Typography key={i} component="li" style={{ marginLeft: '20px', marginBottom: '4px' }}>{line.substring(2)}</Typography>;
              } else if (line === '') {
                return <br key={i} />;
              } else {
                return <Typography key={i} paragraph>{line}</Typography>;
              }
            })}
          </div>
        </div>
        
        {/* Steps Section */}
        <Typography variant="h6" style={{ marginBottom: '16px' }}>Steps</Typography>
        {execDoc.steps.length === 0 ? (
          <Typography color="textSecondary">No steps defined yet.</Typography>
        ) : (
          execDoc.steps.map((step) => (
            <ExecDocStepEditor
              key={step.id}
              step={step}
              onStepChange={handleStepChange}
              onRunStep={handleRunStep}
              currentContext={currentContext}
              currentNamespace={currentNamespace}
              authoringPhase={authoringPhase}
            />
          ))
        )}
      </div>
    );
  };

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <div style={{ padding: '16px' }}>
        {/* Authoring Phase Indicator */}
        <div style={{ 
          marginBottom: '16px',
          padding: '8px 12px',
          backgroundColor: '#e3f2fd',
          borderRadius: '4px',
          border: '1px solid #bbdefb'
        }}>
          <Typography variant="subtitle1" style={{ fontWeight: 'bold' }}>{getPhaseDisplayName()}</Typography>
          
          <div style={{ 
            display: 'flex', 
            marginTop: '8px',
            justifyContent: 'space-between',
            flexWrap: 'wrap'
          }}>
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              margin: '4px 0'
            }}>
              <div style={{ 
                width: '16px', 
                height: '16px', 
                borderRadius: '50%',
                backgroundColor: authoringPhase === 'create-overview' ? '#1976d2' : '#ddd',
                marginRight: '8px'
              }}></div>
              <Typography variant="body2">Create Overview</Typography>
            </div>
            
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              margin: '4px 0'
            }}>
              <div style={{ 
                width: '16px', 
                height: '16px', 
                borderRadius: '50%',
                backgroundColor: authoringPhase === 'refine-overview' ? '#1976d2' : '#ddd',
                marginRight: '8px'
              }}></div>
              <Typography variant="body2">Refine Overview</Typography>
            </div>
            
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              margin: '4px 0'
            }}>
              <div style={{ 
                width: '16px', 
                height: '16px', 
                borderRadius: '50%',
                backgroundColor: authoringPhase === 'implement-content' ? '#1976d2' : '#ddd',
                marginRight: '8px'
              }}></div>
              <Typography variant="body2">Implement Content</Typography>
            </div>
            
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              margin: '4px 0'
            }}>
              <div style={{ 
                width: '16px', 
                height: '16px', 
                borderRadius: '50%',
                backgroundColor: authoringPhase === 'refine-content' ? '#1976d2' : '#ddd',
                marginRight: '8px'
              }}></div>
              <Typography variant="body2">Refine Content</Typography>
            </div>
          </div>
        </div>
      
        {/* Kubernetes Context Selector */}
        <KubernetesContextSelector
          contexts={availableContexts}
          currentContext={currentContext}
          onChangeContext={setCurrentContext}
          namespaces={availableNamespaces}
          currentNamespace={currentNamespace}
          onChangeNamespace={setCurrentNamespace}
          isAdmin={isAdmin}
        />
        
        {/* File Operations (only show in steps view) */}
        {currentView === 'steps' && (
          <FileOperations
            execDoc={execDoc}
            onSave={handleSaveExecDoc}
            onLoad={handleLoadExecDoc}
            onExport={handleExportExecDoc}
            autoSaveEnabled={autoSaveEnabled}
            onToggleAutoSave={() => setAutoSaveEnabled(prev => !prev)}
            autoSaveInterval={autoSaveInterval}
            onChangeAutoSaveInterval={setAutoSaveInterval}
            recentFiles={recentFiles}
            onOpenRecentFile={(file) => alert(`Would open: ${file}`)}
          />
        )}
        
        {/* Toggle Button for View (only show if we have a document) */}
        {execDoc && (
          <div style={{ marginBottom: '16px' }}>
            <div style={{ display: 'flex' }}>
              <button
                onClick={() => setCurrentView('overview')}
                style={{
                  flex: 1,
                  padding: '8px',
                  backgroundColor: currentView === 'overview' ? '#1976d2' : '#f1f1f1',
                  color: currentView === 'overview' ? 'white' : 'black',
                  border: 'none',
                  borderRadius: '4px 0 0 4px',
                  cursor: 'pointer',
                  fontWeight: currentView === 'overview' ? 'bold' : 'normal'
                }}
              >
                Document Overview
              </button>
              <button
                onClick={() => setCurrentView('steps')}
                style={{
                  flex: 1,
                  padding: '8px',
                  backgroundColor: currentView === 'steps' ? '#1976d2' : '#f1f1f1',
                  color: currentView === 'steps' ? 'white' : 'black',
                  border: 'none',
                  borderRadius: '0 4px 4px 0',
                  cursor: 'pointer',
                  fontWeight: currentView === 'steps' ? 'bold' : 'normal'
                }}
              >
                Document Steps
              </button>
            </div>
          </div>
        )}
      </div>
      
      {/* Main Content Area */}
      <div style={{ flex: 1, padding: '0 16px 16px', overflowY: 'auto' }}>
        {currentView === 'overview' ? renderOverviewPanel() : renderStepsPanel()}
      </div>
    </div>
  );
};
