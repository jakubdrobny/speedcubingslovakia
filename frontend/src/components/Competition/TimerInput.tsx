import { CompetitionContextType, ResultEntry } from "../../Types";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";

const Timer = () => {
    const [forceRerender, setForceRerender] = useState(false);
    const { competitionState, updateSolve, saveResults } = useContext(CompetitionContext) as CompetitionContextType;
    const solveProp: keyof ResultEntry= `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
    const formattedTime = competitionState.results[solveProp].toString();
    
    useEffect(() => setForceRerender(!forceRerender), [competitionState.currentEventIdx]);

    return (
        <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', width: '100%'}}>
            <h1>{formattedTime}</h1>
        </div>
    )
}

export default Timer;