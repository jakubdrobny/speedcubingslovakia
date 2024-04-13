import { Alert, Button, Card, CircularProgress, Grid } from "@mui/joy";
import { CompetitionContextType, InputMethod, ResultEntry } from "../../Types";
import { East, Keyboard, Timer, West } from "@mui/icons-material";
import {
  competitionOnGoing,
  getResultsFromCompetitionAndEvent,
} from "../../utils";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import ManualInput from "./ManualInput";
import Penalties from "./Penalties";
import Scramble from "./Scramble";
import TimerInput from "./TimerInput";

const Compete = () => {
  const {
    competitionState,
    updateCurrentSolve,
    toggleInputMethod,
    saveResults,
    setCurrentResults,
    setSuspicousModalOpen,
    loadingState,
  } = useContext(CompetitionContext) as CompetitionContextType;
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    if (
      competitionState.id === undefined ||
      competitionState.events.length === 0
    )
      return;

    setIsLoading(true);

    getResultsFromCompetitionAndEvent(
      competitionState.id,
      competitionState.events[competitionState.currentEventIdx]
    )
      .then((resultEntry: ResultEntry) => {
        setIsLoading(false);
        setCurrentResults(resultEntry);
        if (!resultEntry.status.approvalFinished) setSuspicousModalOpen(true);
      })
      .catch((err) => {
        setIsLoading(false);
        setError(err.message);
      });
  }, []);

  const handleSaveResults = () => {
    setIsLoading(true);
    saveResults()
      .then(() => setIsLoading(false))
      .catch((err) => {
        setIsLoading(false);
        setError(err.message);
      });
  };

  return (
    <>
      {error ? (
        <Alert color="danger">{error}</Alert>
      ) : loadingState.results ? (
        <div style={{ display: "flex", justifyContent: "center" }}>
          <CircularProgress />
        </div>
      ) : (
        <Card>
          <Grid container>
            <Grid xs={4}>
              <Button
                variant="outlined"
                onClick={() =>
                  updateCurrentSolve(
                    (competitionState.currentSolveIdx -
                      1 +
                      competitionState.noOfSolves) %
                      competitionState.noOfSolves
                  )
                }
              >
                <West />
                &nbsp; Previous
              </Button>
            </Grid>
            <Grid
              xs={4}
              sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
              }}
            >
              Solve {competitionState.currentSolveIdx + 1}
            </Grid>
            <Grid xs={4} sx={{ display: "flex", justifyContent: "flex-end" }}>
              <Button
                variant="outlined"
                onClick={() =>
                  updateCurrentSolve(
                    (competitionState.currentSolveIdx + 1) %
                      competitionState.noOfSolves
                  )
                }
              >
                Next&nbsp;
                <East />
              </Button>
            </Grid>
          </Grid>
          <Scramble />
          <Grid container>
            <Grid
              xs={12}
              sx={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <h3
                onClick={toggleInputMethod}
                style={{ display: "flex", alignItems: "center" }}
              >
                {competitionState.inputMethod === InputMethod.Manual ? (
                  <>
                    Manual&nbsp;
                    <Keyboard />
                  </>
                ) : (
                  <>
                    Timer&nbsp;
                    <Timer />
                  </>
                )}
              </h3>
            </Grid>
            {isLoading ? (
              <Grid
                xs={12}
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center",
                }}
              >
                <CircularProgress />
              </Grid>
            ) : (
              <>
                <Grid xs={12}>
                  {competitionState.inputMethod === InputMethod.Manual ||
                  !competitionOnGoing(competitionState) ? (
                    <ManualInput />
                  ) : (
                    <TimerInput />
                  )}
                </Grid>
                <Penalties />
                <Grid xs={12} sx={{ marginTop: 2 }}>
                  <Button
                    color="primary"
                    variant="solid"
                    onClick={handleSaveResults}
                    sx={{ width: "100%" }}
                    disabled={!competitionOnGoing(competitionState)}
                    loading={isLoading}
                  >
                    Save
                  </Button>
                </Grid>
              </>
            )}
          </Grid>
        </Card>
      )}{" "}
    </>
  );
};

export default Compete;
