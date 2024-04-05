import {
  AuthContextType,
  CompetitionContextType,
  CompetitionData,
  CompetitionState,
  InputMethod,
  ResultEntry,
} from "../../Types";
import React, { ReactNode, createContext, useContext, useState } from "react";
import {
  competitionOnGoing,
  getResultsFromCompetitionAndEvent,
  initialCompetitionState,
  initialCurrentResults,
  sendResults,
} from "../../utils";

import { AuthContext } from "../../context/AuthContext";

export const CompetitionContext = createContext<CompetitionContextType | null>(
  null
);

export const CompetitionProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [competitionState, setCompetitionState] = useState<CompetitionState>(
    initialCompetitionState
  );
  const [currentResults, setCurrentResults] = useState<ResultEntry>(
    initialCurrentResults
  );
  const { authState } = useContext(AuthContext) as AuthContextType;
  const [suspicousModalOpen, setSuspicousModalOpen] = useState<boolean>(false);

  const updateBasicInfo = (info: CompetitionData) => {
    const match =
      info.events[competitionState.currentEventIdx].format.match(/\d+$/)?.[0];
    const noOfSolves = match ? parseInt(match) : 1;
    setCompetitionState((ps) => {
      return {
        ...ps,
        ...info,
        noOfSolves: noOfSolves,
        currentEventIdx: 0,
        currentSolveIdx: 0,
        loadingState: {
          ...ps.loadingState,
          results: false,
        },
      };
    });
  };

  const updateCurrentEvent = async (idx: number) => {
    const match = competitionState.events[idx].format.match(/\d+$/)?.[0];
    const noOfSolves = match ? parseInt(match) : 1;
    setCompetitionState((ps) => ({
      ...ps,
      loadingState: { ...ps.loadingState, results: true },
    }));
    getResultsFromCompetitionAndEvent(
      authState.userid,
      competitionState.id,
      competitionState.events[idx]
    )
      .then((resultEntry) => {
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
          loadingState: { ...ps.loadingState, results: false },
        }));
        setCurrentResults(resultEntry);
        if (!resultEntry.status.approvalFinished) setSuspicousModalOpen(true);
      })
      .catch((err) =>
        setCompetitionState((ps) => ({
          ...ps,
          loadingState: {
            ...ps.loadingState,
            results: false,
            error: err.message,
          },
        }))
      );
  };

  const updateCurrentSolve = (idx: number) =>
    setCompetitionState({ ...competitionState, currentSolveIdx: idx });

  const saveResults = async (): Promise<void> => {
    try {
      const resultEntry = await sendResults({
        ...currentResults,
        competitionid: competitionState.id,
      });
      if (!resultEntry.status.approvalFinished) setSuspicousModalOpen(true);
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
        idx == curSolveIdx ? newPenalty : val
      ),
    }));
  };

  const updateSolve = (newTime: string) => {
    const solveProp: keyof ResultEntry = `solve${
      competitionState.currentSolveIdx + 1
    }` as keyof ResultEntry;
    setCurrentResults({
      ...currentResults,
      [solveProp]: newTime,
    });
  };

  const toggleInputMethod = () => {
    if (!competitionOnGoing(competitionState)) return;

    if (
      competitionState.currentEventIdx < competitionState.events.length &&
      competitionState.events[competitionState.currentEventIdx].displayname !==
        "FMC"
    ) {
      setCompetitionState((ps) => ({
        ...ps,
        inputMethod:
          ps.inputMethod === InputMethod.Manual
            ? InputMethod.Timer
            : InputMethod.Manual,
        penalties: ps.penalties.map((val: string, idx: number) =>
          idx == ps.currentSolveIdx ? "0" : val
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
      }}
    >
      {children}
    </CompetitionContext.Provider>
  );
};
