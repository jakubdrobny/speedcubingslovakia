import { ReactElement } from "react";

export type CompetitionData = {
  id: string;
  name: string;
  startdate: string;
  enddate: string;
  events: CompetitionEvent[];
  scrambles: ScrambleSet[];
};

export enum Permission {
  ADMIN,
  USER,
}

export type Scramble = {
  scramble: string;
  img: string;
};

export type ScrambleSet = {
  event: CompetitionEvent;
  scrambles: Scramble[];
};

export enum FilterValue {
  Current = "Current",
  Past = "Past",
  Future = "Future",
}

export type CompetitionEvent = {
  id: number;
  displayname: string;
  format: string;
  iconcode: string;
  scramblingcode: string;
};

export enum InputMethod {
  Manual,
  Timer,
}

export type CompetitionLoadingState = {
  compinfo: boolean;
  results: boolean;
  error: ResponseError;
};

export type CompetitionState = CompetitionData & {
  currentEventIdx: number;
  currentSolveIdx: number;
  noOfSolves: number;
  inputMethod: InputMethod;
  penalties: string[];
};

export type ResponseError = {
  message?: string;
  element?: ReactElement;
};

export type CompetitionContextType = {
  competitionState: CompetitionState;
  currentResults: ResultEntry;
  updateBasicInfo: (info: CompetitionData) => void;
  updateCurrentEvent: (idx: number) => void;
  updateCurrentSolve: (idx: number) => void;
  saveResults: () => Promise<void>;
  updateSolve: (newTime: string) => void;
  toggleInputMethod: () => void;
  addPenalty: (newPenalty: string) => void;
  setCompetitionState: (newState: CompetitionState) => void;
  setCurrentResults: (newResults: ResultEntry) => void;
  suspicousModalOpen: boolean;
  warningModalOpen: boolean;
  setSuspicousModalOpen: (newOpen: boolean) => void;
  setWarningModalOpen: (newModal: boolean) => void;
  results: CompetitionResult[];
  setResults: (newResults: CompetitionResult[]) => void;
  anyComment: boolean;
  setAnyComment: (newAnyComment: boolean) => void;
  resultsCompeteChoice: ResultsCompeteChoiceEnum;
  setResultsCompeteChoice: (newChoice: ResultsCompeteChoiceEnum) => void;
  loadingState: CompetitionLoadingState;
  setLoadingState: (newState: CompetitionLoadingState) => void;
  fetchCompetitionResults: (event?: CompetitionEvent, compId?: string) => void;
  fetchCompeteResultEntry: (event?: CompetitionEvent, compId?: string) => void;
  competitionStateRef: { current: CompetitionState };
  currentResultsRef: { current: ResultEntry };
};

export type ResultEntry = {
  id: number;
  userid: number;
  username: string;
  competitionid: string;
  competitionname: string;
  eventid: number;
  eventname: string;
  iconcode: string;
  format: string;
  solve1: string;
  solve2: string;
  solve3: string;
  solve4: string;
  solve5: string;
  comment: string;
  status: ResultsStatus;
  badFormat: boolean;
};

export enum ResultEntrySolves {
  solve1,
  solve2,
  solve3,
  solve4,
  solve5,
}

export type AuthState = {
  token: string;
  wcaid: string;
  isadmin: boolean;
  avatarUrl: string;
  username: string;
};

export type AuthContextType = {
  authState: AuthState;
  updateAuthToken: (newToken: string) => void;
  setAuthState: (newAuthState: AuthState) => void;
  authStateRef: { current: AuthState };
};

export type WindowSize = {
  width: number;
  height: number;
};

export type WindowSizeContextType = {
  windowSize: WindowSize;
  setWindowSize: (newWindowSize: WindowSize) => void;
};

export type TimerInputContextType = {
  timerInputState: TimerInputState;
  handleTimerInputKeyDown: EventListener;
  handleTimerInputKeyUp: (
    e: Event,
    handleSaveResults?: (moveIndex: boolean) => void
  ) => void;
  timerRef: { current: HTMLDivElement | null };
};

export type TimerInputState = {
  currentState: TimerInputCurrentState;
  color: string;
};

export enum TimerInputCurrentState {
  NotSolving,
  GettingReady,
  Ready,
  Solving,
  Finishing,
}

export enum TimerColors {
  Default = "#32383E",
  Red = "red",
  Green = "green",
}

export enum PenaltyType {
  PlusTwo,
  DNF,
}

export type ManageRolesUser = {
  id: number;
  name: string;
  isadmin: boolean;
};

export enum DashboardPanel {
  ManageRoles,
  None,
}

export type CompetitionEditProps = {
  edit: boolean;
};

export type AnnouncementEditProps = {
  edit: boolean;
};

export type ResultsStatus = {
  id: number;
  approvalFinished: boolean;
  approved?: boolean;
  visible: boolean;
  displayname: string;
};

export enum ResultsCompeteChoiceEnum {
  Results,
  Compete,
}

export type CompetitionResult = {
  place: string;
  username: string;
  wca_id: string;
  country_name: string;
  country_iso2: string;
  single: string;
  average: string;
  times: string[];
  score: string;
  comment: string;
};

export type CompetitionResultStruct = {
  results: CompetitionResult[];
  anyComment: boolean;
};

export type NavContextType = {
  navOpen: boolean;
  openNav: () => void;
  closeNav: () => void;
  toggleNavOpen: () => void;
};

export type LoadingState = {
  isLoading: boolean;
  error: ResponseError;
};

export type ProfileTypeBasics = {
  name: string;
  imageurl: string;
  region: {
    name: string;
    iso2: string;
  };
  wcaid: string;
  sex: string;
  noOfCompetitions: number;
  completedSolves: number;
};

export type PersonalBestEntry = {
  nr: string;
  cr: string;
  wr: string;
  value: string;
};

export type ProfileTypePersonalBests = {
  eventName: string;
  eventIconcode: string;
  average: PersonalBestEntry;
  single: PersonalBestEntry;
};

export type ProfileTypeResultHistoryEntry = {
  competitionId: string;
  competitionName: string;
  place: string;
  single: string;
  singleRecord: string;
  singleRecordColor: string;
  average: string;
  averageRecord: string;
  averageRecordColor: string;
  solves: string[];
};

export type ProfileTypeResultHistory = {
  eventName: string;
  eventIconcode: string;
  eventFormat: string;
  history: ProfileTypeResultHistoryEntry[];
};

export type ProfileType = {
  basics: ProfileTypeBasics;
  personalBests: ProfileTypePersonalBests[];
  medalCollection: { gold: string; silver: string; bronze: string };
  recordCollection: { wr: string; cr: string; nr: string };
  resultsHistory: ProfileTypeResultHistory[];
};

export type SearchUser = {
  username: string;
  wcaid: string;
};

export type RegionSelectGroup = {
  groupName: string;
  groupMembers: string[];
};

export type RankingsEntry = {
  place: string;
  username: string;
  wca_id: string;
  country_iso2: string;
  country_name: string;
  result: string;
  competitionName: string;
  competitionId: string;
  times: string[];
};

export type RecordsItem = {
  eventname: string;
  iconcode: string;
  entries: RecordsItemEntry[];
};

export type RecordsItemEntry = {
  type: string;
  username: string;
  wcaId: string;
  result: string;
  countryIso2: string;
  countryName: string;
  competitionName: string;
  competitionId: string;
  solves: string[];
};

export type Tag = {
  label: string;
  color: "danger" | "warning" | "success" | "primary";
};

export type EmojiCounter = {
  emoji: string;
  by: string;
};

export type AnnouncementState = {
  id: number;
  authorId: number;
  authorWcaId: string;
  authorUsername: string;
  title: string;
  content: string;
  tags: Tag[];
  read: boolean;
  emojiCounters: EmojiCounter[];
};

export const initialAnnouncementState: AnnouncementState = {
  id: 0,
  authorId: 0,
  authorWcaId: "",
  authorUsername: "",
  title: "",
  content: "",
  tags: [],
  read: true,
  emojiCounters: [],
};

export type AnnouncementReactResponse = {
  set: boolean;
};

export type AverageInfo = {
  single: string;
  average: string;
  times: string[];
  bpa: string;
  wpa: string;
  showPossibleAverage: boolean;
  finishedCompeting: boolean;
  place: string;
  singleRecord: string;
  singleRecordColor: string;
  averageRecord: string;
  averageRecordColor: string;
};
