import {
  CompetitionContextType,
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
} from "react";

import { CompetitionContext } from "./CompetitionContext";
import { milisecondsToFormattedTime } from "../utils/utils";
import useState from "react-usestateref";

export const TimerInputContext = createContext<TimerInputContextType | null>(
  null,
);

export const TimerInputProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [timerInputState, setTimerInputState, timerInputStateRef] =
    useState<TimerInputState>(initialState);
  const { updateSolve } = useContext(
    CompetitionContext,
  ) as CompetitionContextType;
  const holdingTimeout: {
    current: ReturnType<typeof setTimeout> | undefined | 1;
  } = useRef(undefined);
  const timingInterval: {
    current: ReturnType<typeof setInterval> | undefined;
  } = useRef(undefined);
  const elapsedTime = useRef(0);
  const timerRef = useRef<HTMLDivElement | null>(null);
  const timerRefStyle = useRef<any>(null);
  const allNodesStyles = useRef<string[]>([]);

  const hideAllElementsExceptTimer = () => {
    allNodesStyles.current = [];
    (document.querySelectorAll("*") as NodeListOf<HTMLElement>).forEach(
      (el: HTMLElement) => {
        if (el.tagName === "A") {
          allNodesStyles.current.push(el.style.pointerEvents);
          el.style.pointerEvents = "none";
        }
      },
    );

    if (timerRef && timerRef.current) {
      timerRefStyle.current = timerRef.current.style;
      timerRef.current.style.position = "fixed";
      timerRef.current.style.width = "100%";
      timerRef.current.style.top = "0px";
      timerRef.current.style.left = "0px";
      timerRef.current.style.height = "100%";
      timerRef.current.style.background = "white";
      timerRef.current.style.zIndex = "10000000";
    }
  };
  const revertHidingAllElementsExceptTimer = () => {
    if (timerRef && timerRef.current) {
      timerRef.current.setAttribute("style", timerRefStyle.current.cssText);
      timerRef.current?.focus();
    }

    (document.querySelectorAll("*") as NodeListOf<HTMLElement>).forEach(
      (el: HTMLElement, idx: number) => {
        if (el.tagName === "A") {
          el.style.pointerEvents = allNodesStyles.current[idx];
        }
      },
    );
  };

  const handleTimerInputKeyDown = useCallback(
    (e: Event) => {
      const ev = e as KeyboardEvent;

      if (ev.key === " " || ev.type === "touchstart") {
        if (
          !holdingTimeout.current &&
          !timingInterval.current &&
          timerInputStateRef.current.currentState ===
          TimerInputCurrentState.NotSolving
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
            hideAllElementsExceptTimer();
          }, 1000);
        }

        if (
          timingInterval.current &&
          timerInputStateRef.current.currentState ===
          TimerInputCurrentState.Solving
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
    [timerInputState.currentState, holdingTimeout, elapsedTime],
  );

  const handleTimerInputKeyUp = useCallback(
    (e: Event, handleSaveResults?: (moveIndex: boolean) => void) => {
      const ev = e as KeyboardEvent;

      if (ev.key === " " || ev.type === "touchend") {
        if (holdingTimeout.current) {
          clearTimeout(holdingTimeout.current);
          holdingTimeout.current = undefined;
          if (
            timerInputStateRef.current.currentState ===
            TimerInputCurrentState.Ready
          ) {
            setTimerInputState((ps) => ({
              ...ps,
              currentState: TimerInputCurrentState.Solving,
              color: TimerColors.Default,
            }));
            elapsedTime.current = 0;
            const start = Date.now();
            timingInterval.current = setInterval(
              () => (elapsedTime.current = Date.now() - start),
              10,
            );
          } else {
            if (
              handleSaveResults !== undefined &&
              timerInputStateRef.current.currentState ===
              TimerInputCurrentState.Finishing
            )
              handleSaveResults(true);
            setTimerInputState((ps) => ({
              ...ps,
              currentState: TimerInputCurrentState.NotSolving,
              color: TimerColors.Default,
            }));
            revertHidingAllElementsExceptTimer();
          }
        } else if (
          timerInputStateRef.current.currentState ===
          TimerInputCurrentState.Finishing
        ) {
          if (handleSaveResults !== undefined) handleSaveResults(true);
          setTimerInputState((ps) => ({
            ...ps,
            currentState: TimerInputCurrentState.NotSolving,
            color: TimerColors.Default,
          }));
          holdingTimeout.current = undefined;
          revertHidingAllElementsExceptTimer();
        }
      }
    },
    [timerInputState.currentState, holdingTimeout, elapsedTime],
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
        timerRef,
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
