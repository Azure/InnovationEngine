// Terminal.js

import {
    BaseButton,
    Dropdown,
    TextField
} from '@fluentui/react';
import { ExecutableDocsScenarioStatus, IExecutableDocsContext, useExecutableDocsContext } from "../Context/ExecutableDocsContext";
import { useTutorials } from '../Hooks/useTutorials';
import { getConfigurableParameters } from '../Metadata/MetadataHelper';
import Accordion from './Accordion';

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
                    }}
                />
                <div style={{ width: '60%', marginBottom: '12px' }}>
                    <TextField
                        label={"Subscription ID"}
                        value={executableDocsContext.deploymentSubscriptionId || ''}
                        onChange={(e) => {
                            executableDocsContext.setDeploymentSubscriptionId(e.target.value);
                        }}
                    />
                </div>
                {getConfigurableParameters(executableDocsContext.selectedTutorial).map((param) => {
                    return (
                        <div style={{ width: '60%', marginBottom: '12px' }}>
                            <TextField
                                label={param.title}
                                value={executableDocsContext.configurableParams.get(param.commandKey) || ''}
                                onChange={(e) => {
                                    executableDocsContext.setConfigurableParams(new Map(executableDocsContext.configurableParams).set(param.commandKey, e.target.value));
                                }}
                            />
                        </div>

                    );
                })}
                <BaseButton
                    text="Run"
                    disabled={executableDocsContext.selectedTutorial === '' || executableDocsContext.deploymentSubscriptionId === ''}
                    onClick={() => {
                        executableDocsContext.setConfirmSelection(true);
                        executableDocsContext.setScenarioStatus(ExecutableDocsScenarioStatus.EXECUTING);

                        let command = `./InnovationEngine/bin/ie interactive scenarios/${executableDocsContext.selectedTutorial} --subscription ${executableDocsContext.deploymentSubscriptionId} --correlation-id ${crypto.randomUUID()} `;
                        const tutorialDirectory = `${executableDocsContext.selectedTutorial.substring(0, executableDocsContext.selectedTutorial.lastIndexOf("/"))}`;
                        const configurableParamKeys = Array.from(executableDocsContext.configurableParams.keys()) || [];
                        for (const key of configurableParamKeys) {
                            command += `--var ${key}=${executableDocsContext.configurableParams.get(key as string)} `;
                        }
                        command += `--environment azure --working-directory scenarios/${tutorialDirectory} && source /tmp/env-vars\n`;


                        executableDocsContext.terminalSocket?.send(command);
                    }}
                />
            </>
        );
    }

    if (executableDocsContext.confirmSelection && (executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.EXECUTING || executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.FAILED)) {
        // return the progress pane
        // add a button to run all of the steps
        return (
            <div style={{ overflowY: 'auto', maxHeight: "500px" }}>
                <Accordion items={executableDocsContext.steps} />
            </div>

        );
    }

    return (
        <div>

        </div>
    );
};

export default LeftPane;
