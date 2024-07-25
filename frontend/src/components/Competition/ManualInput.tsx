import { Alert, Button, Input } from "@mui/joy";
import {
  CompetitionContextType,
  ResponseError,
  ResultEntry,
} from "../../Types";
import {
  competitionOnGoing,
  getError,
  reformatTime,
  renderResponseError,
} from "../../utils";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { MAX_MANUAL_INPUT_LENGTH } from "../../constants";

const ManualInput = () => {
  const [forceRerender, setForceRerender] = useState(false);
  const { competitionState, updateSolve, saveResults, currentResultsRef } =
    useContext(CompetitionContext) as CompetitionContextType;
  const solveProp: keyof ResultEntry = `solve${
    competitionState.currentSolveIdx + 1
  }` as keyof ResultEntry;
  const formattedTime = currentResultsRef.current[solveProp].toString();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<ResponseError>();

  useEffect(
    () => setForceRerender(!forceRerender),
    [competitionState.currentEventIdx]
  );

  const handleTimeInputChange = (e: React.FormEvent<HTMLInputElement>) => {
    const newValue = e.currentTarget.value;

    if (
      competitionState.events[competitionState.currentEventIdx].displayname ===
      "FMC"
    ) {
      updateSolve(newValue);
      return;
    }

    if (newValue.length > MAX_MANUAL_INPUT_LENGTH) return;

    // character deleted
    if (newValue.length + 1 === formattedTime.length) {
      if (newValue.endsWith("N")) {
        updateSolve("0.00");
        return;
      } else {
        updateSolve(reformatTime(newValue));
      }
    } else {
      if (newValue.endsWith("d")) {
        updateSolve("DNF");
      } else if (newValue.endsWith("s")) {
        updateSolve("DNS");
      } else if (/\d$/.test(newValue.slice(-1))) {
        updateSolve(reformatTime(newValue, true));
      } else {
        updateSolve("DNF");
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      setIsLoading(true);
      saveResults()
        .then(() => setIsLoading(false))
        .catch((err) => {
          setIsLoading(false);
          setError(getError(err));
        });
    }
  };

  return (
    <div>
      {error && renderResponseError(error)}
      <Input
        size="lg"
        placeholder="Enter your time or solution..."
        sx={{ marginBottom: 2, marginTop: 2 }}
        value={formattedTime}
        onChange={handleTimeInputChange}
        onKeyDown={handleKeyDown}
        disabled={!competitionOnGoing(competitionState) || isLoading}
      />
    </div>
  );
};

export default ManualInput;
