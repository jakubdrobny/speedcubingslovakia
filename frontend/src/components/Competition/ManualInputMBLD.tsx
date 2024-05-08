import { CompetitionContextType, ResultEntry } from "../../Types";
import { Input, Stack, Typography } from "@mui/joy";
import { useContext, useEffect } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { competitionOnGoing } from "../../utils";
import useState from "react-usestateref";

const ManualInputMBLD = () => {
  const [forceRerender, setForceRerender] = useState(false);
  const { competitionStateRef, updateSolve, currentResultsRef } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const solveProp: keyof ResultEntry = `solve${
    competitionStateRef.current.currentSolveIdx + 1
  }` as keyof ResultEntry;
  const formattedTime = currentResultsRef.current[solveProp].toString();

  const [solvedCubes, setSolvedCubes, solvedCubesRef] = useState(
    formattedTime === "DNS" ? "0" : formattedTime?.split("/")[0]
  );

  const [attemptedCubes, setAttemptedCubes, attemptedCubesRef] =
    useState<string>(
      formattedTime === "DNS" ? "0" : formattedTime?.split(" ")[0].split("/")[1]
    );
  const [hours, setHours, hoursRef] = useState(
    formattedTime === "DNS"
      ? "00"
      : formattedTime?.split(" ")[1]?.split(":")[0].padStart(2, "0")
  );
  const [minutes, setMinutes, minutesRef] = useState<string>(
    formattedTime === "DNS"
      ? "00"
      : formattedTime?.split(" ")[1]?.split(":")[1]?.padStart(2, "0")
  );
  const [seconds, setSeconds, secondsRef] = useState(
    formattedTime === "DNS"
      ? "00"
      : formattedTime?.split(" ")[1]?.split(":")[2].padStart(2, "0")
  );

  const handleSomethingChanged = () => {
    updateSolve(
      `${solvedCubesRef.current}/${attemptedCubesRef.current} ${hoursRef.current}:${minutesRef.current}:${secondsRef.current}`
    );
  };

  useEffect(() => {
    setForceRerender((p) => !p);
  }, [
    competitionStateRef.current.currentEventIdx,
    competitionStateRef.current.currentSolveIdx,
  ]);

  useEffect(() => {
    setSolvedCubes(
      formattedTime === "DNS" ? "0" : formattedTime?.split("/")[0]
    );
    setAttemptedCubes(
      formattedTime === "DNS" ? "0" : formattedTime?.split(" ")[0].split("/")[1]
    );
    setHours(
      formattedTime === "DNS"
        ? "00"
        : formattedTime?.split(" ")[1]?.split(":")[0].padStart(2, "0")
    );
    setMinutes(
      formattedTime === "DNS"
        ? "00"
        : formattedTime?.split(" ")[1]?.split(":")[1].padStart(2, "0")
    );
    setSeconds(
      formattedTime === "DNS"
        ? "00"
        : formattedTime?.split(" ")[1]?.split(":")[2].padStart(2, "0")
    );
  }, [formattedTime]);

  return (
    <div>
      <Stack
        spacing={3}
        sx={{
          marginBottom: 2,
          marginTop: 2,
          display: "flex",
          justifyContent: "center",
        }}
        direction="row"
      >
        <Stack
          direction="row"
          sx={{
            alignItems: "center",
          }}
          spacing={1}
        >
          {[0, 1, 2].map((_, idx) =>
            idx & 1 ? (
              <Typography key={idx}>
                <b>/</b>
              </Typography>
            ) : (
              <Input
                key={idx}
                size="lg"
                type="number"
                slotProps={{ input: { min: 0, max: 69, maxLength: 2 } }}
                value={
                  idx === 0 ? solvedCubesRef.current : attemptedCubesRef.current
                }
                onChange={(e) => {
                  let nval = e.target.value;
                  if (nval === "") nval = "0";
                  else nval = nval.slice(0, 2);
                  if (nval.length > 1 && nval[0] === "0") nval = nval[1];
                  if (idx === 0) setSolvedCubes(nval);
                  else setAttemptedCubes(nval);
                  handleSomethingChanged();
                }}
                disabled={!competitionOnGoing(competitionStateRef.current)}
                sx={{ width: "3em", padding: 0 }}
              />
            )
          )}
        </Stack>
        <Stack
          direction="row"
          sx={{
            alignItems: "center",
          }}
          spacing={1}
        >
          {[0, 1, 2, 3, 4].map((_, idx) =>
            idx & 1 ? (
              <Typography key={idx}>
                <b>:</b>
              </Typography>
            ) : (
              <Input
                key={idx}
                size="lg"
                type="number"
                slotProps={{
                  input: { min: 0, max: idx === 0 ? 1 : 59, maxLength: 2 },
                }}
                value={idx === 0 ? hours : idx === 2 ? minutes : seconds}
                onChange={(e) => {
                  let nval = e.target.value;
                  if (parseInt(nval) > (idx === 0 ? 1 : 59)) return;
                  while (nval.length > 2) nval = nval.slice(1);
                  nval = nval.slice(0, 2);
                  nval = nval.padStart(2, "0");
                  if (idx === 0) setHours(nval);
                  else if (idx === 2) setMinutes(nval);
                  else setSeconds(nval);
                  handleSomethingChanged();
                }}
                disabled={!competitionOnGoing(competitionStateRef.current)}
                sx={{ width: "3em", padding: 0 }}
              />
            )
          )}
        </Stack>
      </Stack>
    </div>
  );
};

export default ManualInputMBLD;
