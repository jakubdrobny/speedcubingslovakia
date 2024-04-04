export type CompetitionData = {
  id: string;
  name: string;
  startdate: string;
  enddate: string;
  events: CompetitionEvent[];
  scrambles: ScrambleSet[];
  results?: ResultEntry;
};

export type ScrambleSet = {
  event: CompetitionEvent;
  scrambles: string[];
};

export type CompetitionDBModel = {
  id: string;
  name: string;
  startdate: string;
  enddate: string;
  events: CompetitionEvent[];
  scrambles: ScrambleSet[];
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

export type CompetitionState = CompetitionData & {
  currentEventIdx: number;
  currentSolveIdx: number;
  noOfSolves: number;
  inputMethod: InputMethod;
  loadingState: {
    compinfo: boolean;
    error: string;
  };
  penalties: string[];
};

export type CompetitionContextType = {
  competitionState: CompetitionState;
  currentResults: ResultEntry;
  updateBasicInfo: (info: CompetitionData) => void;
  updateCurrentEvent: (idx: number) => void;
  updateCurrentSolve: (idx: number) => void;
  saveResults: () => void;
  updateSolve: (newTime: string) => void;
  toggleInputMethod: () => void;
  addPenalty: (newPenalty: string) => void;
  updateCompetitionName: (newName: string) => void;
  updateCompetitionEvents: (newEvents: CompetitionEvent[]) => void;
  setCompetitionState: (newState: CompetitionState) => void;
  setCurrentResults: (newResults: ResultEntry) => void;
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
  authenticated: boolean;
  admin: boolean;
};

export type AuthContextType = {
  authState: AuthState;
  updateAuthToken: (newToken: string) => void;
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
