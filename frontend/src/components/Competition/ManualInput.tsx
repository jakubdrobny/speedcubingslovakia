import { CompetitionContextType, ResultEntry } from "../../Types";
import React, { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { Input } from "@mui/joy";
import { MAX_MANUAL_INPUT_LENGTH } from "../../constants";
import { reformatTime } from "../../utils";

const ManualInput: React.FC<{
  handleSaveResults: () => void;
}> = ({ handleSaveResults }) => {
  const [forceRerender, setForceRerender] = useState(false);
  const { competitionState, updateSolve, currentResultsRef } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const solveProp: keyof ResultEntry = `solve${
    competitionState.currentSolveIdx + 1
  }` as keyof ResultEntry;
  const formattedTime = currentResultsRef.current[solveProp].toString();

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
      handleSaveResults();
    } else if (e.key === "ArrowLeft" || e.key === "ArrowRight") {
      e.preventDefault();
    }
  };

  return (
    <div>
      <Input
        size="lg"
        placeholder="Enter your time or solution..."
        sx={{
          marginBottom: 2,
          marginTop: 2,
          input: { caretColor: "transparent" },
        }}
        value={formattedTime}
        onChange={handleTimeInputChange}
        onKeyDown={handleKeyDown}
        onClick={(e) => {
          const target = e.target as HTMLInputElement;
          if (target.tagName === "INPUT")
            target.setSelectionRange(
              target?.value?.length,
              target?.value?.length
            );
        }}
        autoFocus
      />
    </div>
  );
};

export default ManualInput;
