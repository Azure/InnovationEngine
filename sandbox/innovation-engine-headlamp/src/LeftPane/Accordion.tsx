import { PrimaryButton } from '@fluentui/react/lib/Button';
import { useState, useEffect } from 'react';
import { IExecutableDocsContext, useExecutableDocsContext, ExecutableDocsScenarioStatus } from "../Context/ExecutableDocsContext";
import { Spinner, SpinnerSize } from '@fluentui/react';

const Accordion = ({ items }) => {
    const [activeIndex, setActiveIndex] = useState(null);
    const [clickedButtons, setClickedButtons] = useState(new Set());
    const [loadingSteps, setLoadingSteps] = useState(new Set());
    const [completedSteps, setCompletedSteps] = useState(new Set());
    const executableDocsContext: IExecutableDocsContext = useExecutableDocsContext();

    // Auto-expand the current step accordion
    useEffect(() => {
        if (executableDocsContext.currentStep !== null && executableDocsContext.currentStep !== undefined) {
            setActiveIndex(executableDocsContext.currentStep);
        }
    }, [executableDocsContext.currentStep]);

    // Handle runAllSteps - set all steps to loading when runAllSteps is true
    useEffect(() => {
        if (executableDocsContext.runAllSteps) {
            const allStepIndexes = new Set(items.map((_, index) => index));
            setLoadingSteps(allStepIndexes);
            setCompletedSteps(new Set()); // Clear completed steps
        }
    }, [executableDocsContext.runAllSteps, items]);

    // Handle step progression - mark all previous steps as completed when moving to next
    useEffect(() => {
        if (executableDocsContext.currentStep >= 0) {
            // Mark all steps before current step as completed
            const completedStepIndexes = new Set();
            for (let i = 0; i < executableDocsContext.currentStep; i++) {
                completedStepIndexes.add(i);
            }
            setCompletedSteps(completedStepIndexes);
            
            // Remove all completed steps from loading
            setLoadingSteps(prev => {
                const newSet = new Set(prev);
                for (let i = 0; i < executableDocsContext.currentStep; i++) {
                    newSet.delete(i);
                }
                return newSet;
            });
        }
    }, [executableDocsContext.currentStep]);

    // Handle scenario failure - mark all loading steps as failed
    useEffect(() => {
        if (executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.FAILED) {
            setLoadingSteps(new Set()); // Clear loading steps on failure
        }
    }, [executableDocsContext.scenarioStatus]);

    const renderStatusIndicator = (stepIndex) => {
        // If scenario failed, show red X
        if (executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.FAILED &&
            (loadingSteps.has(stepIndex) || stepIndex === executableDocsContext.currentStep)) {
            return <span style={{ color: 'red', marginLeft: '8px', fontSize: '15px' }}>✗</span>;
        }

        // If step is completed, show green checkmark
        if (completedSteps.has(stepIndex)) {
            return <span style={{ color: 'green', marginLeft: '8px', fontSize: '15px' }}>✓</span>;
        }

        // If step is loading (either from runAllSteps or individual run button), show spinner
        if (loadingSteps.has(stepIndex)) {
            return <Spinner size={SpinnerSize.small} style={{ marginLeft: '8px' }} />;
        }

        return null;
    };

    const toggleItem = (index) => {
        setActiveIndex(prevIndex => (prevIndex === index ? null : index));
    };

    const styles = {
        accordion: {
            width: '100%',
            margin: '0 auto',
            fontFamily: 'Arial, sans-serif',
            maxHeight: '80vh',
            overflowY: 'auto' as const,
        },
        item: {
            borderTop: "1px solid rgba(204, 204, 204, 0.8)",   /* You can change the width, style, and color */
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
            flexDirection: 'row' as const,
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
            width: "calc(50% - 7.5px)", // 50% minus half the gap
            boxSizing: "border-box" as const
        },
        buttonContainer: {
            width: "calc(50% - 7.5px)", // 50% minus half the gap
            display: "flex",
            alignItems: "center" as const,
            justifyContent: "center" as const
        }
    };

    return (
        <div style={styles.accordion}>
            {items.map((item, index) => (
                <div key={index} style={styles.item}>
                    <div style={styles.title} onClick={() => toggleItem(index)}>
                        <div style={{ display: 'flex', alignItems: 'center' }}>
                            {renderStatusIndicator(index)}
                            <h3 style={{ margin: 0, marginLeft: '30px' }}>{item.name?.replace(/^\d+\.\s*/, '')}</h3>
                        </div>
                        <span>{activeIndex === index ? '-' : '+'}</span>
                    </div>
                    {activeIndex === index && item.codeblocks.map((codeblock, codeblockIndex) => {
                        return (
                            <>
                                <h4 style={{ width: "95%", marginLeft: "44px" }}>{codeblock.description}</h4>
                                <div key={codeblock.description} style={styles.content}>
                                    <div style={styles.code}>{codeblock.command}</div>
                                    <div style={styles.buttonContainer}>
                                        <PrimaryButton
                                            disabled={executableDocsContext.runAllSteps || executableDocsContext.currentCodeblock !== codeblockIndex || executableDocsContext.currentStep !== index || clickedButtons.has(`${index}-${codeblockIndex}`)}
                                            text="Run"
                                            style={{ width: "80%", height: "20px" }}
                                            onClick={() => {
                                                const buttonId = `${index}-${codeblockIndex}`;
                                                setClickedButtons(prev => new Set(prev).add(buttonId));

                                                // Add loading state for this step if not from runAllSteps
                                                if (!executableDocsContext.runAllSteps) {
                                                    setLoadingSteps(prev => new Set(prev).add(index));
                                                }

                                                executableDocsContext.terminalSocket?.send('e\n');
                                            }}
                                        />
                                    </div>
                                </div>
                                {executableDocsContext.scenarioStatus === ExecutableDocsScenarioStatus.FAILED && (
                                    <div style={{
                                        backgroundColor: '#FFFACD', // Light yellow background
                                        display: "block",
                                        padding: '16px',
                                        margin: '16px 0',
                                        marginLeft: '44px',
                                        border: '1px solid #DDD',
                                        borderRadius: '4px'
                                    }}>
                                        <div style={{
                                            display: 'flex',
                                            alignItems: 'center',
                                            marginBottom: '8px'
                                        }}>
                                            <span style={{
                                                color: 'red',
                                                fontSize: '18px',
                                                marginRight: '8px'
                                            }}>✗</span>
                                            <h4 style={{
                                                margin: 0,
                                                color: '#333',
                                                fontSize: '16px',
                                                fontWeight: 'bold'
                                            }}>Error</h4>
                                        </div>
                                        <div style={{
                                            color: '#555',
                                            fontSize: '14px',
                                            lineHeight: '1.4',
                                            whiteSpace: 'pre-wrap'
                                        }}>
                                            {executableDocsContext.error || 'An error occurred during execution.'}
                                        </div>
                                    </div>
                                )}
                            </>
                        );
                    })}
                </div>
            ))}
        </div>
    );
};

export default Accordion;
