export type CompetitionData = {
    id: string,
    name: string,
    startdate: Date,
    enddate: Date,
    events: CompetitionEvent[]
}

export enum FilterValue {
    Current = "Current", Past = "Past", Future = "Future"
}

export type CompetitionEvent = {
    id: number
    displayname: string
    format: string
    iconcode: string
}

export type CompetitionState = CompetitionData & {
    currentEventIdx: number,
    currentSolveIdx: number,
    noOfSolves: number
}

export type CompetitionContextType = {
    competitionState: CompetitionState,
    updateBasicInfo: (info: CompetitionData) => void
    updateCurrentEvent: (idx: number) => void
    updateCurrentSolve: (idx: number) => void
}