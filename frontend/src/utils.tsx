import { CompetitionData, FilterValue } from "./Types";

const allCompetitionData = async (): Promise<CompetitionData[]> => {
    const result: CompetitionData[] = [];
    let startdate: Date = new Date();
    startdate.setDate(startdate.getDate() - 23);
    let enddate: Date = new Date(startdate);
    enddate.setDate(enddate.getDate() + 7);

    for (let i = 0; i < 10; i++) {
        if (result.length > 0) {
            startdate = new Date(result[result.length - 1].enddate)
            enddate = new Date(startdate)
            enddate.setDate(enddate.getDate() + 7)
        }

        result.push(
            {
                'id': 'WeeklyCompetition' + (i + 1).toString(),
                'name': 'Weekly Competition ' + (i + 1).toString(),
                startdate,
                enddate
            }
        )
    }

    return result;
}

export const loadCompetitionData = async (filterValue: FilterValue): Promise<CompetitionData[]> => {
    return allCompetitionData().then(res => res.filter((c: CompetitionData, i: number) => {
        const curdate: Date = new Date();
        const startdate: Date = c.startdate;
        const enddate: Date = c.enddate;

        return ((filterValue === FilterValue.Current && startdate.getTime() <= curdate.getTime() && curdate.getTime() <= enddate.getTime()) ||
            (filterValue === FilterValue.Past && enddate.getTime() < curdate.getTime()) ||
            (filterValue === FilterValue.Future && startdate.getTime() > curdate.getTime()))
    }))
}

export const getCompetitionById = async (id: string | undefined): Promise<CompetitionData | undefined> => {
    const competitions: CompetitionData[] = await allCompetitionData();
    return competitions.find((c: CompetitionData) => c.id === id);
}