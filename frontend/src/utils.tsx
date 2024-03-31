import { CompetitionData, CompetitionEvent, FilterValue, ResultEntry, User } from "./Types";

const events: CompetitionEvent[] = [
    {
        'id': 1,
        'displayname': '3x3x3',
        'format': 'ao5',
        'iconcode': '333',
        'puzzlecode': '3x3x3'
    },
    {
        'id': 2,
        'displayname': '2x2x2',
        'format': 'ao5',
        'iconcode': '222',
        'puzzlecode': '2x2x2'
    },
    {
        'id': 3,
        'displayname': '6x6x6',
        'format': 'mo3',
        'iconcode': '666',
        'puzzlecode': '6x6x6'
    },
    {
        'id': 4,
        'displayname': 'Mega',
        'format': 'ao5',
        'iconcode': 'mega',
        'puzzlecode': 'megaminx'
    },
    {
        'id': 5,
        'displayname': 'Pyra',
        'format': 'ao5',
        'iconcode': 'pyra',
        'puzzlecode': 'pyraminx'
    },
    {
        'id': 6,
        'displayname': '3BLD',
        'format': 'bo3',
        'iconcode': '333bld',
        'puzzlecode': '3x3x3'
    },
    {
        'id': 7,
        'displayname': 'FMC',
        'format': 'mo3',
        'iconcode': 'fmc',
        'puzzlecode': '3x3x3'
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
        "R-- D-- R-- D-- R-- D-- R++ D++ R-- D-- U'\n  R++ D-- R-- D-- R-- D++ R++ D-- R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D-- R++ D++ U\n  R++ D++ R-- D-- R++ D-- R++ D++ R++ D++ U \n  R++ D++ R++ D-- R-- D++ R-- D++ R-- D++ U \n  R++ D-- R-- D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D++ R-- D++ R++ D++ R-- D-- R-- D-- U'\n",
        "R-- D++ R-- D-- R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D-- R-- D-- R-- D++ R-- D++ U \n  R-- D-- R-- D-- R-- D-- R++ D++ R++ D-- U'\n  R++ D-- R++ D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D-- R-- D-- R++ D-- R-- D-- R++ D-- U'\n  R++ D-- R++ D-- R-- D++ R-- D-- R++ D++ U \n  R++ D-- R++ D-- R++ D-- R++ D-- R-- D-- U'\n",
        "R-- D++ R-- D++ R-- D-- R-- D-- R-- D++ U \n  R++ D-- R++ D-- R-- D-- R-- D++ R++ D++ U \n  R-- D-- R-- D-- R++ D++ R-- D-- R-- D++ U \n  R-- D++ R-- D-- R++ D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D-- R++ D++ R-- D++ R++ D++ U \n  R-- D++ R++ D++ R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D++ R-- D-- U'\n",
        "R-- D++ R++ D-- R++ D++ R-- D++ R-- D-- U'\n  R-- D++ R++ D-- R-- D++ R-- D-- R-- D-- U'\n  R-- D-- R-- D-- R-- D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D++ R++ D++ R-- D-- R-- D++ U \n  R++ D++ R-- D-- R-- D-- R++ D-- R++ D++ U \n  R-- D++ R++ D++ R-- D++ R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R-- D-- R-- D-- R++ D-- U'\n",
        "R++ D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D++ R++ D-- R++ D-- R++ D++ R-- D-- U'\n  R-- D-- R-- D++ R-- D++ R++ D++ R-- D++ U \n  R++ D++ R-- D++ R++ D++ R-- D-- R++ D++ U \n  R-- D-- R++ D++ R-- D++ R++ D-- R-- D-- U'\n  R++ D-- R++ D++ R-- D-- R++ D++ R++ D-- U'\n",
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

const results: { [key: string]: ResultEntry } = {
    '3x3x3': {
        'id': 1,
        'userid': 1,
        'solve1': '12.55',
        'solve2': '10.14',
        'solve3': '8.81',
        'solve4': 'DNF',
        'solve5': '14.43',
        'comment': '',
        'statusid': 1,
    },
    '2x2x2': {
        'id': 2,
        'userid': 1,
        'solve1': '2.55',
        'solve2': '1.14',
        'solve3': '8.81',
        'solve4': '2.00',
        'solve5': '1.43',
        'comment': '',
        'statusid': 2,
    },
    '6x6x6': {
        'id': 3,
        'userid': 1,
        'solve1': '2:00.55',
        'solve2': '1:59.14',
        'solve3': '1:58.80',
        'solve4': '',
        'solve5': '',
        'comment': '',
        'statusid': 2,
    },
    'Mega': {
        'id': 4,
        'userid': 1,
        'solve1': '42.55',
        'solve2': '41.14',
        'solve3': '48.81',
        'solve4': '42.00',
        'solve5': '41.43',
        'comment': '',
        'statusid': 2,
    },
    'Pyra': {
        'id': 5,
        'userid': 1,
        'solve1': '2.13',
        'solve2': '1.01',
        'solve3': '2.99',
        'solve4': '2.00',
        'solve5': '2.69',
        'comment': '',
        'statusid': 2,
    },
    '3BLD': {
        'id': 6,
        'userid': 1,
        'solve1': 'DNF',
        'solve2': '1:00.05',
        'solve3': 'DNS',
        'solve4': '',
        'solve5': '',
        'comment': '',
        'statusid': 2,
    },
    'FMC': {
        'id': 2,
        'userid': 1,
        'solve1': 'R U R\' U\'',
        'solve2': 'abc',
        'solve3': '',
        'solve4': '',
        'solve5': '',
        'comment': '',
        'statusid': 2,
    },
}

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

export const getResultsFromCompetitionAndEvent = async (token: string, id: string | undefined, event: CompetitionEvent): Promise<ResultEntry> => {
    const resultEntry = results[event.displayname];
    return resultEntry;    
}

const formattedToMiliseconds = (formattedTime: string): number => {
    let res = 0;

    const formattedTimeSplit = formattedTime.split('.');
    const wholePart = formattedTimeSplit[0].split(':').reverse(), decimalPart = formattedTimeSplit[1];

    res += parseInt(decimalPart) * 10;
    if (wholePart.length > 0)
        res += parseInt(wholePart[0]) * 1000;
    if (wholePart.length > 1)
        res += 60 * parseInt(wholePart[1]) * 1000;
    if (wholePart.length > 2)
        res += 60 * 60 * parseInt(wholePart[2]) * 1000;
    if (wholePart.length > 3)
        res += 24 * 60 * 60 * parseInt(wholePart[3]) * 1000;

    return res;
}

export const milisecondsToFormattedTime = (toFormat: number): string => {
    if (toFormat === -1) {
        return "DNF";
    }

    if (toFormat === -2) {
        return "DNS";
    }

    let res = [];

    let pw = 1000 * 60 * 60 * 24;
    for (const mul of [24, 60, 60, 1000, 1]) {
        const toPush = Math.floor(toFormat / pw).toString();
        res.push(mul === 1 ? toPush.padStart(3, '0') : toPush);
        toFormat %= pw;
        pw = Math.floor(pw / mul);
    }

    res[res.length - 1] = res[res.length - 1].slice(0, res[res.length - 1].length - 1);
    let sliceIdx = 0;
    while (sliceIdx < res.length - 2 && res[sliceIdx] === '0')
        sliceIdx += 1;
    res = res.slice(sliceIdx);

    let resString = "";
    let resIdx: number;
    for (resIdx = 0; resIdx < res.length - 1; resIdx++) {
        resString += resIdx > 0 ? res[resIdx].padStart(2, '0') : res[resIdx];
        resString += resIdx == res.length - 2 ? '.' : ':';
    }
    resString += res[resIdx].padStart(2, '0');

    return resString;
}

export const reformatWithPenalties = (oldFormattedTime: string, penalty: string) => {
    if (oldFormattedTime === "DNF") {
        return oldFormattedTime;
    }

    if (oldFormattedTime === "DNS") {
        return penalty === "DNF" ? "DNF" : "DNS";
    }

    let miliseconds = formattedToMiliseconds(oldFormattedTime);

    if (penalty === "DNF") {
        miliseconds = -1;
    } else {
        miliseconds += parseInt(penalty) * 1000;
    }

    let newFormattedTime = milisecondsToFormattedTime(miliseconds);

    return newFormattedTime;
}

const users: User[] = [
    {
        id: 1,
        name: "Janko Hrasko",
        isadmin: true,
    },
    {
        id: 2,
        name: "Ferko Mrkvicka",
        isadmin: false
    }
]

export const getUsers = async () => {
    return users;
}

export const updateUsers = async (users: User[]) => {
    
}