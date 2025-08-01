/* eslint-disable no-unused-vars */
import React from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { ExecutableDocsScenarioStatus,IExecutableDocsContext } from '../Context/ExecutableDocsContext';

let jsonData = '';

function getDelimitterCount(data: string) {
    const delimitterRegex = /ie_(us|ue)/g;
    const matches = data.match(delimitterRegex);
    const matchLength = matches ? matches.length : 0;
    return matchLength;
}

const parseStatusData = (data: string, execDocsContext: IExecutableDocsContext) => {
    try {
        const statusData = JSON.parse(data);
        console.log('Parsed status data:', statusData);

        if (statusData.status === 'Executing') execDocsContext.setScenarioStatus?.(ExecutableDocsScenarioStatus.EXECUTING);
        if (statusData.status === 'Succeeded') execDocsContext.setScenarioStatus?.(ExecutableDocsScenarioStatus.SUCCESS);
        if (statusData.status === 'Failed') execDocsContext.setScenarioStatus?.(ExecutableDocsScenarioStatus.FAILED);

        execDocsContext.setCurrentStep?.(statusData.currentStep - 1); // innovation engine steps are 1-indexed so decrement by 1
        execDocsContext.setCurrentCodeblock?.(statusData.currentCodeBlock);
        execDocsContext.setError?.(statusData.error || null);
        execDocsContext.setResourceURI?.(statusData.resourceURIs || null);
        execDocsContext.setSteps?.(statusData.steps || []);
    } catch (e) {
        console.error('Error parsing status data:', e);
    }
}

export const configureTerminalSocket = (terminalInstance: React.RefObject<Terminal>, execDocsContext: IExecutableDocsContext) => {
    const terminalSocket = new WebSocket("ws://localhost:4001/ws/term");

    terminalSocket.onopen = () => {
        // set up the environment to execute the selected tutorial
        const innovationEngineVersion = 'v0.2.3';
        const scenariosVersion = 'v1.0.14764361219';

        terminalSocket.send('rm -rf scenarios\n');
        // install scenarios
        terminalSocket.send(`curl -Lks https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scripts/install_docs_from_release.sh | /bin/bash -s -- en-us ${scenariosVersion}\n`);
        // install Innovation Engine
        terminalSocket.send('git clone https://github.com/Azure/InnovationEngine && cd InnovationEngine && make build-ie && cd ..\n');

    };

    terminalSocket.onmessage = (event) => {
        // check if we are dealing with array buffer or string
        let eventData = '';
        if (typeof event.data === "object") {
            try {
                const enc = new TextDecoder("utf-8");
                eventData = enc.decode(event.data as any);
            } catch (e) {
                // not array buffer
            }
        }
        if (typeof event.data === 'string') {
            eventData = event.data;
        }

        jsonData += eventData;

        const matchLength = getDelimitterCount(jsonData);

        if (matchLength >= 2) {
            const statusData = jsonData.split('ie_us')[1].split('ie_ue', 2)[0];
            const indexOfFirstEndDelimitter = jsonData.indexOf('ie_ue') + 5;
            // add parsing logic here
            parseStatusData(statusData, execDocsContext);
            jsonData = jsonData.slice(indexOfFirstEndDelimitter);
        }
    };
    return terminalSocket;
}

export const configureResizeSocket = (terminalInstance: React.RefObject<Terminal>, fitAddon: React.RefObject<FitAddon>) => {
    const resizeSocket = new WebSocket("ws://localhost:4001/ws/resize");

    const resizeObserver = new ResizeObserver(() => {
        fitAddon.current.fit();
        const { cols, rows } = terminalInstance.current;
        resizeSocket.send(JSON.stringify({ cols, rows }));
    });

    resizeObserver.observe(terminalInstance.current.element);
    return { resizeSocket, resizeObserver };
}