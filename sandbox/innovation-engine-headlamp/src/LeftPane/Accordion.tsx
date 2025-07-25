import { PrimaryButton } from '@fluentui/react/lib/Button';
import React, { useState } from 'react';

const Accordion = ({ items }) => {
    const [activeIndex, setActiveIndex] = useState(null);

    const toggleItem = (index) => {
        setActiveIndex(prevIndex => (prevIndex === index ? null : index));
    };

    const styles = {
        accordion: {
            width: '100%',
            margin: '0 auto',
            fontFamily: 'Arial, sans-serif',
        },
        item: {
            borderBottom: '1px solid #ccc',
            marginBottom: '10px'
        },
        title: {
            display: 'flex',
            justifyContent: 'space-between',
            cursor: 'pointer',
            padding: '1em',
            backgroundColor: '#f0f0f0',
        },
        content: {
            backgroundColor: 'rgb(237, 235, 233)',
            display: 'flex',
            flexDirection: 'row',
            gap: "15px",
            width: "95%",
        },
        code: {
            wordWrap: "break-word",
            whiteSpace: "pre-wrap",
            fontFamily: "Consolas, monospace",
            padding: "20px"
        }
    };

    return (
        <div style={styles.accordion}>
            {items.map((item, index) => (
                <div key={index} style={styles.item}>
                    <div style={styles.title} onClick={() => toggleItem(index)}>
                        <h3 style={{ margin: 0 }}>{item.name?.replace(/^\d+\.\s*/, '')}</h3>
                        <span>{activeIndex === index ? '-' : '+'}</span>
                    </div>
                    {activeIndex === index && item.codeblocks.map((codeblock) => {
                        return (
                            <>
                                <h4 style={{ width: "95%" }}>{codeblock.description}</h4>
                                <div key={codeblock.description} style={styles.content}>
                                    <div style={styles.code}>{codeblock.command}</div>
                                    <div style={{ flex: 1, margin: "0 auto" }}>
                                        <PrimaryButton
                                            text="Run"
                                            style={{ maxWidth: "20px", maxHeight: "30px", width: "20px", height: "30px" }}
                                            onClick={() => {
                                                // Handle button click
                                            }}
                                        />
                                    </div>
                                </div>
                            </>
                        );
                    })}
                </div>
            ))}
        </div>
    );
};

export default Accordion;
