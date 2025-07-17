/* eslint-disable no-unused-vars */
import React from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';

let jsonData = '';

function getDelimitterCount(data: string) {
    const delimitterRegex = /ie_(us|ue)/g;
    const matches = data.match(delimitterRegex);
    const matchLength = matches ? matches.length : 0;
    return matchLength;
}

export const configureTerminalSocket = (terminalInstance: React.RefObject<Terminal>) => {
    const terminalSocket = new WebSocket("ws://localhost:4001/ws/term");

    terminalSocket.onopen = () => {
        // set up the environment to execute the selected tutorial
        const innovationEngineVersion = 'v0.2.3';
        const scenariosVersion = 'v1.0.14764361219';

        terminalSocket.send('rm -rf ExecDocs && mkdir ExecDocs && cd ExecDocs\n');
        terminalSocket.send(`wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/$VERSION/ie\n`);
        terminalSocket.send('chmod +x ie && mkdir -p ~/.local/bin && mv ie ~/.local/bin\n')

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
            jsonData = jsonData.slice(indexOfFirstEndDelimitter);
        }
    };

    terminalInstance.current.onData(data => {
        // to be implemented
    });

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