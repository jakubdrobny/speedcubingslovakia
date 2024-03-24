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
    },
    {
        'id': 7,
        'displayname': 'FMC',
        'format': 'mo3',
        'iconcode': 'fmc'
    }
]

const scrambles: string[][] = [
    [
        "R2 U B2 D' R2 U L2 B' D U' L' F2 U' L F' D'",
        "F2 B U2 F2 D' B D2 L R2 U' F2 D F2 U' L2 U R2 U2 B'",
        "L D2 R2 B2 U' R2 D B2 U2 L2 R2 F2 R F' D' R' B U2 B",
        "B' L2 F' L2 R2 D2 R2 F' L2 B' L2 B R' U B' F2 D' R U' B' R",
        "U L2 B2 D' L2 F2 L2 U2 L2 U R2 U2 L' D' R2 B' D2 B2 D2"
    ],
    [
        "R' U2 R F2 U' R' U2 R' F2",
        "F R2 U2 R F2 U F' R' F2",
        "U R2 U F' R2 U2 F' U2 R'",
        "R U2 R2 U' F R2 U F2 U2",
        "F2 U2 R' F R F2 U' F R2 F2"
    ],
    [
        "R' U2 Uw2 3Rw2 Fw' 3Fw' D' Fw 3Fw2 R' Uw2 Lw Dw R D' Bw2 R2 U2 Rw2 3Rw U F' L 3Fw2 R' F2 3Rw2 D Dw' Lw' B R' Fw Bw2 3Uw2 Fw' U2 3Fw' Fw' D L2 F2 Uw 3Fw2 3Uw' Bw Uw2 R2 Rw' 3Fw2 R Lw B Dw2 U2 Bw 3Rw R2 3Fw Fw R' 3Uw' Fw Uw 3Rw2 L2 Lw' U2 Lw U2 Bw' F 3Fw Dw R2 Rw2 L' 3Rw 3Fw Fw2",
        "L' R2 Bw F 3Uw D' 3Fw Lw2 Rw' Bw' R Bw2 D2 Bw' D2 F' D2 L2 Rw2 Lw' 3Rw' F Bw' D2 3Uw Bw2 Lw' U 3Uw Rw Bw2 Lw' F' B Bw2 U 3Fw' F2 R2 Bw' Fw 3Rw2 Uw Fw R F2 Lw U2 Bw2 Uw' B Uw' Lw 3Uw2 Dw F Uw' F2 L2 3Fw Dw' Bw2 Rw2 Lw' Dw' F' Lw' B' Rw' D' Dw2 Fw Lw 3Fw Dw' D2 F' D2 3Fw2 Fw'",
        "R' Fw' D Fw Uw' U2 F' L2 Rw' 3Rw' Lw' Dw U2 Lw' 3Rw' R' B2 3Uw2 Uw2 F' 3Uw Rw2 F2 R2 Lw 3Uw2 Uw' 3Fw2 Fw2 D' 3Uw Fw 3Rw' Fw Dw2 3Rw2 Lw L F Lw' B2 Uw' 3Fw2 Dw D Lw' F2 R Bw' Rw' Fw' 3Rw' Fw 3Rw' F2 Fw 3Fw2 D2 F L B2 Lw' L2 D2 3Fw' 3Uw Uw Rw' Uw F' Rw' L' U Fw2 U Uw 3Uw F R 3Rw'"
    ],
    [
        "R-- D-- R-- D-- R-- D-- R++ D++ R-- D-- U'\n  R++ D-- R-- D-- R-- D++ R++ D-- R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D-- R++ D++ U~\n  R++ D++ R-- D-- R++ D-- R++ D++ R++ D++ U~\n  R++ D++ R++ D-- R-- D++ R-- D++ R-- D++ U~\n  R++ D-- R-- D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D++ R-- D++ R++ D++ R-- D-- R-- D-- U'\n",
        "R-- D++ R-- D-- R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D-- R-- D-- R-- D++ R-- D++ U~\n  R-- D-- R-- D-- R-- D-- R++ D++ R++ D-- U'\n  R++ D-- R++ D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D-- R-- D-- R++ D-- R-- D-- R++ D-- U'\n  R++ D-- R++ D-- R-- D++ R-- D-- R++ D++ U~\n  R++ D-- R++ D-- R++ D-- R++ D-- R-- D-- U'\n",
        "R-- D++ R-- D++ R-- D-- R-- D-- R-- D++ U~\n  R++ D-- R++ D-- R-- D-- R-- D++ R++ D++ U~\n  R-- D-- R-- D-- R++ D++ R-- D-- R-- D++ U~\n  R-- D++ R-- D-- R++ D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D-- R++ D++ R-- D++ R++ D++ U~\n  R-- D++ R++ D++ R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D++ R-- D-- U'\n",
        "R-- D++ R++ D-- R++ D++ R-- D++ R-- D-- U'\n  R-- D++ R++ D-- R-- D++ R-- D-- R-- D-- U'\n  R-- D-- R-- D-- R-- D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D++ R++ D++ R-- D-- R-- D++ U~\n  R++ D++ R-- D-- R-- D-- R++ D-- R++ D++ U~\n  R-- D++ R++ D++ R-- D++ R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R-- D-- R-- D-- R++ D-- U'\n",
        "R++ D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D++ R++ D-- R++ D-- R++ D++ R-- D-- U'\n  R-- D-- R-- D++ R-- D++ R++ D++ R-- D++ U~\n  R++ D++ R-- D++ R++ D++ R-- D-- R++ D++ U~\n  R-- D-- R++ D++ R-- D++ R++ D-- R-- D-- U'\n  R++ D-- R++ D++ R-- D-- R++ D++ R++ D-- U'\n",
    ],
    [
        "R U B U' B' L' U R' l' r' u'",
        "B' U' R B R' L B' U r u'",
        "L' U L R' L B' R B' l b u",
        "B' R U B L U' B' L l' r' b' u",
        "B L U' L B' R' U' B' l r b' u"
    ],
    [
        "F2 L2 D' F2 D2 L2 U B2 U F2 U B' U L' B R2 D2 F R' D2 Rw' Uw'",
        "L2 B D2 B2 U2 L2 B' R2 F' D2 F L2 D' B' F' U' L R' B U' F Rw Uw",
        "R' D L2 B2 U2 R2 B2 L' U2 L2 B2 D2 R' F' L U2 R2 B F D' Rw"
    ],
    [
        "R' U' F D2 U2 L2 B U2 B2 L2 D2 U2 F' R B2 F R U F U' B2 F2 R' U' F",
        "R' U' F D2 F2 D2 B2 L2 F2 L2 D F2 B' R F' L F U2 F' D R' U' F",
        "R' U' F L U2 R' D2 L B2 R F2 L B2 U2 F2 B' L F' R D F U L2 B2 R' U' F"
    ]
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
                events,
                scrambles
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