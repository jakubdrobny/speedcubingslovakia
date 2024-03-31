import { AuthContextType, CompetitionContextType, CompetitionData, CompetitionState, InputMethod, PenaltyType, ResultEntry } from "../../Types";
import React, { ReactNode, createContext, useContext, useState } from "react";
import { getResultsFromCompetitionAndEvent, reformatWithPenalties } from "../../utils";

import { AuthContext } from "../../context/AuthContext";

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
        setCompetitionState(ps => ({
            ...ps,
            currentEventIdx: idx,
            noOfSolves: noOfSolves,
            currentSolveIdx: 0,
            results: resultEntry,
            penalties: Array(5).fill('0'),
            inputMethod: competitionState.events[idx].displayname === "FMC" ? InputMethod.Manual : ps.inputMethod
        }));
    }

    const updateCurrentSolve = (idx: number) => setCompetitionState({...competitionState, currentSolveIdx: idx });

    const saveResults = () => {
        const solveProp: keyof ResultEntry= `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
        const formattedTime = competitionState.results[solveProp].toString();
        const finalFormattedTime = reformatWithPenalties(formattedTime, competitionState.penalties[competitionState.currentSolveIdx]);
        console.log(`You saved a time of ${finalFormattedTime}!`);
    }

    const addPenalty = (newPenalty: string) => {
        const curSolveIdx = competitionState.currentSolveIdx
        const oldPenalty = competitionState.penalties[curSolveIdx]
        
        if (newPenalty === "DNF") {
            newPenalty = oldPenalty === "DNF" ? "0" : "DNF";
        } else {
            newPenalty = oldPenalty === "DNF" ? "2" : ((parseInt(oldPenalty) + parseInt(newPenalty)) % 18).toString();
        }

        setCompetitionState(ps => ({...ps, penalties: ps.penalties.map((val: string, idx: number) => idx == curSolveIdx ? newPenalty : val)}));
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
            setCompetitionState(ps => ({
                ...ps,
                inputMethod: ps.inputMethod === InputMethod.Manual ? InputMethod.Timer : InputMethod.Manual,
                penalties: ps.penalties.map((val: string, idx: number) => idx == ps.currentSolveIdx ? '0' : val)
            }))
        }
    }
    
    return (
        <CompetitionContext.Provider value={{
            competitionState, updateBasicInfo, updateCurrentEvent,
            updateCurrentSolve, saveResults, toggleInputMethod,
            addPenalty, updateSolve
        }}>
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
    },
    penalties: Array(5).fill('0')
};