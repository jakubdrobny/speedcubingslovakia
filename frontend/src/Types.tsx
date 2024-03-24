export type CompetitionData = {
    id: String,
    name: String,
    startdate: Date,
    enddate: Date
}

export enum FilterValue {
    Current = "Current", Past = "Past", Future = "Future"
}