import { CompetitionContextType, TimerInputContextType, TimerInputCurrentState } from "../../Types";
import { useEffect, useRef } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { useContext } from "react";

const Scramble = () => {
    const containerRef = useRef<HTMLDivElement>(null);
    const { competitionState } = useContext(CompetitionContext) as CompetitionContextType
    const { timerInputState } = useContext(TimerInputContext) as TimerInputContextType
    const scramble = competitionState && competitionState.currentEventIdx < competitionState.scrambles.length && competitionState.currentSolveIdx < competitionState.scrambles[competitionState.currentEventIdx].length ? competitionState.scrambles[competitionState.currentEventIdx][competitionState.currentSolveIdx] : "";
    const puzzlecode = competitionState && competitionState.currentEventIdx < competitionState.events.length ? competitionState.events[competitionState.currentEventIdx].puzzlecode : "";

    useEffect(() => {
        if (!containerRef.current) return;

        const scrambleDisplay = document.createElement('twisty-player');
        scrambleDisplay.setAttribute('alg', scramble)
        scrambleDisplay.setAttribute('hint-facelets', 'none')
        scrambleDisplay.setAttribute('background', 'none')
        scrambleDisplay.setAttribute('control-panel', 'none')
        scrambleDisplay.setAttribute('puzzle', puzzlecode)
        scrambleDisplay.setAttribute('visualization', '2D')
        containerRef.current.appendChild(scrambleDisplay);

        return () => {
            if (containerRef.current) {
                containerRef.current.removeChild(scrambleDisplay);
            }
        };
    }, [competitionState]);

    return (
        <div style={timerInputState.currentState !== TimerInputCurrentState.NotSolving ? {display: 'none'} : {display: 'flex', flexDirection: 'column', alignItems: 'center'}}>
            <h3>Scramble:</h3>
            <p style={{whiteSpace: 'pre-line'}}>{scramble}</p>
            <h3>Preview:</h3>
            <div ref={containerRef}></div>
        </div>
    );
}

export default Scramble;