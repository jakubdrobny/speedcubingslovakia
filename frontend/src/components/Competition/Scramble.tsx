import {
  CompetitionContextType,
  ScrambleSet,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import { useContext, useEffect } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { Typography } from "@mui/joy";
import { useRef } from "react";

const Scramble = () => {
  const scrambleImgRef = useRef<HTMLDivElement>(null);
  const { competitionState } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const { timerInputState } = useContext(
    TimerInputContext
  ) as TimerInputContextType;
  const scramble =
    competitionState &&
    competitionState.scrambles &&
    competitionState.events &&
    competitionState.currentEventIdx < competitionState.events.length &&
    competitionState.scrambles.find(
      (s: ScrambleSet) =>
        s.event.displayname ===
        competitionState.events[competitionState.currentEventIdx].displayname
    ) != undefined
      ? (
          competitionState.scrambles.find(
            (s: ScrambleSet) =>
              s.event.displayname ===
              competitionState.events[competitionState.currentEventIdx]
                .displayname
          ) as ScrambleSet
        ).scrambles[competitionState.currentSolveIdx].scramble
      : "";

  useEffect(() => {
    if (
      competitionState &&
      competitionState.scrambles &&
      competitionState.events &&
      competitionState.currentEventIdx < competitionState.events.length &&
      competitionState.scrambles.find(
        (s: ScrambleSet) =>
          s.event.displayname ===
          competitionState.events[competitionState.currentEventIdx].displayname
      ) != undefined &&
      scrambleImgRef != null &&
      scrambleImgRef.current != null
    ) {
      const scrambleSet = competitionState.scrambles.find(
        (s: ScrambleSet) =>
          s.event.displayname ===
          competitionState.events[competitionState.currentEventIdx].displayname
      ) as ScrambleSet;
      scrambleImgRef.current.innerHTML =
        scrambleSet.scrambles[competitionState.currentSolveIdx].svgimg;
    }
  }, []);

  return (
    <div
      style={
        [
          TimerInputCurrentState.NotSolving,
          TimerInputCurrentState.GettingReady,
        ].includes(timerInputState.currentState)
          ? { display: "flex", flexDirection: "column", alignItems: "center" }
          : { display: "none" }
      }
    >
      <h3>Scramble:</h3>
      <Typography style={{ whiteSpace: "pre-line" }}>{scramble}</Typography>
      <h3>Preview:</h3>
      <div ref={scrambleImgRef}></div>
    </div>
  );
};

export default Scramble;
