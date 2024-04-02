import { CompetitionData, CompetitionEvent, CompetitionState, FilterValue, InputMethod, ResultEntry, ResultsStatus, User } from "./Types";

import axios from 'axios';

export const loadFilteredCompetitions = async (filterValue: FilterValue): Promise<CompetitionData[]> => {
    const response = await axios.get(`/api/competitions/${filterValue}`)
    console.log(response.data, 'data')
    return response.data;
}

export const getCompetitionById = async (id: string | undefined): Promise<CompetitionData> => {
    if (id === undefined) {
        return Promise.reject("Invalid competition id.");
    }
    
    const response = await axios.get(`/api/competitions/${id}`)
    return response.data;
}

export const getResultsFromCompetitionAndEvent = async (token: string, id: string | undefined, event: CompetitionEvent): Promise<ResultEntry> => {
    const response = await axios.get(`/api/results/${id}/${event.displayname}`)
    return response.data;
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

let users: User[] = [
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

export const updateUsers = async (newUsers: User[]) => {
    users = [...newUsers];
}

export const getAvailableEvents = async () => {
    const response = await axios.get('/api/events')
    return response.data;
}

export const updateCompetition = (state: CompetitionState, edit: boolean) => {
    console.log('editujem')
    return;
}

export const getResults = async (competitorName: string, competitionName: string, competeEvent: CompetitionEvent | undefined) => {
    const response = await axios.get(
        '/api/results', { data: {
            competitorName: competitorName,
            competitionName: competitionName,
            eventName: competeEvent?.displayname
        }}
    )
    const data = response.data;
    return data;
}

export const initialCompetitionState: CompetitionState = {
    id: "",
    name: "",
    startdate: new Date().toISOString(),
    enddate: new Date().toISOString(),
    events: [],
    currentEventIdx: 0,
    noOfSolves: 1,
    currentSolveIdx: 0,
    scrambles: [],
    inputMethod: InputMethod.Manual,
    loadingState: {
        results: false,
        compinfo: false,
        error: ''
    },
    results: {
        id: 0,
        userid: 0,
        username: '',
        competitionid: 0,
        competitionname: '',
        eventid: 0,
        eventname: '',
        iconcode: '',
        format: '',
        solve1: '',
        solve2: '',
        solve3: '',
        solve4: '',
        solve5: '',
        comment: '',
        status: {
            id: 0,
            approvalFinished: true,
            visible: true,
            displayname: '',
        }
    },
    penalties: Array(5).fill('0')
};

export const reformatTime = (oldFormattedTime: string, added: boolean = false): string => {
    if (added) {
        let idx = 0;
        while (idx < oldFormattedTime.length && /^\D/.test(oldFormattedTime[idx]) || oldFormattedTime[idx] === '0')
            idx++;
        oldFormattedTime = oldFormattedTime.slice(idx);
    }

    const matchedDigits = oldFormattedTime.match(/\d+/g);
    let digits = !matchedDigits ? '' : matchedDigits.join('');
    if (digits.length < 3)
        digits = digits.padStart(3, '0');

    let newFormattedTime = `${digits[digits.length - 1]}${digits[digits.length - 2]}.`;
    let idx = digits.length - 3;
    while (idx >= 0) {
        newFormattedTime += digits[idx--];
        if (idx >= 0)
            newFormattedTime += digits[idx--];
        if (idx >= 0)
            newFormattedTime += ':';
    }

    newFormattedTime = newFormattedTime.split('').reverse().join('');

    return newFormattedTime;
}

export const sendResults = async (resultEntry: ResultEntry) => {
    console.log('zatial sa nic neudeje', resultEntry);
}

export const saveValidation = async (resultEntry: ResultEntry, verdict: boolean) => {
    console.log('zatial sa nic neudeje', verdict, resultEntry);
}

export const competitionOnGoing = (state: CompetitionState): boolean => {
    const startdate = new Date(state.startdate)
    const now = new Date();
    const enddate = new Date(state.enddate)
    return startdate < now && now < enddate;
}