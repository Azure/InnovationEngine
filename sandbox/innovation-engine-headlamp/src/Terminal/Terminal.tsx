// Terminal.js
import 'xterm/css/xterm.css';
import { AttachAddon } from '@xterm/addon-attach';
import { useContext, useEffect, useRef } from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { ExecutableDocsContext } from '../Context/ExecutableDocsContext';
import { configureResizeSocket, configureTerminalSocket } from './TerminalHelperFunctions';

const XTermTerminal = () => {
    const terminalRef = useRef(null);
    const fitAddon = useRef(new FitAddon());
    const terminalInstance = useRef(new Terminal({
        cursorBlink: true
    }));

    const execDocsContext = useContext(ExecutableDocsContext);

    useEffect(() => {

        terminalInstance.current.loadAddon(fitAddon.current);
        terminalInstance.current.open(terminalRef.current);
        fitAddon.current.fit();

        // configure main input / output socket
        const terminalSocket = configureTerminalSocket(terminalInstance);
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
