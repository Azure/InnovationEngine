import { metadata } from "./Metadata";

export const getTutorialOptions = () => {
    return metadata.map((tutorial) => {
        return {
            text: tutorial.title,
            key: tutorial.key
        };
    });
}


export const getConfigurableParameters = (tutorialKey: string) => {
    const tutorial = metadata.find(t => t.key === tutorialKey);
    return tutorial?.configurations?.configurableParams || []
}