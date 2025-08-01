// Terminal.js

import {
    Dropdown,
    Spinner,
    SpinnerSize,
    TextField
} from '@fluentui/react';
import { PrimaryButton } from '@fluentui/react/lib/Button';
import { useState } from 'react';
import { ExecutableDocsScenarioStatus, IExecutableDocsContext, useExecutableDocsContext } from "../Context/ExecutableDocsContext";
import { useTutorials } from '../Hooks/useTutorials';
import { metadata } from '../Metadata/Metadata';
import { getConfigurableParameters } from '../Metadata/MetadataHelper';
import Accordion from './Accordion';

const LeftPane = () => {
    const executableDocsContext: IExecutableDocsContext = useExecutableDocsContext();
    const { tutorials } = useTutorials();

    // Success Accordion component (without run buttons)
    const SuccessAccordion = ({ items }) => {
        const [activeIndex, setActiveIndex] = useState(null);

        const toggleItem = (index) => {
            setActiveIndex(prevIndex => (prevIndex === index ? null : index));
        };

        const styles = {
            accordion: {
                width: '100%',
                margin: '0 auto',
                fontFamily: 'Arial, sans-serif',
                maxHeight: '70vh',
                overflowY: 'auto' as const,
            },
            item: {
                borderTop: "1px solid rgba(204, 204, 204, 0.8)",
                borderBottom: "1px solid rgba(204, 204, 204, 0.8)",
                marginBottom: '15px',
            },
            title: {
                display: 'flex',
                justifyContent: 'space-between',
                cursor: 'pointer',
                padding: '1em',
            },
            content: {
                display: 'flex',
                flexDirection: 'column' as const,
                gap: "15px",
                width: "95%",
                marginLeft: "44px",
                marginBottom: "10px"
            },
            code: {
                wordWrap: "break-word" as const,
                whiteSpace: "pre-wrap" as const,
                fontFamily: "Consolas, monospace",
                padding: "20px",
                backgroundColor: "rgba(204, 204, 204, 0.8)",
                border: "1px solid black",
                boxSizing: "border-box" as const
            }
        };

        return (
            <div style={styles.accordion}>
                {items.map((item, index) => (
                    <div key={index} style={styles.item}>
                        <div style={styles.title} onClick={() => toggleItem(index)}>
                            <div style={{ display: 'flex', alignItems: 'center' }}>
                                <span style={{
                                    color: 'green',
                                    marginLeft: '8px',
                                    marginRight: '8px',
                                    fontSize: '15px'
                                }}>✓</span>
                                <h3 style={{ margin: 0, marginLeft: '30px' }}>{item.name?.replace(/^\d+\.\s*/, '')}</h3>
                            </div>
                            <span>{activeIndex === index ? '-' : '+'}</span>
                        </div>
                        {activeIndex === index && item.codeblocks.map((codeblock, codeblockIndex) => {
                            return (
                                <div key={codeblockIndex}>
                                    <h4 style={{ width: "95%", marginLeft: "44px" }}>{codeblock.description}</h4>
                                    <div style={styles.content}>
                                        <div style={styles.code}>{codeblock.command}</div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                ))}
            </div>
        );
    };
    if (!executableDocsContext.selectedTutorial || !executableDocsContext.confirmSelection || executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.NOTSTARTED) {
        // return a dropdown with all of the
        return (
            <>
                <title style={{ fontWeight: 'bold', fontSize: '24px', marginBottom: "10px", display: "block" }}>Executable Docs Headlamp Integration</title>
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
                <PrimaryButton
                    style={{ marginTop: '12px' }}
                    disabled={executableDocsContext.selectedTutorial === '' || executableDocsContext.deploymentSubscriptionId === ''}
                    text="Run"
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

    if (executableDocsContext.steps.length === 0 && executableDocsContext.confirmSelection && (executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.EXECUTING)) {
        return (
            <div style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                padding: '40px',
                minHeight: '200px'
            }}>
                <Spinner size={SpinnerSize.large} />
                <div style={{ marginTop: '16px', fontSize: '16px', fontWeight: '500' }}>
                    Loading
                </div>
            </div>
        );
    }

    if (executableDocsContext.confirmSelection && (executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.EXECUTING || executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.FAILED)) {
        // return the progress pane
        // add a button to run all of the steps
        const selectedMetadata = metadata.find(meta => meta.key === executableDocsContext.selectedTutorial);

        return (
            <>
                <title style={{ fontWeight: 'bold', fontSize: '24px', marginBottom: "10px", display: "block" }}>{selectedMetadata?.title || 'Tutorial'}</title>
                <p style={{ marginBottom: '12px' }}>Please follow the steps below to complete the deployment:</p>
                <PrimaryButton
                    style={{ width: "100%", height: "20px", marginBottom: '12px' }}
                    disabled={executableDocsContext.runAllSteps === true}
                    text="Run All Steps"
                    onClick={() => {
                        // send the signal to the cli to run all steps
                        executableDocsContext.terminalSocket?.send('a\n');
                        executableDocsContext.setRunAllSteps(true);
                    }}
                />
                <Accordion items={executableDocsContext.steps} />
            </>
        );
    }

    // Success Scenario
    const selectedMetadata = metadata.find(meta => meta.key === executableDocsContext.selectedTutorial);

    return (
        <div style={{
            maxHeight: '100vh',
            overflowY: 'auto',
            padding: '16px'
        }}>
            <div style={{
                display: 'flex',
                alignItems: 'center',
                marginBottom: '20px'
            }}>
                <span style={{
                    color: 'green',
                    fontSize: '24px',
                    marginRight: '12px'
                }}>✓</span>
                <h1 style={{
                    margin: 0,
                    fontSize: '24px',
                    fontWeight: 'bold',
                    color: '#333'
                }}>
                    Congratulations! Your deployment was successful!
                </h1>
            </div>
            <h2 style={{
                fontWeight: 'bold',
                fontSize: '18px',
                marginBottom: '16px',
                marginTop: '20px',
                color: '#333'
            }}>
                Resources Deployed To Azure:
            </h2>
            {executableDocsContext.resourceURI && executableDocsContext.resourceURI.length > 0 && (
                <table style={{
                    width: '100%',
                    borderCollapse: 'collapse',
                    marginBottom: '12px',
                    border: '1px solid #ddd',
                    fontFamily: 'Arial, sans-serif'
                }}>
                    <thead>
                        <tr style={{ backgroundColor: '#f5f5f5' }}>
                            <th style={{
                                padding: '12px',
                                textAlign: 'left',
                                fontWeight: 'bold',
                                border: '1px solid #ddd',
                                fontSize: '14px'
                            }}>Resource Type</th>
                            <th style={{
                                padding: '12px',
                                textAlign: 'left',
                                fontWeight: 'bold',
                                border: '1px solid #ddd',
                                fontSize: '14px'
                            }}>Resource Name</th>
                            <th style={{
                                padding: '12px',
                                textAlign: 'left',
                                fontWeight: 'bold',
                                border: '1px solid #ddd',
                                fontSize: '14px'
                            }}>Link</th>
                        </tr>
                    </thead>
                    <tbody>
                        {executableDocsContext.resourceURI.map((uri, index) => {
                            // Parse Azure resource URI to extract type and name
                            // Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/{provider}/{type}/{name}
                            const parts = uri.split('/');
                            const typeIndex = parts.findIndex(part => part === 'providers');
                            const resourceType = typeIndex !== -1 && typeIndex + 2 < parts.length ? parts[typeIndex + 2] : 'Unknown';
                            const resourceName = parts[parts.length - 1] || 'Unknown';
                            const azurePortalLink = `https://portal.azure.com/#@/resource${uri}`;
                            
                            return (
                                <tr key={index} style={{ 
                                    backgroundColor: index % 2 === 0 ? '#ffffff' : '#f9f9f9' 
                                }}>
                                    <td style={{
                                        padding: '10px 12px',
                                        border: '1px solid #ddd',
                                        fontSize: '13px'
                                    }}>{resourceType}</td>
                                    <td style={{
                                        padding: '10px 12px',
                                        border: '1px solid #ddd',
                                        fontSize: '13px'
                                    }}>{resourceName}</td>
                                    <td style={{
                                        padding: '10px 12px',
                                        border: '1px solid #ddd',
                                        fontSize: '13px'
                                    }}>
                                        <a 
                                            href={azurePortalLink} 
                                            target="_blank" 
                                            rel="noopener noreferrer"
                                            style={{
                                                color: '#0078d4',
                                                textDecoration: 'underline'
                                            }}
                                        >
                                            View in Azure Portal
                                        </a>
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            )}
            <h2 style={{
                fontWeight: 'bold',
                fontSize: '18px',
                marginBottom: '16px',
                marginTop: '20px',
                color: '#333'
            }}>
                Steps Taken For Deployment:
            </h2>
            <SuccessAccordion items={executableDocsContext.steps} />
            {selectedMetadata?.nextSteps && selectedMetadata.nextSteps.length > 0 && (
                <div style={{ marginTop: '30px' }}>
                    <h2 style={{
                        fontWeight: 'bold',
                        fontSize: '18px',
                        marginBottom: '16px',
                        color: '#333'
                    }}>
                        Next Steps:
                    </h2>
                    <ul style={{
                        listStyleType: 'disc',
                        paddingLeft: '20px',
                        margin: 0
                    }}>
                        {selectedMetadata.nextSteps.map((nextStep, index) => (
                            <li key={index} style={{ marginBottom: '8px' }}>
                                <a
                                    href={nextStep.url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    style={{
                                        color: '#0078d4',
                                        textDecoration: 'underline',
                                        fontSize: '16px'
                                    }}
                                >
                                    {nextStep.title}
                                </a>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};

export default LeftPane;
