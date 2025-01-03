import { Box, Button, Card, Grid } from "@mui/joy";
import {
  CompetitionContextType,
  InputMethod,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";
import { East, Keyboard, Timer, West } from "@mui/icons-material";
import {
  competitionOnGoing,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import { useContext, useEffect, useRef, useState } from "react";

import AveragePreview from "../AveragePreview/AveragePreview";
import { CompetitionContext } from "../../context/CompetitionContext";
import LoadingComponent from "../Loading/LoadingComponent";
import ManualInput from "./ManualInput";
import ManualInputMBLD from "./ManualInputMBLD";
import Penalties from "./Penalties";
import Scramble from "./Scramble";
import TimerInput from "./TimerInput";
import { TimerInputContext } from "../../context/TimerInputContext";

const Compete = () => {
  const {
    competitionState,
    updateCurrentSolve,
    toggleInputMethod,
    saveResults,
    loadingState,
    setLoadingState,
    fetchCompeteResultEntry,
  } = useContext(CompetitionContext) as CompetitionContextType;
  const { timerInputState } = useContext(
    TimerInputContext,
  ) as TimerInputContextType;
  const ismbld =
    competitionState?.events[competitionState?.currentEventIdx]?.iconcode ===
    "333mbf";
  const [showResultsModal, setShowResultsModal] = useState<boolean>(false);
  const competeRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (
      competitionState.id === undefined ||
      competitionState.events.length === 0
    )
      return;

    fetchCompeteResultEntry();
  }, []);

  const handleSaveResults = (moveIndex: boolean) => {
    setLoadingState({ ...loadingState, results: true, compinfo: false });
    saveResults()
      .then(() => {
        setLoadingState({ ...loadingState, results: false, error: {} });
        if (moveIndex)
          updateCurrentSolve(
            (competitionState.currentSolveIdx + 1) %
              competitionState.noOfSolves,
          );
        setShowResultsModal(true);
        if (competeRef && competeRef.current) {
          competeRef.current.scrollIntoView();
        }
      })
      .catch((err) =>
        setLoadingState({
          ...loadingState,
          results: false,
          error: getError(err),
        }),
      );
  };

  return (
    <div ref={competeRef}>
      {!isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <Card>
          <Grid container>
            <Grid xs={4} sx={{ zIndex: 10, backgroundColor: "#FBFCFE" }}>
              <Button
                variant="outlined"
                onClick={() =>
                  updateCurrentSolve(
                    (competitionState.currentSolveIdx -
                      1 +
                      competitionState.noOfSolves) %
                      competitionState.noOfSolves,
                  )
                }
                disabled={
                  timerInputState.currentState !==
                    TimerInputCurrentState.NotSolving || loadingState.results
                }
                sx={{ backgroundColor: "#FBFCFE" }}
              >
                <West />
                &nbsp;
                <div> Previous</div>
              </Button>
            </Grid>
            <Grid
              xs={4}
              sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                zIndex: 0,
              }}
            >
              {ismbld ? "Attempt" : "Solve"}{" "}
              {competitionState.currentSolveIdx + 1}
            </Grid>
            <Grid xs={4} sx={{ display: "flex", justifyContent: "flex-end" }}>
              <Button
                variant="outlined"
                onClick={() =>
                  updateCurrentSolve(
                    (competitionState.currentSolveIdx + 1) %
                      competitionState.noOfSolves,
                  )
                }
                disabled={
                  timerInputState.currentState !==
                    TimerInputCurrentState.NotSolving || loadingState.results
                }
              >
                Next&nbsp;
                <East />
              </Button>
            </Grid>
          </Grid>
          <AveragePreview
            showResultsModal={showResultsModal}
            loadingResults={loadingState.results}
          />
          <Scramble ismbld={ismbld} />
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
                onClick={() => {
                  if (
                    timerInputState.currentState ===
                    TimerInputCurrentState.NotSolving
                  )
                    toggleInputMethod();
                }}
                style={{
                  display: "flex",
                  alignItems: "center",
                  marginTop: "0.5em",
                  marginBottom: "0.25em",
                }}
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
            <Grid xs={12}>
              {loadingState.results ? (
                <Box sx={{ pt: 2.25, pb: 2.25 }}>
                  <LoadingComponent title="" />
                </Box>
              ) : competitionState.inputMethod === InputMethod.Manual ||
                !competitionOnGoing(competitionState) ? (
                ismbld ? (
                  <ManualInputMBLD />
                ) : (
                  <ManualInput handleSaveResults={handleSaveResults} />
                )
              ) : (
                <TimerInput handleSaveResults={handleSaveResults} />
              )}
            </Grid>
            <Penalties />
            <Grid xs={12} sx={{ marginTop: 2 }}>
              <Button
                color="primary"
                variant="solid"
                onClick={() => handleSaveResults(false)}
                sx={{ width: "100%" }}
                disabled={
                  !competitionOnGoing(competitionState) ||
                  timerInputState.currentState !==
                    TimerInputCurrentState.NotSolving ||
                  loadingState.results
                }
              >
                Save
              </Button>
            </Grid>
          </Grid>
        </Card>
      )}
    </div>
  );
};

export default Compete;
