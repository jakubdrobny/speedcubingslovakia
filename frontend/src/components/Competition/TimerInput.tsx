import {
  CompetitionContextType,
  ResultEntry,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import React, { useContext, useEffect, useRef } from "react";

import { CompetitionContext } from "../../context/CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { Typography } from "@mui/joy";
import { reformatWithPenalties } from "../../utils/utils";
import { useLocation } from "react-router-dom";

const Timer: React.FC<{
  handleSaveResults: (moveIndex: boolean) => void;
}> = ({ handleSaveResults }) => {
  const { competitionState, currentResultsRef, competitionStateRef } =
    useContext(CompetitionContext) as CompetitionContextType;
  const { timerInputState, timerRef } = useContext(
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

  const _handleTimerInputKeyUp = (e: Event) =>
    handleTimerInputKeyUp(e, handleSaveResults);

  const removeTimerListeners = () => {
    window.removeEventListener("keydown", handleTimerInputKeyDown);
    window.removeEventListener("keyup", _handleTimerInputKeyUp);
    if (timerRef && timerRef.current) {
      timerRef.current.removeEventListener(
        "touchstart",
        handleTimerInputKeyDown
      );
      timerRef.current.removeEventListener("touchend", _handleTimerInputKeyUp);
    }
  };

  useEffect(() => {
    const routePattern = /^\/competition(?:\/.*)?$/;
    if (routePattern.test(location.pathname)) {
      window.addEventListener("keydown", handleTimerInputKeyDown);
      window.addEventListener("keyup", _handleTimerInputKeyUp);
      if (timerRef && timerRef.current) {
        timerRef.current.addEventListener(
          "touchstart",
          handleTimerInputKeyDown
        );
        timerRef.current.addEventListener("touchend", _handleTimerInputKeyUp);
      }
    } else {
      removeTimerListeners();
    }

    return () => {
      removeTimerListeners();
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
