/**
 * Copyright (c) Microsoft Corporation.  All rights reserved.
 */

import React, { createContext, useContext, useMemo, useState } from "react";

export enum ExecutableDocsScenarioStatus {
    NOTSTARTED = "NOTSTARTED",
    EXECUTING = "EXECUTING",
    SUCCESS = "SUCCESS",
    FAILED = "FAILED"
}


export type IExecutableDocsContext = {
    terminalSocket: WebSocket | null;
    resizeSocket: WebSocket | null;
    selectedTutorial: string;
    confirmSelection: boolean;
    scenarioStatus: ExecutableDocsScenarioStatus;
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
    const [scenarioStatus, setScenarioStatus] = useState<ExecutableDocsScenarioStatus>(ExecutableDocsScenarioStatus.NOTSTARTED);
    const [selectedTutorial, setSelectedTutorial] = useState<string>(""); // this will be the tutorial id that we will end up executing
    const [confirmSelection, setConfirmSelection] = useState<boolean>(false); // this will be used to confirm the selection of the tutorial

    const context = useMemo(() => ({
        resizeSocket,
        terminalSocket,
        selectedTutorial,
        confirmSelection,
        scenarioStatus,
        setScenarioStatus,
        setConfirmSelection,
        setSelectedTutorial,
        setTerminalSocket,
        setResizeSocket
    }), [resizeSocket, terminalSocket, scenarioStatus, selectedTutorial, confirmSelection]);


  return (
    <ExecutableDocsContext.Provider value={context}>
      {children}
    </ExecutableDocsContext.Provider>
  );
};


