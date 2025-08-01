/**
 * Copyright (c) Microsoft Corporation.  All rights reserved.
 */

import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { getConfigurableParameters } from "../Metadata/MetadataHelper";

export enum ExecutableDocsScenarioStatus {
  NOTSTARTED = "NOTSTARTED",
  EXECUTING = "EXECUTING",
  SUCCESS = "SUCCESS",
  FAILED = "FAILED"
}

export type Codeblock = {
  command: string;
  description: string;
}

export type Step = {
  name: string;
  codeblocks: Array<Codeblock>;
}


export type IExecutableDocsContext = {
  terminalSocket: WebSocket | null;
  resizeSocket: WebSocket | null;
  selectedTutorial: string;
  confirmSelection: boolean;
  scenarioStatus: ExecutableDocsScenarioStatus;
  configurableParams: Map<string, string>;
  deploymentSubscriptionId?: string;
  currentStep: number;
  error: string | null;
  resourceURI: Array<string> | null;
  steps: Array<Step>;
  currentCodeblock: number;
  runAllSteps: boolean;
  setRunAllSteps?: React.Dispatch<React.SetStateAction<boolean>>;
  setCurrentCodeblock?: React.Dispatch<React.SetStateAction<number>>;
  setError?: React.Dispatch<React.SetStateAction<string | null>>;
  setCurrentStep?: React.Dispatch<React.SetStateAction<number>>;
  setResourceURI?: React.Dispatch<React.SetStateAction<Array<string> | null>>;
  setSteps?: React.Dispatch<React.SetStateAction<Array<Step>>>;
  setDeploymentSubscriptionId?: React.Dispatch<React.SetStateAction<string>>;
  setConfigurableParams?: React.Dispatch<React.SetStateAction<Map<string, string>>>;
  setScenarioStatus?: React.Dispatch<React.SetStateAction<ExecutableDocsScenarioStatus>>;
  setTerminalSocket?: React.Dispatch<React.SetStateAction<WebSocket | null>>;
  setResizeSocket?: React.Dispatch<React.SetStateAction<WebSocket | null>>;
  setSelectedTutorial?: React.Dispatch<React.SetStateAction<string>>;
  setConfirmSelection?: React.Dispatch<React.SetStateAction<boolean>>;
};

export const ExecutableDocsContext = createContext<IExecutableDocsContext | null>(null);

export const useExecutableDocsContext = (): IExecutableDocsContext => {
  const context = useContext(ExecutableDocsContext);
  if (!context) {
    throw new Error('useExecutableDocsContext must be used within an ExecutableDocsProvider');
  }
  return context;
};

// Create the provider component
export const ExecutableDocsContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [terminalSocket, setTerminalSocket] = useState<WebSocket | null>(null);
  const [resizeSocket, setResizeSocket] = useState<WebSocket | null>(null);
  const [selectedTutorial, setSelectedTutorial] = useState<string>(""); // this will be the tutorial id that we will end up executing
  const [confirmSelection, setConfirmSelection] = useState<boolean>(false); // this will be used to confirm the selection of the tutorial

  const [runAllSteps, setRunAllSteps] = useState<boolean>(false);


  // values needed for making up initial execution command
  const [configurableParams, setConfigurableParams] = useState<Map<string, string>>(new Map<string, string>());
  const [deploymentSubscriptionId, setDeploymentSubscriptionId] = useState<string>('');

  // values needed for deploying
  const [scenarioStatus, setScenarioStatus] = useState<ExecutableDocsScenarioStatus>(ExecutableDocsScenarioStatus.NOTSTARTED);
  const [currentStep, setCurrentStep] = useState<number>(0); // this tells us which step we are on in the deployment process
  const [currentCodeblock, setCurrentCodeblock] = useState<number>(0); // this tells us which codeblock we are on in the current step
  const [error, setError] = useState<string | null>(null); // this will be used to display any errors that occur during the deployment process
  const [resourceURI, setResourceURI] = useState<Array<string> | null>(null); // this will be used to display the resource URI of the deployed resource
  const [steps, setSteps] = useState<Array<Step>>([]);

  useEffect(() => {
    const configurableParams = getConfigurableParameters(selectedTutorial);
    setConfigurableParams(new Map<string, string>());
    for (const param of configurableParams) {
      setConfigurableParams(prev => new Map(prev).set(param.name, param.defaultValue));
    }
  }, [selectedTutorial])

  const context = useMemo(() => ({
    resizeSocket,
    terminalSocket,
    selectedTutorial,
    confirmSelection,
    scenarioStatus,
    configurableParams,
    deploymentSubscriptionId,
    currentStep,
    error,
    resourceURI,
    steps,
    currentCodeblock,
    runAllSteps,
    setRunAllSteps,
    setCurrentCodeblock,
    setError,
    setResourceURI,
    setSteps,
    setCurrentStep,
    setConfigurableParams,
    setScenarioStatus,
    setConfirmSelection,
    setSelectedTutorial,
    setTerminalSocket,
    setResizeSocket,
    setDeploymentSubscriptionId
  }), [runAllSteps, currentCodeblock, currentStep, error, resourceURI, steps, deploymentSubscriptionId, resizeSocket, terminalSocket, scenarioStatus, selectedTutorial, confirmSelection, configurableParams, setConfigurableParams, setScenarioStatus, setConfirmSelection, setSelectedTutorial, setTerminalSocket, setResizeSocket]);


  return (
    <ExecutableDocsContext.Provider value={context}>
      {children}
    </ExecutableDocsContext.Provider>
  );
};


