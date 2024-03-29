import { AuthContextType, CompetitionContextType, CompetitionData, CompetitionState, ResultEntry } from "../../Types";
import React, { ReactNode, createContext, useContext, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import { getResultsFromCompetitionAndEvent } from "../../utils";

export const CompetitionContext = createContext<CompetitionContextType | null>(null);

export const CompetitionProvider: React.FC<{ children?: ReactNode }> = ({ children }) => {
    const [competitionState, setCompetitionState] = useState<CompetitionState>(initialState);
    const { authState } = useContext(AuthContext) as AuthContextType;

    const updateBasicInfo = (info: CompetitionData) => {
        const match = info.events[competitionState.currentEventIdx].format.match(/\d+$/)?.[0]
        const noOfSolves = match ? parseInt(match) : 1
        setCompetitionState({...competitionState, ...info, noOfSolves: noOfSolves});
    }

    const updateCurrentEvent = (idx: number) => {
        const match = competitionState.events[idx].format.match(/\d+$/)?.[0]
        const noOfSolves = match ? parseInt(match) : 1
        const results = getResultsFromCompetitionAndEvent(authState.token, competitionState.id, competitionState.events[competitionState.currentEventIdx]);
        setCompetitionState({...competitionState, currentEventIdx: idx, noOfSolves: noOfSolves, currentSolveIdx: 0, ...results });
    }

    const updateCurrentSolve = (idx: number) => setCompetitionState({...competitionState, currentSolveIdx: idx });

    const saveResults = () => console.log('Your results were saved!');

    const updateSolve = (newTime: string) => {
        const solveProp: keyof ResultEntry = `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
        console.log(newTime);
        setCompetitionState({
            ...competitionState,
            results:
            {
                ...competitionState.results,
                [solveProp]: newTime
            }
        });
    } 
    
    return (
        <CompetitionContext.Provider value={{competitionState, updateBasicInfo, updateCurrentEvent, updateCurrentSolve, saveResults, updateSolve}}>
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
    currentSolveIdx: 0,
    scrambles: [],
    results: {
        id: 0,
        userid: 0,
        solve1: '',
        solve2: '',
        solve3: '',
        solve4: '',
        solve5: '',
        comment: '',
        statusid: 0,
    }
};