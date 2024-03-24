import { CompetitionContextType, CompetitionData, CompetitionState } from "../../Types";
import React, { ReactNode, createContext, useState } from "react";

export const CompetitionContext = createContext<CompetitionContextType | null>(null);

export const CompetitionProvider: React.FC<{ children?: ReactNode }> = ({ children }) => {
    const [competitionState, setCompetitionState] = useState<CompetitionState>(initialState);

    const updateBasicInfo = (info: CompetitionData) => {
        const match = info.events[competitionState.currentEventIdx].format.match(/\d+$/)?.[0]
        const noOfSolves = match ? parseInt(match) : 1
        setCompetitionState({...competitionState, ...info, noOfSolves: noOfSolves});
    }

    const updateCurrentEvent = (idx: number) => {
        const match = competitionState.events[idx].format.match(/\d+$/)?.[0]
        const noOfSolves = match ? parseInt(match) : 1
        setCompetitionState({...competitionState, currentEventIdx: idx, noOfSolves: noOfSolves, currentSolveIdx: 0 });
    }

    const updateCurrentSolve = (idx: number) => setCompetitionState({...competitionState, currentSolveIdx: idx });

    return (
        <CompetitionContext.Provider value={{competitionState, updateBasicInfo, updateCurrentEvent, updateCurrentSolve}}>
            {children}
        </CompetitionContext.Provider>
    );
}

const initialState: CompetitionState = {
    id: "",
    name: "",
    startdate: new Date(),
    enddate: new Date(),
    events: [],
    currentEventIdx: 0,
    noOfSolves: 1,
    currentSolveIdx: 0
};