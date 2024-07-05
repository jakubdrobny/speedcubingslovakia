import { Button, Typography } from "@mui/joy";
import {
  CompetitionContextType,
  ScrambleSet,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import { East, West } from "@mui/icons-material";
import React, { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import DefaultScramble from "../../images/DefaultScramble";
import { Stack } from "@mui/system";
import { TimerInputContext } from "../../context/TimerInputContext";

const Scramble: React.FC<{ ismbld: boolean }> = ({ ismbld }) => {
  const [scrambleImg, setScrambleImg] = useState<string>();
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
      ) !== undefined
    ) {
      const scrambleSet = competitionState.scrambles.find(
        (s: ScrambleSet) =>
          s.event.displayname ===
          competitionState.events[competitionState.currentEventIdx].displayname
      ) as ScrambleSet;
      setScrambleImg(
        scrambleSet.scrambles[competitionState.currentSolveIdx].img
      );
    }
  }, [competitionState.currentSolveIdx, competitionState.scrambles]);

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
      {scrambleImg === undefined ? (
        <DefaultScramble />
      ) : (
        <img
          src={`${process.env.REACT_APP_SCRAMBLE_IMAGES_PATH}/${scrambleImg}`}
          alt={`${competitionState?.id}/${
            competitionState?.events[competitionState?.currentEventIdx]
              ?.displayname
          }/scramble${competitionState?.currentSolveIdx + 1}`}
        />
      )}
    </div>
  );
};

export default Scramble;
