import { AuthContextType, CompetitionContextType, CompetitionData, CompetitionState, InputMethod, ResultEntry } from "../../Types";
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

    const updateCurrentEvent = async (idx: number) => {
        const match = competitionState.events[idx].format.match(/\d+$/)?.[0]
        const noOfSolves = match ? parseInt(match) : 1
        const resultEntry = await getResultsFromCompetitionAndEvent(authState.token, competitionState.id, competitionState.events[idx]);
        setCompetitionState({...competitionState, currentEventIdx: idx, noOfSolves: noOfSolves, currentSolveIdx: 0, results: resultEntry });
    }

    const updateCurrentSolve = (idx: number) => setCompetitionState({...competitionState, currentSolveIdx: idx });

    const saveResults = () => {
        const solveProp: keyof ResultEntry= `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
        const formattedTime = competitionState.results[solveProp].toString();
        console.log(`You saved a time of ${formattedTime}!`);
    }

    const updateSolve = (newTime: string) => {
        const solveProp: keyof ResultEntry = `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
        setCompetitionState({
            ...competitionState,
            results:
            {
                ...competitionState.results,
                [solveProp]: newTime
            }
        });
    }

    const toggleInputMethod = () => {
        if (competitionState.currentEventIdx < competitionState.events.length && competitionState.events[competitionState.currentEventIdx].displayname !== "FMC") {
            setCompetitionState({...competitionState, inputMethod: competitionState.inputMethod === InputMethod.Manual ? InputMethod.Timer : InputMethod.Manual})
        }
    }
    
    return (
        <CompetitionContext.Provider value={{competitionState, updateBasicInfo, updateCurrentEvent, updateCurrentSolve, saveResults, updateSolve, toggleInputMethod}}>
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
    inputMethod: InputMethod.Manual,
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