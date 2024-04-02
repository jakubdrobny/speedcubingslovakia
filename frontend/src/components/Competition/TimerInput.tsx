import { CompetitionContextType, ResultEntry, TimerInputContextType, TimerInputCurrentState } from "../../Types";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { Typography } from "@mui/joy";
import { reformatWithPenalties } from "../../utils";

const Timer = () => {
    const [forceRerender, setForceRerender] = useState(false);
    const { competitionState } = useContext(CompetitionContext) as CompetitionContextType;
    const { timerInputState } = useContext(TimerInputContext) as TimerInputContextType;
    const solveProp: keyof ResultEntry = `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
    const formattedTime = competitionState.results[solveProp].toString();
    
    useEffect(() => setForceRerender(!forceRerender), [competitionState.currentEventIdx]);

    return (
        <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', width: '100%'}}>
            <Typography level="h1" style={{color: timerInputState.color}}>
                {timerInputState.currentState === TimerInputCurrentState.Ready
                ?
                    "Ready"
                :
                    timerInputState.currentState === TimerInputCurrentState.Solving
                    ?
                        'Solving...'
                    :
                        reformatWithPenalties(formattedTime, competitionState.penalties[competitionState.currentSolveIdx])
                }
            </Typography>
        </div>
    )
}

export default Timer;