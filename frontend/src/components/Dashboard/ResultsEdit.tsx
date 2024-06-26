import {
  Alert,
  Box,
  Button,
  Card,
  Chip,
  CircularProgress,
  Divider,
  FormControl,
  FormHelperText,
  FormLabel,
  Grid,
  Input,
  Option,
  Select,
  Stack,
  Textarea,
  Typography,
} from "@mui/joy";
import { AuthContextType, CompetitionEvent, ResultEntry } from "../../Types";
import { Check, Close } from "@mui/icons-material";
import {
  authorizeAdmin,
  getAvailableEvents,
  getError,
  getResults,
  logOut,
  reformatTime,
  saveValidation,
  sendResults,
} from "../../utils";
import { useEffect, useState } from "react";

import { MAX_MANUAL_INPUT_LENGTH } from "../../constants";

const ResultsEdit = () => {
  const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>(
    []
  );
  const [competitorName, setCompetitorName] = useState<string>("");
  const [competitionName, setCompetitionName] = useState<string>("");
  const [competitionEvent, setCompetitionEvent] = useState<string>();
  const [results, setResults] = useState<ResultEntry[]>([]);
  const [selectError, setSelectError] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<{
    results: boolean;
    events: boolean;
  }>({ results: false, events: false });
  const [error, setError] = useState<string>("");

  useEffect(() => {
    setIsLoading((ps) => ({ ...ps, events: true }));

    getAvailableEvents()
      .then((res) => {
        setAvailableEvents(res);
        if (res.length > 0) setCompetitionEvent(res[0].displayname);
        setIsLoading((ps) => ({ ...ps, events: false }));
      })
      .catch((err) => {
        setError(getError(err));
        setIsLoading((ps) => ({ ...ps, events: false }));
      });
  }, []);

  const fetchResults = () => {
    setIsLoading((ps) => ({ ...ps, results: true }));

    getResults(
      competitorName,
      competitionName,
      availableEvents.find((e) => e.displayname === competitionEvent)
    )
      .then((res) => {
        setResults(res);
        setIsLoading((ps) => ({ ...ps, results: false }));
        setError("");
      })
      .catch((err) => {
        setIsLoading((ps) => ({ ...ps, results: false }));
        setError(getError(err));
      });
  };

  const handleQuery = () => {
    if (competitionEvent === undefined) {
      setSelectError(true);
      return;
    }

    fetchResults();
  };

  const updateSolve = (
    newTime: string,
    resultsIdx: number,
    solveProp: string
  ) => {
    const newResults = results.map((val, idx) =>
      idx === resultsIdx ? { ...val, [solveProp]: newTime } : { ...val }
    );
    setResults(newResults);
  };

  const handleTimeInputChange = (
    newValue: string,
    oldValue: string,
    resultsIdx: number,
    solveProp: string
  ) => {
    if (results[resultsIdx].eventname === "FMC") {
      updateSolve(newValue, resultsIdx, solveProp);
      return;
    }

    if (newValue.length > MAX_MANUAL_INPUT_LENGTH) return;

    // character deleted
    if (newValue.length + 1 === oldValue.length) {
      if (newValue.endsWith("N")) {
        updateSolve("0.00", resultsIdx, solveProp);
        return;
      } else {
        updateSolve(reformatTime(newValue), resultsIdx, solveProp);
      }
    } else {
      if (newValue.endsWith("d")) {
        updateSolve("DNF", resultsIdx, solveProp);
      } else if (newValue.endsWith("s")) {
        updateSolve("DNS", resultsIdx, solveProp);
      } else if (/\d$/.test(newValue.slice(-1))) {
        updateSolve(reformatTime(newValue, true), resultsIdx, solveProp);
      } else {
        updateSolve("DNF", resultsIdx, solveProp);
      }
    }
  };

  const getSolveProps = (resultsIdx: number): string[] => {
    const match = results[resultsIdx].format.match(/\d+$/)?.[0];
    const noOfSolves = match ? parseInt(match) : 1;
    return ["solve1", "solve2", "solve3", "solve4", "solve5"].slice(
      0,
      noOfSolves
    );
  };

  const saveResult = async (resultsIdx: number) => {
    setIsLoading((ps) => ({ ...ps, results: true }));
    await sendResults(results[resultsIdx]);
    await fetchResults();
  };

  const validateResult = (resultsIdx: number, verdict: boolean) => {
    setIsLoading((ps) => ({ ...ps, results: true }));
    saveValidation(results[resultsIdx], verdict)
      .then((res) => fetchResults())
      .catch((err) => {
        setIsLoading((ps) => ({ ...ps, results: false }));
        setError(getError(err));
      });
  };

  const handleCommentChange = (newComment: string, resultsIdx: number) => {
    const newResults = results.map((val, idx) =>
      idx === resultsIdx ? { ...val, comment: newComment } : { ...val }
    );
    setResults(newResults);
  };

  return (
    <Stack spacing={4} sx={{ marginTop: "2em", marginBottom: "2em" }}>
      <Typography
        level="h2"
        sx={{ pl: 1, borderBottom: "1px solid #636d7433" }}
      >
        Edit results
      </Typography>
      <Card>
        <Stack spacing={2}>
          <Typography level="h3" className="bottom-divider">
            Query builder
          </Typography>
          <FormControl>
            <FormLabel>Competitor name</FormLabel>
            <Input
              placeholder="Enter exact competitor name..."
              value={competitorName}
              onChange={(e) => setCompetitorName(e.target.value)}
            />
            <FormHelperText>Leave empty for all competitors.</FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>Competition name</FormLabel>
            <Input
              placeholder="Enter exact competition name..."
              value={competitionName}
              onChange={(e) => setCompetitionName(e.target.value)}
            />
            <FormHelperText>Leave empty for all competitions.</FormHelperText>
          </FormControl>
          {competitionEvent && (
            <FormControl>
              <FormLabel>Event</FormLabel>
              <Select
                value={competitionEvent}
                onChange={(e, val) => {
                  setCompetitionEvent(val || "");
                  setSelectError(false);
                }}
                required
                renderValue={(event) => (
                  <Box sx={{ display: "flex", gap: "0.25rem" }}>
                    <Chip variant="soft" color="primary">
                      <span className={`cubing-icon event-${event?.label}`} />
                      &nbsp;
                      {event?.value}
                    </Chip>
                  </Box>
                )}
                color={selectError ? "danger" : "neutral"}
                disabled={isLoading.events}
              >
                {availableEvents.map((event: CompetitionEvent) => (
                  <Option
                    key={event.id}
                    value={event.displayname}
                    label={event.iconcode}
                  >
                    <span className={`cubing-icon event-${event.iconcode}`} />
                    {event.displayname}
                  </Option>
                ))}
              </Select>
              {selectError && (
                <FormHelperText sx={{ color: "red" }}>
                  This field is required. Please choose an event.
                </FormHelperText>
              )}
            </FormControl>
          )}
          <Button
            type="submit"
            onClick={handleQuery}
            loading={isLoading.results || isLoading.events}
          >
            Query
          </Button>
        </Stack>
      </Card>
      <Card>
        {error && <Alert color="danger">{error}</Alert>}
        <Stack spacing={2}>
          <Typography level="h3" className="bottom-divider">
            Results
          </Typography>
          {isLoading.results ? (
            <CircularProgress />
          ) : (
            <>
              {results.map((result: ResultEntry, resultIdx: number) => (
                <Card key={result.id}>
                  <Stack spacing={3} sx={{ marginBottom: "0.25em" }}>
                    <Grid container>
                      <Grid
                        xs={6}
                        sx={{
                          display: "flex",
                          justifyContent: "center",
                          alignItems: "center",
                        }}
                      >
                        <Stack spacing={2}>
                          <div>
                            <Typography level="h4">Name:</Typography>
                            <Typography>{result.username}</Typography>
                          </div>
                          <div>
                            <Typography level="h4">Competition:</Typography>
                            <Typography>{result.competitionname}</Typography>
                          </div>
                          <div>
                            <Typography level="h4">Event:</Typography>
                            <Typography component="div">
                              <Chip size="lg" color="primary">
                                <span
                                  className={`cubing-icon event-${result.iconcode}`}
                                >
                                  &nbsp;{result.eventname}
                                </span>
                              </Chip>
                            </Typography>
                          </div>
                          <div
                            style={
                              result.status.approvalFinished
                                ? { display: "none" }
                                : {}
                            }
                          >
                            <Typography level="h4">Resolve status:</Typography>
                            <Stack spacing={2} direction="row">
                              <Button
                                color="danger"
                                variant="soft"
                                onClick={() => validateResult(resultIdx, false)}
                              >
                                <Close />
                                Deny
                              </Button>
                              <Button
                                color="success"
                                variant="soft"
                                onClick={() => validateResult(resultIdx, true)}
                              >
                                <Check />
                                Approve
                              </Button>
                            </Stack>
                          </div>
                          <div
                            style={
                              !result.status.approvalFinished
                                ? { display: "none" }
                                : {}
                            }
                          >
                            <Typography level="h4">Status:</Typography>
                            {result.status.approvalFinished &&
                            result.status.approved === true ? (
                              <div className="mui-joy-btn mui-joy-btn-soft-success">
                                {result.status.displayname}
                              </div>
                            ) : result.status.approvalFinished &&
                              result.status.approved === false ? (
                              <div className="mui-joy-btn mui-joy-btn-soft-danger">
                                {result.status.displayname}
                              </div>
                            ) : (
                              <></>
                            )}
                          </div>
                        </Stack>
                      </Grid>
                      <Grid xs={6}>
                        <Stack spacing={1}>
                          {getSolveProps(resultIdx).map(
                            (solveProp, solveIdx) => {
                              return (
                                <FormControl key={solveProp}>
                                  <FormLabel>Solve {solveIdx + 1}</FormLabel>
                                  <Input
                                    autoFocus
                                    size="sm"
                                    placeholder="Enter your time or solution..."
                                    value={results[resultIdx][
                                      solveProp as keyof ResultEntry
                                    ].toString()}
                                    onChange={(e) =>
                                      handleTimeInputChange(
                                        e.target.value,
                                        result[
                                          solveProp as keyof ResultEntry
                                        ].toString(),
                                        resultIdx,
                                        solveProp
                                      )
                                    }
                                  />
                                </FormControl>
                              );
                            }
                          )}
                          <FormControl>
                            <FormLabel>Comment:</FormLabel>
                            <Textarea
                              value={results[resultIdx].comment}
                              onChange={(e) =>
                                handleCommentChange(e.target.value, resultIdx)
                              }
                              placeholder="Enter a comment to your solutions..."
                              minRows={4}
                              style={{ marginBottom: "1.25em" }}
                            />
                          </FormControl>
                        </Stack>
                      </Grid>
                    </Grid>
                    <Button type="submit" onClick={() => saveResult(resultIdx)}>
                      Save
                    </Button>
                  </Stack>
                </Card>
              ))}
            </>
          )}
        </Stack>
      </Card>
    </Stack>
  );
};

export default ResultsEdit;
