// Terminal.js

import { BaseButton, Dropdown } from '@fluentui/react';
import { ExecutableDocsScenarioStatus, IExecutableDocsContext, useExecutableDocsContext } from "../Context/ExecutableDocsContext";
import { useTutorials } from '../Hooks/useTutorials';

const LeftPane = () => {
    const executableDocsContext: IExecutableDocsContext = useExecutableDocsContext();
    const { tutorials } = useTutorials();

    if (!executableDocsContext.selectedTutorial || !executableDocsContext.confirmSelection || executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.NOTSTARTED) {
        // return a dropdown with all of the
        return (
            <>
                <p style={{ marginBottom: '12px' }}>Please select one of the following Executable Docs from the dropdown. Then, click "Run" to execute the selected doc.</p>
                <p style={{ marginBottom: '12px' }}>Then, click "Run" to execute the selected doc.</p>
                <Dropdown
                    label="Select an Executable Doc:"
                    aria-label="Select an Executable Doc"
                    style={{ width: '60%', marginBottom: '12px' }}
                    options={tutorials}
                    selectedKey={executableDocsContext.selectedTutorial}
                    onChange={(event, option) => {
                        executableDocsContext.setSelectedTutorial(option.key as string);
                        executableDocsContext.setScenarioStatus(ExecutableDocsScenarioStatus.EXECUTING);
                    }}
                />
                <BaseButton
                    text="Run"
                    disabled={executableDocsContext.selectedTutorial === ''}
                    onClick={() => {
                        executableDocsContext.setConfirmSelection(true);
                        // add ie execute command here
                        executableDocsContext.terminalSocket?.send(`ie execute ${executableDocsContext.selectedTutorial}`);
                    }}
                    color='primary'
                />
            </>
        );
    }

    if (executableDocsContext.confirmSelection && (executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.EXECUTING || executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.FAILED)) {
        // return the progress pane
        return <div>

        </div>
    }

    return (
        <div>
            
        </div>
    );
};

export default LeftPane;
