import { CompetitionContextType, TimerColors, TimerInputContextType, TimerInputCurrentState, TimerInputState } from "../Types";
import React, { ReactNode, createContext, useCallback, useContext, useEffect, useRef, useState } from "react";

import { CompetitionContext } from "../components/Competition/CompetitionContext";

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
                updateSolve(formatTime(elapsedTime.current))
                setTimerInputState(ps => ({...ps, currentState: TimerInputCurrentState.Finishing, color: TimerColors.Red}))
                holdingTimeout.current = 1;
            }
        }
    }, [timerInputState.currentState, holdingTimeout, elapsedTime])

    const formatTime = (toFormat: number) => {
        let res = [];

        let pw = 1000 * 60 * 60 * 24;
        for (const mul of [24, 60, 60, 1000, 1]) {
            const toPush = Math.floor(toFormat / pw).toString();
            res.push(mul === 1 ? toPush.padStart(3, '0') : toPush);
            toFormat %= pw;
            pw = Math.floor(pw / mul);
        }

        res[res.length - 1] = res[res.length - 1].slice(0, res[res.length - 1].length - 1);
        let sliceIdx = 0;
        while (sliceIdx < res.length - 2 && res[sliceIdx] === '0')
            sliceIdx += 1;
        res = res.slice(sliceIdx);

        let resString = "";
        let resIdx: number;
        for (resIdx = 0; resIdx < res.length - 1; resIdx++) {
            resString += resIdx > 0 ? res[resIdx].padStart(2, '0') : res[resIdx];
            resString += resIdx == res.length - 2 ? '.' : ':';
        }
        resString += res[resIdx].padStart(2, '0');

        return resString;
    }

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