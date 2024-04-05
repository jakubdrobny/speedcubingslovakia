import {
  CompetitionContextType,
  ResultEntry,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { Typography } from "@mui/joy";
import { reformatWithPenalties } from "../../utils";
import { useLocation } from "react-router-dom";

const Timer = () => {
  const [forceRerender, setForceRerender] = useState(false);
  const { competitionState, currentResults } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const { timerInputState } = useContext(
    TimerInputContext
  ) as TimerInputContextType;
  const solveProp: keyof ResultEntry = `solve${
    competitionState.currentSolveIdx + 1
  }` as keyof ResultEntry;
  const formattedTime = currentResults[solveProp].toString();
  const location = useLocation();
  const { handleTimerInputKeyDown, handleTimerInputKeyUp } = useContext(
    TimerInputContext
  ) as TimerInputContextType;

  useEffect(
    () => setForceRerender(!forceRerender),
    [competitionState.currentEventIdx]
  );

  useEffect(() => {
    const routePattern = /^\/competition(?:\/.*)?$/;
    if (routePattern.test(location.pathname)) {
      window.addEventListener("keydown", handleTimerInputKeyDown);
      window.addEventListener("keyup", handleTimerInputKeyUp);
    }

    return () => {
      window.removeEventListener("keydwn", handleTimerInputKeyDown);
      window.removeEventListener("keyup", handleTimerInputKeyUp);
    };
  }, [location.pathname, handleTimerInputKeyDown, handleTimerInputKeyUp]);

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        width: "100%",
      }}
    >
      <Typography level="h1" style={{ color: timerInputState.color }}>
        {timerInputState.currentState === TimerInputCurrentState.Ready
          ? "Ready"
          : timerInputState.currentState === TimerInputCurrentState.Solving
          ? "Solving..."
          : reformatWithPenalties(
              formattedTime,
              competitionState.penalties[competitionState.currentSolveIdx]
            )}
      </Typography>
    </div>
  );
};

export default Timer;
