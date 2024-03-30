import { CompetitionContextType, TimerColors, TimerInputContextType, TimerInputCurrentState, TimerInputState } from "../Types";
import React, { ReactNode, createContext, useCallback, useContext, useEffect, useRef, useState } from "react";

import { CompetitionContext } from "../components/Competition/CompetitionContext";
import { milisecondsToFormattedTime } from "../utils";

export const TimerInputContext = createContext<TimerInputContextType | null>(null);

export const TimerInputProvider: React.FC<{ children?: ReactNode }> = ({ children }) => {
    const [timerInputState, setTimerInputState] = useState<TimerInputState>(initialState);
    const { saveResults, updateSolve } = useContext(CompetitionContext) as CompetitionContextType
    let holdingTimeout: { current: ReturnType<typeof setTimeout> | undefined | 1} = useRef(undefined);
    let timingInterval: { current: ReturnType<typeof setInterval> | undefined } = useRef(undefined);
    const elapsedTime = useRef(0)

    const handleTimerInputKeyDown = useCallback((e: Event) => {
        const ev = e as KeyboardEvent

        if (ev.key === ' ') {
            if (!holdingTimeout.current && !timingInterval.current && timerInputState.currentState === TimerInputCurrentState.NotSolving) {
                setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.GettingReady, color: TimerColors.Red}))
                holdingTimeout.current = setTimeout(() => {
                    setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.Ready, color: TimerColors.Green}))
                }, 1000);
            }

            if (timingInterval.current && timerInputState.currentState === TimerInputCurrentState.Solving) {
                clearInterval(timingInterval.current);
                timingInterval.current = undefined;
                updateSolve(milisecondsToFormattedTime(elapsedTime.current))
                setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.Finishing, color: TimerColors.Red}))
                holdingTimeout.current = 1;
            }
        }
    }, [timerInputState.currentState, holdingTimeout, elapsedTime])

    const handleTimerInputKeyUp = useCallback((e: Event) => {
        const ev = e as KeyboardEvent

        if (ev.key === ' ') {
            if (holdingTimeout.current) {
                clearTimeout(holdingTimeout.current)
                holdingTimeout.current = undefined;
                if (timerInputState.currentState === TimerInputCurrentState.Ready) {
                    setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.Solving, color: TimerColors.Default}))
                    elapsedTime.current = 0;
                    const start = Date.now();
                    timingInterval.current = setInterval(() => elapsedTime.current = Date.now() - start, 10);
                } else {
                    setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.NotSolving, color: TimerColors.Default}))
                }
            } else {
                if (timerInputState.currentState === TimerInputCurrentState.Finishing) {
                    setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.NotSolving, color: TimerColors.Default}))
                    saveResults()
                    holdingTimeout.current = undefined;
                }
            }
        }
    }, [timerInputState.currentState, holdingTimeout, elapsedTime])

    useEffect(() => {
        return () => {
            if (holdingTimeout.current)
                clearTimeout(holdingTimeout.current)
            if (timingInterval.current)
                clearInterval(timingInterval.current)
        }
    }, []);

    return (
        <TimerInputContext.Provider value={{timerInputState, handleTimerInputKeyDown, handleTimerInputKeyUp}}>
            {children}
        </TimerInputContext.Provider>
    );
}

const initialState: TimerInputState = {
    currentState: TimerInputCurrentState.NotSolving,
    color: TimerColors.Default
};