import { CompetitionData, CompetitionEvent, FilterValue } from "./Types";

const events: CompetitionEvent[] = [
    {
        'id': 1,
        'displayname': '3x3x3',
        'format': 'ao5',
        'iconcode': '333'
    },
    {
        'id': 2,
        'displayname': '2x2x2',
        'format': 'ao5',
        'iconcode': '222'
    },
    {
        'id': 3,
        'displayname': '6x6x6',
        'format': 'mo3',
        'iconcode': '666'
    },
    {
        'id': 4,
        'displayname': 'Mega',
        'format': 'ao5',
        'iconcode': 'mega'
    },
    {
        'id': 5,
        'displayname': 'Pyra',
        'format': 'ao5',
        'iconcode': 'pyra'
    },
    {
        'id': 6,
        'displayname': '3BLD',
        'format': 'bo3',
        'iconcode': '333bld'
    }
]

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
                enddate,
                events
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