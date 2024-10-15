import {
  AnnouncementReactResponse,
  AnnouncementState,
  AuthState,
  AverageInfo,
  CompetitionData,
  CompetitionEvent,
  CompetitionLoadingState,
  CompetitionResultStruct,
  CompetitionState,
  FilterValue,
  InputMethod,
  LoadingState,
  ManageRolesUser,
  ProfileType,
  RankingsEntry,
  RecordsItem,
  RegionSelectGroup,
  ResponseError,
  ResultEntry,
  ResultsStatus,
  SearchUser,
  Tag,
  WCACompetitionType
} from "../Types";
import { FeatureCollection, GeoJsonObject } from "geojson";
import axios, { AxiosError } from "axios";
import { useEffect, useState } from "react";

import { Alert } from "@mui/joy";
import Cookies from "universal-cookie";
import { Link } from "react-router-dom";

export const loadFilteredCompetitions = async (
  filterValue: FilterValue,
): Promise<CompetitionData[]> => {
  const response = await axios.get(`/api/competitions/filter/${filterValue}`);
  return response.data;
};

export const getCompetitionById = async (
  id: string | undefined,
): Promise<CompetitionData> => {
  if (id === undefined) {
    return Promise.reject("Invalid competition id.");
  }

  const response = await axios.get(`/api/competitions/id/${id}`);
  return !response.data ? undefined : response.data;
};

export const getResultsFromCompetitionAndEvent = async (
  cid: string | undefined,
  event: CompetitionEvent | undefined,
): Promise<ResultEntry> => {
  if (cid === undefined || event === undefined)
    return Promise.reject("invalid competition/event id");
  const response = await axios.get(`/api/results/compete/${cid}/${event.id}`);
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
    res[res.length - 1].length - 1,
  );
  let sliceIdx = 0;
  while (sliceIdx < res.length - 2 && res[sliceIdx] === "0") sliceIdx += 1;
  res = res.slice(sliceIdx);

  let resString = "";
  let resIdx: number;
  for (resIdx = 0; resIdx < res.length - 1; resIdx++) {
    resString += resIdx > 0 ? res[resIdx].padStart(2, "0") : res[resIdx];
    resString += resIdx === res.length - 2 ? "." : ":";
  }
  resString += res[resIdx].padStart(2, "0");

  return resString;
};

export const reformatWithPenalties = (
  oldFormattedTime: string,
  penalty: string,
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
  newUsers: ManageRolesUser[],
): Promise<ManageRolesUser[]> => {
  const response = await axios.put("/api/users/manage-roles", newUsers);
  return response.data;
};

export const getAvailableEvents = async (): Promise<CompetitionEvent[]> => {
  const response = await axios.get("/api/events");
  return response.data;
};

export const getAvailableResultsStatuses = async (): Promise<
  ResultsStatus[]
> => {
  const response = await axios.get("/api/resultsStatuses");
  return response.data;
};

export const updateCompetition = async (
  state: CompetitionState,
  edit: boolean,
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
    events: state.events.sort(
      (e1: CompetitionEvent, e2: CompetitionEvent) => e1.id - e2.id,
    ),
    scrambles: [],
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
  competeEvent: CompetitionEvent | undefined,
  resultsStatus: string,
) => {
  if (competeEvent === undefined || !resultsStatus)
    return Promise.reject("proste nie");

  username = username === "" ? "_" : username;
  cid = cid === "" ? "_" : cid;

  const response = await axios.get(
    `/api/results/edit/${username}/${cid}/${competeEvent.id}/${resultsStatus}`,
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
  penalties: Array(5).fill("0"),
};

export const isObjectEmpty = (obj: Object) => {
  return Object.keys(obj).length === 0;
};

export const renderResponseError = (error: ResponseError) => {
  if (error.message) return <Alert color="danger">{error.message}</Alert>;
  return error.element;
};

export const initialCurrentResults: ResultEntry = {
  id: 0,
  userid: 0,
  username: "",
  wcaid: "",
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
  badFormat: false,
};

export const reformatTime = (
  oldFormattedTime: string,
  added: boolean = false,
): string => {
  oldFormattedTime = sanitazeFromBrackets(oldFormattedTime);

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

  let newFormattedTime = `${digits[digits.length - 1]}${digits[digits.length - 2]
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
  resultEntry: ResultEntry,
): Promise<ResultEntry> => {
  const response = await axios.post("/api/results/save", resultEntry);
  return response.data;
};

export const saveValidation = async (
  resultEntry: ResultEntry,
  verdict: boolean,
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
  searchParams: URLSearchParams,
): Promise<AuthState> => {
  const code = searchParams.get("code");
  if (code === null) {
    return Promise.reject("Missing code.");
  }

  const response = await axios.post("/api/users/login", code);
  const data: {
    access_token: string;
    expires_in: number;
    isadmin: boolean;
    avatarUrl: string;
    wcaid: string;
    username: string;
  } = response.data;

  setBearerIfPresent(data.access_token);

  const result: AuthState = {
    token: data.access_token,
    isadmin: data.isadmin,
    avatarUrl: data.avatarUrl,
    wcaid: data.wcaid,
    username: data.username,
  };

  const cookies = new Cookies(null, { path: "/" });

  let key: keyof AuthState;
  for (key in result) {
    if (key === "isadmin") continue;
    cookies.set(key, result[key], {
      expires: new Date(new Date().getTime() + data.expires_in * 1000),
    });
  }

  return result;
};

const cookies = new Cookies(null, { path: "/" });

export const initialAuthState: AuthState = {
  token: cookies.get("token") || "",
  isadmin: false,
  avatarUrl: cookies.get("avatarUrl") || "",
  wcaid: cookies.get("wcaid") || "",
  username: cookies.get("username") || "",
};

export const logOut = async () => {
  let key: keyof AuthState;
  for (key in initialAuthState) {
    cookies.remove(key);
  }
};

export const authorizeAdmin = async () => {
  return axios.get("/api/users/auth/admin");
};

export const setBearerIfPresent = (token: string) => {
  axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
};

export const getCompetitionResults = async (
  competitionId: string,
  event: CompetitionEvent,
): Promise<CompetitionResultStruct> => {
  const response = await axios.get(
    `/api/competitions/results/${competitionId}/${event.id}`,
  );
  return response.data;
};

export const initialCompetitionLoadingState: CompetitionLoadingState = {
  results: false,
  compinfo: false,
  error: {},
};

export const emptyEvent: CompetitionEvent = {
  id: -1,
  displayname: "",
  format: "",
  iconcode: "",
  scramblingcode: "",
};

export const initialLoadingState: LoadingState = {
  isLoading: false,
  error: {},
};

export const getProfile = async (id: string): Promise<ProfileType> => {
  const response = await axios.get(`/api/results/profile/${id}`);
  return response.data;
};

export const defaultProfile: ProfileType = {
  basics: {
    name: "",
    imageurl: "",
    region: {
      name: "",
      iso2: "",
    },
    wcaid: "",
    sex: "",
    noOfCompetitions: 0,
    completedSolves: 0,
  },
  personalBests: [],
  medalCollection: { gold: "", silver: "", bronze: "" },
  recordCollection: { wr: "", cr: "", nr: "" },
  resultsHistory: [],
};

export const getError = (err: AxiosError): ResponseError => {
  const status = err.response?.status;
  if (status === 401 || status === 200) {
    return {
      element:
        status === 200 ? (
          <Alert color="success">
            <>{err.response?.data}</>
          </Alert>
        ) : (
          <Alert color="danger" sx={{ gap: 0 }}>
            Unauthorized/token expired. Try to{" "}
            <span style={{ padding: "0 2px" }}></span>
            <Link
              to={process.env.REACT_APP_WCA_GET_CODE_URL || ""}
              onClick={() => saveCurrentLocation(window.location.pathname)}
            >
              re-login
            </Link>
            .
          </Alert>
        ),
    };
  }
  return { message: err.response?.data as string };
};

export const getUsers = async (searchQuery: string): Promise<SearchUser[]> => {
  if (searchQuery === "") searchQuery = "_";
  const response = await axios.get(`/api/users/search?query=${searchQuery}`);
  return response.data;
};

export const getRegionGroups = async (): Promise<RegionSelectGroup[]> => {
  const response = await axios.get("/api/results/regions/grouped");
  return response.data;
};

export const getRankings = async (
  eid: number,
  single: boolean,
  regionGroup: string,
  region: string,
  queryType: string,
): Promise<RankingsEntry[]> => {
  const response = await axios.get(
    `/api/results/rankings?eid=${eid}&type=${single ? "single" : "average"
    }&regionGroup=${regionGroup}&region=${region}&queryType=${queryType.split("+")[1]
    }&numOfEntries=${queryType.split("+")[0]}`,
  );
  return response.data;
};

export const getRecords = async (
  eid: number,
  regionGroup: string,
  region: string,
): Promise<RecordsItem[]> => {
  const response = await axios.get(
    `/api/results/records?eid=${eid}&regionGroup=${regionGroup}&region=${region}`,
  );
  return response.data;
};

export const sanitazeFromBrackets = (oldTime: string): string => {
  if (oldTime.length == 0) return oldTime;

  let newTime = oldTime;
  if (newTime.length > 0 && newTime[0] == "(") newTime = newTime.slice(1);
  if (newTime.length > 0 && newTime[newTime.length - 1] == ")")
    newTime = newTime.slice(0, -1);

  return newTime;
};

export const reformatMultiTime = (startingTime: string): string => {
  startingTime = sanitazeFromBrackets(startingTime);

  if (startingTime === "DNS" || startingTime === "DNF") return startingTime;
  if (startingTime.indexOf(":") === -1) return startingTime;

  const cubePart = startingTime.split(" ")[0];
  let res = startingTime.split(" ")[1];

  while (res[0] === "0" || res[0] === ":") res = res.slice(1);

  const len = res.split(":").length;
  if (len === 1) {
    res = "00:" + res;
  } else if (len === 2) {
    res = res.padStart(5, "0");
  }

  return cubePart + " " + res;
};

export const reformatFMCSolve = (solve: string): string => {
  return solve.split(".")[0];
};

export const saveCurrentLocation = (locationPathname: string) => {
  const cookies = new Cookies(null, { path: "/" });
  cookies.set("backlink", locationPathname, {
    expires: new Date(new Date().getTime() + 60 * 1000), // expires in 1 minute
  });
};

export const initialResultsStruct: CompetitionResultStruct = {
  results: [],
  anyComment: false,
};

export const getCubingIconClassName = (iconcode: any): string => {
  return `cubing-icon ${iconcode.toString().startsWith("unofficial") ? "" : "event-"
    }${iconcode.toString()}`;
};

export const getAnnouncementById = async (
  id: string | undefined,
): Promise<AnnouncementState> => {
  if (id === undefined) {
    return Promise.reject("Invalid competition id.");
  }

  const response = await axios.get(`/api/announcements/id/${id}`);
  return !response.data ? undefined : response.data;
};

export const getAvailableTags = async (): Promise<Tag[]> => {
  const response = await axios.get("/api/tags");
  return response.data;
};

export const updateAnnoncement = async (
  state: AnnouncementState,
  edit: boolean,
): Promise<AnnouncementState> => {
  const response = await axios({
    method: edit ? "PUT" : "POST",
    url: "/api/announcements",
    data: state,
  });

  return response.data;
};

export const ReadAnnouncement = async (
  state: AnnouncementState,
): Promise<AnnouncementState> => {
  const response = await axios.get(`/api/announcements/read/${state.id}`);
  return response.data;
};

export const GetAnnouncements = async (): Promise<AnnouncementState[]> => {
  const response = await axios.get(`/api/announcements`);
  return response.data;
};

export const GetNoOfNewAnnouncements = async (): Promise<number> => {
  const response = await axios.get(`/api/announcements/noOfNew`);
  return response.data;
};

export const DeleteAnnouncement = async (
  announcementId: number,
): Promise<string> => {
  const response = await axios.delete(
    `/api/announcements/delete/${announcementId}`,
  );
  return response.data;
};

export const AddReactionToAnnouncement = async (
  announcementId: number,
  emoji: string,
  by: string,
): Promise<AnnouncementReactResponse> => {
  const response = await axios({
    method: "POST",
    url: `/api/announcements/react/${announcementId}`,
    data: { announcementId, emoji, by },
  });
  return response.data;
};

export const shouldHideAverageColumn = (
  eventFormat: string,
  eventIconcode: string,
): boolean => {
  return eventFormat == "bo1" || eventIconcode === "333mbf";
};

export const initialAverageInfo: AverageInfo = {
  single: "",
  average: "",
  times: [],
  bpa: "",
  wpa: "",
  showPossibleAverage: false,
  finishedCompeting: false,
  place: "",
  singleRecord: "",
  singleRecordColor: "",
  averageRecord: "",
  averageRecordColor: "",
};

export const GetAverageInfo = async (
  resultEntry: ResultEntry,
): Promise<AverageInfo> => {
  const response = await axios({
    method: "POST",
    url: `/api/results/averageinfo`,
    data: resultEntry,
  });
  return response.data;
};

export const GetAverageInfoRecords = async (
  resultEntry: ResultEntry,
  averageInfo: AverageInfo,
): Promise<AverageInfo> => {
  const response = await axios({
    method: "POST",
    url: `/api/results/averageinfo/records`,
    data: { resultEntry, averageInfo },
  });
  return response.data;
};

export const isBrowser = typeof window !== "undefined";

export const GetMapData = async (): Promise<FeatureCollection> => {
  const response = await axios.get(`/api/users/map`);
  return response.data;
};

export const GetWCACompetitions = async (region: string): Promise<WCACompetitionType[]> => {
  const response = await axios.get(`/api/competitions/wca?region=${region}`);
  return response.data;
} 
