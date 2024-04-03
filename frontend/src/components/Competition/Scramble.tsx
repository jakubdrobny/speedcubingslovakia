import { CompetitionContextType, ScrambleSet, TimerInputContextType, TimerInputCurrentState } from "../../Types";
import { useEffect, useRef } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { Typography } from "@mui/joy";
import { useContext } from "react";

const Scramble = () => {
    const containerRef = useRef<HTMLDivElement>(null);
    const { competitionState } = useContext(CompetitionContext) as CompetitionContextType
    const { timerInputState } = useContext(TimerInputContext) as TimerInputContextType
    const scramble = competitionState && competitionState.scrambles && competitionState.events && competitionState.currentEventIdx < competitionState.events.length && competitionState.scrambles.find((s: ScrambleSet) => s.event.displayname === competitionState.events[competitionState.currentEventIdx].displayname) != undefined  ? (competitionState.scrambles.find((s: ScrambleSet) => s.event.displayname === competitionState.events[competitionState.currentEventIdx].displayname) as ScrambleSet).scrambles[competitionState.currentSolveIdx] : "";
    const puzzlecode = competitionState && competitionState.events && competitionState.currentEventIdx < competitionState.events.length ? competitionState.events[competitionState.currentEventIdx].puzzlecode : "";

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
        <div style={[TimerInputCurrentState.NotSolving, TimerInputCurrentState.GettingReady].includes(timerInputState.currentState) ? {display: 'flex', flexDirection: 'column', alignItems: 'center'} : {display: 'none'}}>
            <h3>Scramble:</h3>
            <Typography style={{whiteSpace: 'pre-line'}}>{scramble}</Typography>
            <h3>Preview:</h3>
            <div ref={containerRef}></div>
        </div>
    );
}

export default Scramble;