import { Button, Typography } from "@mui/joy";
import {
  CompetitionContextType,
  ScrambleSet,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import { East, West } from "@mui/icons-material";
import React, { useContext, useEffect, useMemo, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { Stack } from "@mui/system";
import { TimerInputContext } from "../../context/TimerInputContext";
import { useRef } from "react";

const Scramble: React.FC<{ ismbld: boolean }> = ({ ismbld }) => {
  const scrambleImgRef = useRef<HTMLDivElement>(null);
  const { competitionState } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const { timerInputState } = useContext(
    TimerInputContext
  ) as TimerInputContextType;
  const scramble =
    competitionState &&
    competitionState.scrambles &&
    competitionState.events &&
    competitionState.currentEventIdx < competitionState.events.length &&
    competitionState.scrambles.find(
      (s: ScrambleSet) =>
        s.event.displayname ===
        competitionState.events[competitionState.currentEventIdx].displayname
    ) !== undefined
      ? (
          competitionState.scrambles.find(
            (s: ScrambleSet) =>
              s.event.displayname ===
              competitionState.events[competitionState.currentEventIdx]
                .displayname
          ) as ScrambleSet
        ).scrambles[competitionState.currentSolveIdx].scramble
      : "";

  useEffect(() => {
    if (
      competitionState &&
      competitionState.scrambles &&
      competitionState.events &&
      competitionState.currentEventIdx < competitionState.events.length &&
      competitionState.scrambles.find(
        (s: ScrambleSet) =>
          s.event.displayname ===
          competitionState.events[competitionState.currentEventIdx].displayname
      ) !== undefined &&
      scrambleImgRef !== null &&
      scrambleImgRef.current !== null
    ) {
      const scrambleSet = competitionState.scrambles.find(
        (s: ScrambleSet) =>
          s.event.displayname ===
          competitionState.events[competitionState.currentEventIdx].displayname
      ) as ScrambleSet;
      scrambleImgRef.current.innerHTML =
        scrambleSet.scrambles[competitionState.currentSolveIdx].svgimg;
    }
  }, [competitionState.currentSolveIdx]);

  const [scramblePage, setScramblePage] = useState(0);
  const scrambles = scramble.split("\n").length;
  const pages = Math.ceil(scrambles / 10);

  const formatScramble = () => {
    if (!ismbld)
      return (
        <Typography style={{ whiteSpace: "pre-line" }}>{scramble}</Typography>
      );

    let newScramble = scramble
      .split("\n")
      .slice(scramblePage * 10, Math.min((scramblePage + 1) * 10, scrambles))
      .map((scr, idx) => (
        <div key={idx}>
          <b>
            {scramblePage * 10 + 1 + idx}
            {"."}
          </b>
          &nbsp;
          {scr}
        </div>
      ));

    return (
      <Stack spacing={3} sx={{ mb: "1em" }}>
        <Stack spacing={0.25}>{newScramble}</Stack>
        <Stack
          direction="row"
          spacing={2}
          style={{ display: "flex", justifyContent: "center" }}
        >
          <Button
            variant="outlined"
            onClick={() => setScramblePage((scramblePage - 1 + pages) % pages)}
          >
            <West />
            &nbsp;
            <div>Previous</div>
          </Button>
          <Button
            variant="outlined"
            onClick={() => setScramblePage((scramblePage + 1) % pages)}
          >
            <div>Next</div>
            &nbsp;
            <East />
          </Button>
        </Stack>
      </Stack>
    );
  };

  return (
    <div
      style={
        [
          TimerInputCurrentState.NotSolving,
          TimerInputCurrentState.GettingReady,
        ].includes(timerInputState.currentState)
          ? { display: "flex", flexDirection: "column", alignItems: "center" }
          : { display: "none" }
      }
    >
      <h3>Scramble{ismbld ? "s" : ""}:</h3>
      {formatScramble()}
      <h3>Preview:</h3>
      <div ref={scrambleImgRef}></div>
    </div>
  );
};

export default Scramble;
