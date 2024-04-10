import {
  AuthState,
  CompetitionData,
  CompetitionEvent,
  CompetitionState,
  FilterValue,
  InputMethod,
  ManageRolesUser,
  ResultEntry,
  ScrambleSet,
} from "./Types";

import { Alg } from "cubing/alg";
import Cookies from "universal-cookie";
import axios from "axios";
import { randomScrambleForEvent } from "cubing/scramble";
import { setSearchDebug } from "cubing/search";

setSearchDebug({
  logPerf: false,
  scramblePrefetchLevel: "none",
});

export const loadFilteredCompetitions = async (
  filterValue: FilterValue
): Promise<CompetitionData[]> => {
  const response = await axios.get(`/api/competitions/filter/${filterValue}`);
  return response.data;
};

export const getCompetitionById = async (
  id: string | undefined
): Promise<CompetitionData> => {
  if (id === undefined) {
    return Promise.reject("Invalid competition id.");
  }

  const response = await axios.get(`/api/competitions/id/${id}`);
  return !response.data ? undefined : response.data;
};

export const getResultsFromCompetitionAndEvent = async (
  uid: number,
  cid: string | undefined,
  event: CompetitionEvent | undefined,
  token: string
): Promise<ResultEntry> => {
  if (cid === undefined || event === undefined)
    return Promise.reject("invalid competition/event id");
  const response = await axios.get(
    `/api/results/compete/${uid}/${cid}/${event.id}`,
    { headers: { Authorization: `Bearer ${token}`, UserId: uid } }
  );
  return response.data;
};

const formattedToMiliseconds = (formattedTime: string): number => {
  let res = 0;

  const formattedTimeSplit = formattedTime.split(".");
  const wholePart = formattedTimeSplit[0].split(":").reverse(),
    decimalPart = formattedTimeSplit[1];

  res += parseInt(decimalPart) * 10;
  if (wholePart.length > 0) res += parseInt(wholePart[0]) * 1000;
  if (wholePart.length > 1) res += 60 * parseInt(wholePart[1]) * 1000;
  if (wholePart.length > 2) res += 60 * 60 * parseInt(wholePart[2]) * 1000;
  if (wholePart.length > 3) res += 24 * 60 * 60 * parseInt(wholePart[3]) * 1000;

  return res;
};

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
    res.push(mul === 1 ? toPush.padStart(3, "0") : toPush);
    toFormat %= pw;
    pw = Math.floor(pw / mul);
  }

  res[res.length - 1] = res[res.length - 1].slice(
    0,
    res[res.length - 1].length - 1
  );
  let sliceIdx = 0;
  while (sliceIdx < res.length - 2 && res[sliceIdx] === "0") sliceIdx += 1;
  res = res.slice(sliceIdx);

  let resString = "";
  let resIdx: number;
  for (resIdx = 0; resIdx < res.length - 1; resIdx++) {
    resString += resIdx > 0 ? res[resIdx].padStart(2, "0") : res[resIdx];
    resString += resIdx == res.length - 2 ? "." : ":";
  }
  resString += res[resIdx].padStart(2, "0");

  return resString;
};

export const reformatWithPenalties = (
  oldFormattedTime: string,
  penalty: string
) => {
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
};

export const getManageUsers = async (): Promise<ManageRolesUser[]> => {
  const response = await axios.get("/api/users/manage-roles");
  return response.data;
};

export const updateUserRoles = async (
  newUsers: ManageRolesUser[]
): Promise<ManageRolesUser[]> => {
  const response = await axios.put("/api/users/manage-roles", newUsers);
  return response.data;
};

export const getAvailableEvents = async (): Promise<CompetitionEvent[]> => {
  const response = await axios.get("/api/events");
  return response.data;
};

const cancels = async (scramble: string) => {
  let cancelledScramble = new Alg(scramble)
    .experimentalSimplify({ cancel: true })
    .toString();
  return scramble !== cancelledScramble;
};

const generateScrambleSetsFromEvents = async (
  events: CompetitionEvent[]
): Promise<ScrambleSet[]> => {
  let scrambleSets: ScrambleSet[] = [];

  for (const event of events) {
    const match = event.format.match(/\d+$/)?.[0];
    const noOfSolves = match ? parseInt(match) : 1;

    let scrambles: string[] = [];
    for (let i = 0; i < noOfSolves; i++) {
      let scramble: string = (
        await randomScrambleForEvent(event.iconcode)
      ).toString();
      while (event.displayname === "FMC" && (await cancels(scramble))) {
        scramble = (await randomScrambleForEvent(event.iconcode)).toString();
      }
      scrambles.push(scramble);
    }

    scrambleSets.push({ event, scrambles });
  }

  return scrambleSets;
};

export const updateCompetition = async (
  state: CompetitionState,
  edit: boolean
): Promise<CompetitionState> => {
  const reqBody: CompetitionData = {
    id: state.id,
    name: state.name,
    startdate: state.startdate.endsWith("Z")
      ? state.startdate
      : state.startdate + ":00Z",
    enddate: state.enddate.endsWith("Z")
      ? state.enddate
      : state.enddate + ":00Z",
    events: state.events.toSorted(
      (e1: CompetitionEvent, e2: CompetitionEvent) => e1.id - e2.id
    ),
    scrambles: edit ? [] : await generateScrambleSetsFromEvents(state.events),
  };

  const response = await axios({
    method: edit ? "PUT" : "POST",
    url: "/api/competitions",
    data: reqBody,
  });

  return response.data;
};

export const getResults = async (
  username: string,
  cid: string,
  competeEvent: CompetitionEvent | undefined
) => {
  if (competeEvent === undefined) return Promise.reject("proste nie");

  username = username === "" ? "_" : username;
  cid = cid === "" ? "_" : cid;

  const response = await axios.get(
    `/api/results/edit/${username}/${cid}/${competeEvent.id}`
  );
  return response.data;
};

export const formatCompetitionDateForInput = (originalDate: string): string => {
  const originalDateSplit = originalDate.split(".")[0].split(":");
  return originalDateSplit.slice(0, 2).join(":");
};

export const initialCompetitionState: CompetitionState = {
  id: "",
  name: "",
  startdate: formatCompetitionDateForInput(new Date().toISOString()),
  enddate: formatCompetitionDateForInput(new Date().toISOString()),
  events: [],
  currentEventIdx: 0,
  noOfSolves: 1,
  currentSolveIdx: 0,
  scrambles: [],
  inputMethod: InputMethod.Manual,
  loadingState: {
    compinfo: false,
    results: false,
    error: "",
  },
  penalties: Array(5).fill("0"),
};

export const initialCurrentResults: ResultEntry = {
  id: 0,
  userid: 0,
  username: "",
  competitionid: "",
  competitionname: "",
  eventid: 0,
  eventname: "",
  iconcode: "",
  format: "",
  solve1: "",
  solve2: "",
  solve3: "",
  solve4: "",
  solve5: "",
  comment: "",
  status: {
    id: 0,
    approvalFinished: true,
    visible: true,
    displayname: "",
  },
};

export const reformatTime = (
  oldFormattedTime: string,
  added: boolean = false
): string => {
  if (added) {
    let idx = 0;
    while (
      (idx < oldFormattedTime.length && /^\D/.test(oldFormattedTime[idx])) ||
      oldFormattedTime[idx] === "0"
    )
      idx++;
    oldFormattedTime = oldFormattedTime.slice(idx);
  }

  const matchedDigits = oldFormattedTime.match(/\d+/g);
  let digits = !matchedDigits ? "" : matchedDigits.join("");
  if (digits.length < 3) digits = digits.padStart(3, "0");

  let newFormattedTime = `${digits[digits.length - 1]}${
    digits[digits.length - 2]
  }.`;
  let idx = digits.length - 3;
  while (idx >= 0) {
    newFormattedTime += digits[idx--];
    if (idx >= 0) newFormattedTime += digits[idx--];
    if (idx >= 0) newFormattedTime += ":";
  }

  newFormattedTime = newFormattedTime.split("").reverse().join("");

  return newFormattedTime;
};

export const sendResults = async (
  resultEntry: ResultEntry
): Promise<ResultEntry> => {
  const response = await axios.post("/api/results/save", resultEntry);
  return response.data;
};

export const saveValidation = async (
  resultEntry: ResultEntry,
  verdict: boolean
) => {
  const response = await axios.post("/api/results/save-validation", {
    resultId: resultEntry.id,
    verdict,
  });
  return response.data;
};

export const competitionOnGoing = (state: CompetitionState): boolean => {
  const startdate = new Date(state.startdate);
  const now = new Date();
  const enddate = new Date(state.enddate);
  return startdate < now && now < enddate;
};

export const formatDate = (dateString: string): String => {
  const date = new Date(dateString);
  return date.toLocaleDateString() + " " + date.toLocaleTimeString();
};

export const logIn = async (
  searchParams: URLSearchParams
): Promise<AuthState> => {
  const code = searchParams.get("code");
  if (code === null) {
    return Promise.reject("Missing code.");
  }

  const response = await axios.post("/api/login", code);
  const data: {
    access_token: string;
    expires_in: number;
    userid: number;
    isadmin: boolean;
    avatarUrl: string;
    wcaid: string;
  } = response.data;

  const result: AuthState = {
    token: data.access_token,
    userid: data.userid,
    authenticated: true,
    admin: data.isadmin,
    avatarUrl: data.avatarUrl,
    wcaid: data.wcaid,
  };

  const cookies = new Cookies(null, { path: "/" });

  let key: keyof AuthState;
  for (key in result) {
    cookies.set(key, result[key], {
      expires: new Date(new Date().getTime() + data.expires_in * 60 * 1000),
    });
  }

  return result;
};

const cookies = new Cookies(null, { path: "/" });

export const initialAuthState: AuthState = {
  token: cookies.get("token") || "",
  userid: parseInt(cookies.get("userid")) || -1,
  authenticated: Boolean(cookies.get("authenticated")) || false,
  admin: Boolean(cookies.get("admin")) || false,
  avatarUrl: cookies.get("avatarUrl") || "",
  wcaid: cookies.get("wcaid") || "",
};

export const logOut = () => {
  let key: keyof AuthState;
  for (key in initialAuthState) {
    cookies.remove(key);
  }
};
