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
    let newValue = e.currentTarget.value;
    console.log("original newValue: ", newValue);
    const target = e.currentTarget;

    if (
      competitionState.events[competitionState.currentEventIdx].displayname ===
      "FMC"
    ) {
      updateSolve(newValue);
      target.setSelectionRange(newValue.length, newValue.length);
      return;
    }

    if (newValue.length > MAX_MANUAL_INPUT_LENGTH) return;

    // character deleted
    if (newValue.length + 1 === formattedTime.length) {
      if (newValue.endsWith("N")) {
        newValue = "0.00";
      } else {
        newValue = reformatTime(newValue);
      }
    } else {
      if (newValue.endsWith("d")) {
        newValue = "DNF";
      } else if (newValue.endsWith("s")) {
        newValue = "DNS";
      } else if (/\d$/.test(newValue.slice(-1))) {
        newValue = reformatTime(newValue, true);
      } else {
        newValue = "DNF";
      }
    }

    console.log("new newValue: ", newValue);
    updateSolve(newValue);
    target.setSelectionRange(newValue.value.length, newValue.value.length);
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
          //input: { caretColor: "transparent" },
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
