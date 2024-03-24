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