import { useEffect, useRef } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { CompetitionContextType } from "../../Types";
import { useContext } from "react";

const Scramble = () => {
    const containerRef = useRef<HTMLDivElement>(null);
    const { competitionState } = useContext(CompetitionContext) as CompetitionContextType
    const scramble = competitionState && competitionState.currentEventIdx < competitionState.scrambles.length && competitionState.currentSolveIdx < competitionState.scrambles[competitionState.currentEventIdx].length ? competitionState.scrambles[competitionState.currentEventIdx][competitionState.currentSolveIdx].replaceAll("\n", "") : "";
    console.log(scramble);
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
        <div>
            <h3>Scramble:</h3>
            <p>{scramble}</p>
            <h3>Preview:</h3>
            <div ref={containerRef}></div>
        </div>
    );
}

export default Scramble;