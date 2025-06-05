import React from 'react';
import Typography from '@mui/material/Typography';
import { ExecDoc, ExecDocStep } from './ExecDocTypes';
import { ExecDocStepEditor } from './ExecDocStepEditor';

interface FileOperationsProps {
  execDoc: ExecDoc | null;
  onSave: (execDoc: ExecDoc) => void;
  onLoad: () => void;
  onExport: (format: 'markdown') => void;
  autoSaveEnabled: boolean;
  onToggleAutoSave: () => void;
  autoSaveInterval: number;
  onChangeAutoSaveInterval: (interval: number) => void;
  recentFiles: string[];
  onOpenRecentFile: (filePath: string) => void;
}

export const FileOperations: React.FC<FileOperationsProps> = ({
  execDoc,
  onSave,
  onLoad,
  onExport,
  autoSaveEnabled,
  onToggleAutoSave,
  autoSaveInterval,
  onChangeAutoSaveInterval,
  recentFiles,
  onOpenRecentFile
}) => {
  const handleSave = () => {
    if (execDoc) {
      onSave(execDoc);
    }
  };

  const handleAutoSaveIntervalChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    onChangeAutoSaveInterval(parseInt(e.target.value, 10));
  };

  return (
    <div style={{ borderBottom: '1px solid #e0e0e0', paddingBottom: '16px', marginBottom: '16px' }}>
      <Typography variant="h6" style={{ marginBottom: '12px' }}>File Operations</Typography>
      
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px' }}>
        <button
          onClick={handleSave}
          disabled={!execDoc}
          style={{
            padding: '8px 16px',
            backgroundColor: '#1976d2',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: !execDoc ? 'not-allowed' : 'pointer',
            opacity: !execDoc ? 0.7 : 1
          }}
        >
          Save to Disk
        </button>
        
        <button
          onClick={onLoad}
          style={{
            padding: '8px 16px',
            backgroundColor: '#f0f0f0',
            border: '1px solid #ddd',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          Load from Disk
        </button>
        
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <label style={{ marginRight: '8px', display: 'flex', alignItems: 'center' }}>
            <input
              type="checkbox"
              checked={autoSaveEnabled}
              onChange={onToggleAutoSave}
              style={{ marginRight: '4px' }}
            />
            Auto-save
          </label>
          
          {autoSaveEnabled && (
            <select
              value={autoSaveInterval}
              onChange={handleAutoSaveIntervalChange}
              style={{
                padding: '6px',
                borderRadius: '4px',
                border: '1px solid #ddd'
              }}
            >
              <option value={30}>Every 30 seconds</option>
              <option value={60}>Every minute</option>
              <option value={300}>Every 5 minutes</option>
              <option value={600}>Every 10 minutes</option>
            </select>
          )}
        </div>
      </div>
      
      {/* Save as Markdown */}
      <div style={{ marginTop: '12px' }}>
        <Typography variant="subtitle2" style={{ marginBottom: '8px' }}>Save as:</Typography>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={() => onExport('markdown')}
            disabled={!execDoc}
            style={{
              padding: '6px 12px',
              backgroundColor: '#f0f0f0',
              border: '1px solid #ddd',
              borderRadius: '4px',
              cursor: !execDoc ? 'not-allowed' : 'pointer',
              opacity: !execDoc ? 0.7 : 1
            }}
          >
            Markdown
          </button>
        </div>
      </div>
      
      {/* Recent Files */}
      {recentFiles.length > 0 && (
        <div style={{ marginTop: '16px' }}>
          <Typography variant="subtitle2" style={{ marginBottom: '8px' }}>Recent Files:</Typography>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
            {recentFiles.map((file, index) => (
              <button
                key={index}
                onClick={() => onOpenRecentFile(file)}
                style={{
                  padding: '6px 12px',
                  backgroundColor: '#f5f5f5',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  textOverflow: 'ellipsis',
                  overflow: 'hidden',
                  maxWidth: '200px',
                  whiteSpace: 'nowrap'
                }}
              >
                {file.split('/').pop()}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
