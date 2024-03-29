export type CompetitionData = {
    id: string,
    name: string,
    startdate: Date,
    enddate: Date,
    events: CompetitionEvent[],
    scrambles: string[][],
    results?: ResultEntry
}

export enum FilterValue {
    Current = "Current", Past = "Past", Future = "Future"
}

export type CompetitionEvent = {
    id: number
    displayname: string
    format: string
    iconcode: string
    puzzlecode: string
}

export enum InputMethod {
    Manual, Timer
}

export type CompetitionState = CompetitionData & {
    currentEventIdx: number,
    currentSolveIdx: number,
    noOfSolves: number,
    inputMethod: InputMethod,
    results: ResultEntry
}

export type CompetitionContextType = {
    competitionState: CompetitionState,
    updateBasicInfo: (info: CompetitionData) => void
    updateCurrentEvent: (idx: number) => void
    updateCurrentSolve: (idx: number) => void
    saveResults: () => void
    updateSolve: (newTime: string) => void
    toggleInputMethod: () => void
}

export type ResultEntry = {
    id: number,
    userid: number,
    solve1: string,
    solve2: string,
    solve3: string,
    solve4: string,
    solve5: string,
    comment: string,
    statusid: number,
}

export enum ResultEntrySolves {
    solve1, solve2, solve3, solve4, solve5
}

export type AuthState = {
    token: string
}

export type AuthContextType = {
    authState: AuthState,
    updateAuthToken: (newToken: string) => void
}