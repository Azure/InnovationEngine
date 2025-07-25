import React from 'react';
import { getTutorialOptions } from '../Metadata/MetadataHelper';

export const useTutorials = () => {
    const [tutorials, setTutorials] = React.useState<Array<any>>([]);

    React.useEffect(() => {
        setTutorials(getTutorialOptions());
        // update to retrieve tutorials from storage account once scalability effort is complete
    }, []);

    return { tutorials };
}