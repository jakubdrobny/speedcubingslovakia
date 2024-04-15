export type CompetitionData = {
  id: string;
  name: string;
  startdate: string;
  enddate: string;
  events: CompetitionEvent[];
  scrambles: ScrambleSet[];
};

export type ScrambleSet = {
  event: CompetitionEvent;
  scrambles: string[];
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
  puzzlecode: string;
};

export enum InputMethod {
  Manual,
  Timer,
}

export type CompetitionLoadingState = {
  compinfo: boolean;
  results: boolean;
  error: string;
};

export type CompetitionState = CompetitionData & {
  currentEventIdx: number;
  currentSolveIdx: number;
  noOfSolves: number;
  inputMethod: InputMethod;
  penalties: string[];
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
  setSuspicousModalOpen: (newOpen: boolean) => void;
  results: CompetitionResult[];
  setResults: (newResults: CompetitionResult[]) => void;
  resultsCompeteChoice: ResultsCompeteChoiceEnum;
  setResultsCompeteChoice: (newChoice: ResultsCompeteChoiceEnum) => void;
  loadingState: CompetitionLoadingState;
  setLoadingState: (newState: CompetitionLoadingState) => void;
  fetchCompetitionResults: (event?: CompetitionEvent, compId?: string) => void;
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
};

export type AuthContextType = {
  authState: AuthState;
  updateAuthToken: (newToken: string) => void;
  setAuthState: (newAuthState: AuthState) => void;
};

export type TimerInputContextType = {
  timerInputState: TimerInputState;
  handleTimerInputKeyDown: EventListener;
  handleTimerInputKeyUp: EventListener;
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
  username: string;
  country_name: string;
  country_iso2: string;
  single: string;
  average: string;
  times: string[];
  score: string;
};

export type NavContextType = {
  navOpen: boolean;
  openNav: () => void;
  closeNav: () => void;
  toggleNavOpen: () => void;
};
