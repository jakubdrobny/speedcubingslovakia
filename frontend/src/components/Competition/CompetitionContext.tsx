import { CompetitionContextType, CompetitionData, CompetitionState } from "../../Types";
import React, { ReactNode, createContext, useState } from "react";

export const CompetitionContext = createContext<CompetitionContextType | null>(null);

export const CompetitionProvider: React.FC<{ children?: ReactNode }> = ({ children }) => {
    const [competitionState, setCompetitionState] = useState<CompetitionState>(initialState);

    const updateBasicInfo = (info: CompetitionData) => setCompetitionState({...competitionState, ...info});

    const updateCurrentEvent = (idx: number) => setCompetitionState({...competitionState, currentEventIdx: idx });

    return (
        <CompetitionContext.Provider value={{competitionState, updateBasicInfo, updateCurrentEvent}}>
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
    currentEventIdx: 0
};