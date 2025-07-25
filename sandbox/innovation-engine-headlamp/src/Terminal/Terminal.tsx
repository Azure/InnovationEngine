// Terminal.js
import 'xterm/css/xterm.css';
import { AttachAddon } from '@xterm/addon-attach';
import { useContext, useEffect, useRef } from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { ExecutableDocsContext, ExecutableDocsScenarioStatus } from '../Context/ExecutableDocsContext';
import { configureResizeSocket, configureTerminalSocket } from './TerminalHelperFunctions';

const XTermTerminal = () => {
    const terminalRef = useRef(null);
    const fitAddon = useRef(new FitAddon());
    const terminalInstance = useRef(new Terminal({
        cursorBlink: true
    }));

    const execDocsContext = useContext(ExecutableDocsContext);


    // lock the terminal when it is in an executing state
    useEffect(() => {
        if (!terminalInstance.current) return;

        terminalInstance.current.attachCustomKeyEventHandler(event => {
            // This check always need to happen first
            // This snippet disables any typing in the terminal when the deployment is in progress. This
            // is an alternative to terminal.options.disableStdin = true, which messes with scrolling.
            if (execDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.NOTSTARTED || execDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.EXECUTING) {
                return false;
            }

            if (event.ctrlKey || event.metaKey) {
                // Enable copy/paste. Returning false tells xterm not to process the key event.
                // Which means the default behavior (copying/pasting) will be performed instead.

                const terminalText = terminalInstance.current.getSelection();
                if (terminalText.toString().length > 0 && event.key === "c") {
                    return false;
                }

                if (event.key === "v") {
                    return false;
                }

                // Enable screen reader.
                if (event.altKey && event.key === "r") {
                    const screenReaderMode = window.localStorage.getItem("screenReaderMode") === "off" ? "on" : "off";
                    window.localStorage.setItem("screenReaderMode", screenReaderMode);
                    terminalInstance.current.options.screenReaderMode = (screenReaderMode === "on");
                }
                return true;
            }
        });
    }, [terminalInstance, execDocsContext.scenarioStatus]);

    useEffect(() => {

        terminalInstance.current.loadAddon(fitAddon.current);
        terminalInstance.current.open(terminalRef.current);
        fitAddon.current.fit();

        // configure main input / output socket
        const terminalSocket = configureTerminalSocket(terminalInstance, execDocsContext);
        const { resizeObserver, resizeSocket } = configureResizeSocket(terminalInstance, fitAddon);

        execDocsContext.setTerminalSocket?.(terminalSocket);
        execDocsContext.setResizeSocket?.(resizeSocket);
        const attachAddon = new AttachAddon(terminalSocket);
        terminalInstance.current.loadAddon(attachAddon);

        // Log intro message
        terminalInstance.current.writeln('Welcome to the Executable Docs Terminal!');
        terminalInstance.current.writeln('We will leverage this cli experience to orchestrate the Executable Docs deployment.');
        terminalInstance.current.writeln('');

        return () => {
            terminalInstance.current.dispose();
            terminalSocket.close();
            resizeSocket.close();
            resizeObserver.disconnect();
        };
    }, []);

    return (
        <div
            ref={terminalRef}
            style={{
                backgroundColor: 'black',
                height: '100%',
                width: '100%',
            }}
        ></div>
    );
};

export default XTermTerminal;
