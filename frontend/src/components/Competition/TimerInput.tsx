import {
  CompetitionContextType,
  ResultEntry,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import { useContext, useEffect, useRef, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { Typography } from "@mui/joy";
import { reformatWithPenalties } from "../../utils";
import { useLocation } from "react-router-dom";

const Timer = () => {
  const { competitionState, currentResultsRef, competitionStateRef } =
    useContext(CompetitionContext) as CompetitionContextType;
  const { timerInputState /*,timerElementRef*/ } = useContext(
    TimerInputContext
  ) as TimerInputContextType;
  const formattedTime =
    currentResultsRef.current[
      `solve${
        competitionStateRef.current.currentSolveIdx + 1
      }` as keyof ResultEntry
    ].toString();
  const location = useLocation();
  const { handleTimerInputKeyDown, handleTimerInputKeyUp } = useContext(
    TimerInputContext
  ) as TimerInputContextType;
  const timerRef = useRef<any>(null);

  useEffect(() => {
    const routePattern = /^\/competition(?:\/.*)?$/;
    if (routePattern.test(location.pathname)) {
      window.addEventListener("keydown", handleTimerInputKeyDown);
      window.addEventListener("keyup", handleTimerInputKeyUp);
      if (timerRef && timerRef.current) {
        timerRef.current.addEventListener(
          "touchstart",
          handleTimerInputKeyDown
        );
        timerRef.current.addEventListener("touchend", handleTimerInputKeyUp);
      }
    }

    return () => {
      window.removeEventListener("keydown", handleTimerInputKeyDown);
      window.removeEventListener("keyup", handleTimerInputKeyUp);
      if (timerRef && timerRef.current) {
        timerRef.current.removeEventListener(
          "touchstart",
          handleTimerInputKeyDown
        );
        timerRef.current.removeEventListener("touchend", handleTimerInputKeyUp);
      }
    };
  }, [location.pathname, handleTimerInputKeyDown, handleTimerInputKeyUp]);

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        width: "100%",
        userSelect: "none",
        MozUserSelect: "none",
        msUserSelect: "none",
        WebkitUserSelect: "none",
      }}
      ref={timerRef}
    >
      <Typography
        level="h1"
        style={{ color: timerInputState.color }}
        // ref={timerElementRef}
      >
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
