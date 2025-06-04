// Common types for Executable Document authoring
export interface Message {
  role: 'user' | 'assistant';
  content: string;
}

export interface ExecDocStep {
  id: string;
  title: string;
  description: string;
  code?: string;
  isExpanded: boolean;
  isCodeBlock: boolean;
  executed?: boolean;
  executionStatus?: 'success' | 'failure' | 'running' | null;
  executionOutput?: string;
}

export interface ExecDoc {
  id: string;
  title: string;
  overview: string;
  steps: ExecDocStep[];
  createdAt: Date;
  updatedAt: Date;
  author?: string;
  kubeContext?: string;
  kubeNamespace?: string;
}
