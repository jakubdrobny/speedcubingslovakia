import {
  CompetitionContextType,
  CompetitionData,
  CompetitionEvent,
  CompetitionLoadingState,
  CompetitionResult,
  CompetitionState,
  InputMethod,
  ResultEntry,
  ResultsCompeteChoiceEnum,
} from "../../Types";
import React, { ReactNode, createContext } from "react";
import {
  competitionOnGoing,
  getCompetitionResults,
  getError,
  getResultsFromCompetitionAndEvent,
  initialCompetitionLoadingState,
  initialCompetitionState,
  initialCurrentResults,
  sendResults,
} from "../../utils";

import useState from "react-usestateref";

export const CompetitionContext = createContext<CompetitionContextType | null>(
  null
);

export const CompetitionProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [competitionState, setCompetitionState, competitionStateRef] =
    useState<CompetitionState>(initialCompetitionState);
  const [currentResults, setCurrentResults, currentResultsRef] =
    useState<ResultEntry>(initialCurrentResults);
  const [suspicousModalOpen, setSuspicousModalOpen] = useState<boolean>(false);
  const [warningModalOpen, setWarningModalOpen] = useState<boolean>(false);
  const [resultsCompeteChoice, setResultsCompeteChoice] =
    useState<ResultsCompeteChoiceEnum>(ResultsCompeteChoiceEnum.Results);
  const [results, setResults] = useState<CompetitionResult[]>([]);
  const [loadingState, setLoadingState] = useState<CompetitionLoadingState>(
    initialCompetitionLoadingState
  );

  const updateBasicInfo = (info: CompetitionData) => {
    const match =
      info.events[competitionStateRef.current.currentEventIdx].format.match(
        /\d+$/
      )?.[0];
    const noOfSolves = match ? parseInt(match) : 1;

    setCompetitionState((ps) => {
      return {
        ...ps,
        ...info,
        noOfSolves: noOfSolves,
        currentEventIdx: 0,
        currentSolveIdx: 0,
      };
    });

    setLoadingState({ ...loadingState, compinfo: false, error: {} });
  };

  const fetchCompeteResultEntry = (
    event: CompetitionEvent = competitionStateRef.current.events[
      competitionStateRef.current.currentEventIdx
    ],
    compId: string = competitionState.id
  ) => {
    if (event.displayname === "Overall") {
      event =
        competitionStateRef.current.events[
          competitionStateRef.current.currentEventIdx - 1
        ];
      setCompetitionState((ps) => ({
        ...ps,
        currentEventIdx: ps.currentEventIdx - 1,
      }));
    }

    setLoadingState((ps) => ({ ...ps, results: true, error: {} }));
    getResultsFromCompetitionAndEvent(compId, event)
      .then((resultEntry) => {
        setCurrentResults(resultEntry);
        if (!resultEntry.status.approvalFinished) {
          setSuspicousModalOpen(true);
        } else if (resultEntry.badFormat) {
          setWarningModalOpen(true);
        }
        setLoadingState((ps) => ({
          ...ps,
          results: false,
        }));
      })
      .catch((err) => {
        setLoadingState((ps) => ({
          ...ps,
          results: false,
          error: getError(err),
        }));
      });
  };

  const fetchCompetitionResults = (
    event: CompetitionEvent = competitionStateRef.current.events[
      competitionStateRef.current.currentEventIdx
    ],
    compId: string = competitionState.id
  ) => {
    setLoadingState((ps) => ({ ...ps, results: true, error: {} }));
    getCompetitionResults(compId, event)
      .then((res) => {
        setResults(res);
        setLoadingState((ps) => ({
          ...ps,
          results: false,
        }));
      })
      .catch((err) => {
        setLoadingState((ps) => ({
          ...ps,
          results: false,
          error: getError(err),
        }));
      });
  };

  const updateCurrentEvent = async (idx: number) => {
    const match = competitionState.events[idx].format.match(/\d+$/)?.[0];
    const noOfSolves = match ? parseInt(match) : 1;
    setCompetitionState((ps) => ({
      ...ps,
      currentEventIdx: idx,
      noOfSolves: noOfSolves,
      currentSolveIdx: 0,
      penalties: Array(5).fill("0"),
      inputMethod:
        competitionState.events[idx].displayname === "FMC"
          ? InputMethod.Manual
          : ps.inputMethod,
    }));
    setLoadingState((ps) => ({ ...ps, results: true, error: {} }));

    if (resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete)
      fetchCompeteResultEntry();
    else fetchCompetitionResults(competitionStateRef.current.events[idx]);
  };

  const updateCurrentSolve = (idx: number) =>
    setCompetitionState({ ...competitionState, currentSolveIdx: idx });

  const saveResults = async (): Promise<void> => {
    try {
      const resultEntry = await sendResults({
        ...currentResultsRef.current,
        competitionid: competitionState.id,
      });
      if (!resultEntry.status.approvalFinished) {
        setSuspicousModalOpen(true);
      } else if (resultEntry.badFormat) {
        setWarningModalOpen(true);
      }
      setCurrentResults(resultEntry);
      return Promise.resolve();
    } catch (e) {
      return Promise.reject(e);
    }
  };

  const addPenalty = (newPenalty: string) => {
    const curSolveIdx = competitionState.currentSolveIdx;
    const oldPenalty = competitionState.penalties[curSolveIdx];

    if (newPenalty === "DNF") {
      newPenalty = oldPenalty === "DNF" ? "0" : "DNF";
    } else {
      newPenalty =
        oldPenalty === "DNF"
          ? "2"
          : ((parseInt(oldPenalty) + parseInt(newPenalty)) % 18).toString();
    }

    setCompetitionState((ps) => ({
      ...ps,
      penalties: ps.penalties.map((val: string, idx: number) =>
        idx === curSolveIdx ? newPenalty : val
      ),
    }));
  };

  const updateSolve = (newTime: string) => {
    const solveProp: keyof ResultEntry = `solve${
      competitionStateRef.current.currentSolveIdx + 1
    }` as keyof ResultEntry;
    setCurrentResults((ps) => ({
      ...ps,
      [solveProp]: newTime,
    }));
  };

  const toggleInputMethod = () => {
    if (!competitionOnGoing(competitionState)) return;

    if (
      competitionState.currentEventIdx < competitionState.events.length &&
      competitionState.events[competitionState.currentEventIdx].displayname !==
        "FMC" &&
      competitionState.events[competitionState.currentEventIdx].displayname !==
        "MBLD"
    ) {
      setCompetitionState((ps) => ({
        ...ps,
        inputMethod:
          ps.inputMethod === InputMethod.Manual
            ? InputMethod.Timer
            : InputMethod.Manual,
        penalties: ps.penalties.map((val: string, idx: number) =>
          idx === ps.currentSolveIdx ? "0" : val
        ),
      }));
    }
  };

  return (
    <CompetitionContext.Provider
      value={{
        competitionState,
        updateBasicInfo,
        updateCurrentEvent,
        updateCurrentSolve,
        saveResults,
        toggleInputMethod,
        addPenalty,
        updateSolve,
        setCompetitionState,
        currentResults,
        setCurrentResults,
        suspicousModalOpen,
        setSuspicousModalOpen,
        warningModalOpen,
        setWarningModalOpen,
        results,
        setResults,
        resultsCompeteChoice,
        setResultsCompeteChoice,
        loadingState,
        setLoadingState,
        fetchCompetitionResults,
        fetchCompeteResultEntry,
        competitionStateRef,
        currentResultsRef,
      }}
    >
      {children}
    </CompetitionContext.Provider>
  );
};
