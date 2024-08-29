import {
  CompetitionContextType,
  InputMethod,
  TimerColors,
  TimerInputContextType,
  TimerInputCurrentState,
  TimerInputState,
} from "../Types";
import React, {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";

import { CompetitionContext } from "./CompetitionContext";
import { milisecondsToFormattedTime } from "../utils/utils";

export const TimerInputContext = createContext<TimerInputContextType | null>(
  null
);

export const TimerInputProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [timerInputState, setTimerInputState] =
    useState<TimerInputState>(initialState);
  const { saveResults, updateSolve, updateCurrentSolve, competitionState } =
    useContext(CompetitionContext) as CompetitionContextType;
  let holdingTimeout: {
    current: ReturnType<typeof setTimeout> | undefined | 1;
  } = useRef(undefined);
  let timingInterval: { current: ReturnType<typeof setInterval> | undefined } =
    useRef(undefined);
  const elapsedTime = useRef(0);
  //   const timerElementRef = useRef<HTMLElement | null>(null);
  //   const allNodesPositions = useRef<string[]>([]);

  //   const hideAllElementsExceptTimer = () => {
  //     allNodesPositions.current = [];
  //     (document.querySelectorAll("*") as NodeListOf<HTMLElement>).forEach(
  //       (el: HTMLElement) => {
  //         el.style.visibility = "hidden";
  //         allNodesPositions.current.push(el.style.position);
  //         el.style.position = "static";
  //       }
  //     );

  //     if (timerElementRef.current !== null) {
  //       timerElementRef.current.style.visibility = "visible";
  //       timerElementRef.current.style.position = "absolute";
  //       timerElementRef.current.style.top = "50%";
  //       timerElementRef.current.style.bottom = "50%";
  //       timerElementRef.current.style.transform = "translate(-50%, -50%)";
  //     }
  //   };
  //   const revertHidingAllElementsExceptTimer = () => {
  //     (document.querySelectorAll("*") as NodeListOf<HTMLElement>).forEach(
  //       (el: HTMLElement, idx: number) => {
  //         el.style.removeProperty("visibility");
  //         el.style.position = allNodesPositions.current[idx];
  //       }
  //     );
  //     if (timerElementRef.current !== null) {
  //       timerElementRef.current.style.position = "relative";
  //       timerElementRef.current.style.top = "";
  //       timerElementRef.current.style.bottom = "";
  //       timerElementRef.current.style.transform = "";
  //     }
  //   };

  const handleTimerInputKeyDown = useCallback(
    (e: Event) => {
      const ev = e as KeyboardEvent;

      if (ev.key === " " || ev.type === "touchstart") {
        if (
          !holdingTimeout.current &&
          !timingInterval.current &&
          timerInputState.currentState === TimerInputCurrentState.NotSolving
        ) {
          setTimerInputState((ps) => ({
            ...ps,
            currentState: TimerInputCurrentState.GettingReady,
            color: TimerColors.Red,
          }));
          holdingTimeout.current = setTimeout(() => {
            setTimerInputState((ps) => ({
              ...ps,
              currentState: TimerInputCurrentState.Ready,
              color: TimerColors.Green,
            }));
            // hideAllElementsExceptTimer();
          }, 1000);
        }

        if (
          timingInterval.current &&
          timerInputState.currentState === TimerInputCurrentState.Solving
        ) {
          clearInterval(timingInterval.current);
          timingInterval.current = undefined;
          updateSolve(milisecondsToFormattedTime(elapsedTime.current));
          setTimerInputState((ps) => ({
            ...ps,
            currentState: TimerInputCurrentState.Finishing,
            color: TimerColors.Red,
          }));
          holdingTimeout.current = 1;
        }
      }
    },
    [timerInputState.currentState, holdingTimeout, elapsedTime]
  );

  const handleTimerInputKeyUp = useCallback(
    (e: Event) => {
      const ev = e as KeyboardEvent;

      if (ev.key === " " || ev.type === "touchend") {
        if (holdingTimeout.current) {
          clearTimeout(holdingTimeout.current);
          holdingTimeout.current = undefined;
          if (timerInputState.currentState === TimerInputCurrentState.Ready) {
            setTimerInputState((ps) => ({
              ...ps,
              currentState: TimerInputCurrentState.Solving,
              color: TimerColors.Default,
            }));
            elapsedTime.current = 0;
            const start = Date.now();
            timingInterval.current = setInterval(
              () => (elapsedTime.current = Date.now() - start),
              10
            );
          } else {
            setTimerInputState((ps) => ({
              ...ps,
              currentState: TimerInputCurrentState.NotSolving,
              color: TimerColors.Default,
            }));
            updateCurrentSolve(
              (competitionState.currentSolveIdx + 1) %
                competitionState.noOfSolves
            );
            // revertHidingAllElementsExceptTimer();
          }
        } else {
          if (
            timerInputState.currentState === TimerInputCurrentState.Finishing
          ) {
            setTimerInputState((ps) => ({
              ...ps,
              currentState: TimerInputCurrentState.NotSolving,
              color: TimerColors.Default,
            }));
            saveResults();
            updateCurrentSolve(
              (competitionState.currentSolveIdx + 1) %
                competitionState.noOfSolves
            );
            holdingTimeout.current = undefined;
            // revertHidingAllElementsExceptTimer();
          }
        }
      }
    },
    [timerInputState.currentState, holdingTimeout, elapsedTime]
  );

  useEffect(() => {
    return () => {
      if (holdingTimeout.current) clearTimeout(holdingTimeout.current);
      if (timingInterval.current) clearInterval(timingInterval.current);
    };
  }, []);

  return (
    <TimerInputContext.Provider
      value={{
        timerInputState,
        handleTimerInputKeyDown,
        handleTimerInputKeyUp,
        // timerElementRef,
      }}
    >
      {children}
    </TimerInputContext.Provider>
  );
};

const initialState: TimerInputState = {
  currentState: TimerInputCurrentState.NotSolving,
  color: TimerColors.Default,
};
