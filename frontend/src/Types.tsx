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
    results: ResultEntry,
    penalties: string[]
}

export type CompetitionContextType = {
    competitionState: CompetitionState,
    updateBasicInfo: (info: CompetitionData) => void
    updateCurrentEvent: (idx: number) => void
    updateCurrentSolve: (idx: number) => void
    saveResults: () => void
    updateSolve: (newTime: string) => void
    toggleInputMethod: () => void
    addPenalty: (newPenalty: string) => void
    updateCompetitionName: (newName: string) => void, 
    updateCompetitionStartDate: (newStartingDate: Date) => void,
    updateCompetitionEndDate: (newEndingDate: Date) => void,
    updateCompetitionEvents: (newEvents: CompetitionEvent[]) => void
    setCompetitionState: (newState: CompetitionState) => void
}

export type ResultEntry = {
    id: number,
    userid: number,
    username: string,
    competitionid: number,
    competitionname: string,
    eventid: number,
    eventname: string,
    iconcode: string,
    format: string,
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
    token: string,
    authenticated: boolean,
    admin: boolean
}

export type AuthContextType = {
    authState: AuthState,
    updateAuthToken: (newToken: string) => void
}

export type TimerInputContextType = {
    timerInputState: TimerInputState,
    handleTimerInputKeyDown: EventListener
    handleTimerInputKeyUp: EventListener
}

export type TimerInputState = {
    currentState: TimerInputCurrentState,
    color: string
}

export enum TimerInputCurrentState {
    NotSolving, GettingReady, Ready, Solving, Finishing
}

export enum TimerColors {
    Default = '#32383E',
    Red = 'red',
    Green = 'green'
}

export enum PenaltyType { 
    PlusTwo, DNF
}

export type User = {
    id: number
    name: string
    isadmin: boolean
}

export enum DashboardPanel {
    ManageRoles, None
}

export type CompetitionEditProps = {
    edit: boolean
}