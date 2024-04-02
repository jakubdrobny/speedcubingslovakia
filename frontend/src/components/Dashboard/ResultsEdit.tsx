import { Box, Button, Card, Chip, Divider, FormControl, FormHelperText, FormLabel, Grid, Input, Option, Select, Stack, Typography } from "@mui/joy";
import { CompetitionEvent, ResultEntry } from "../../Types";
import { getAvailableEvents, getResults, reformatTime, sendResults } from "../../utils";
import { useCallback, useEffect, useState } from "react";

const ResultsEdit = () => {
    const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>([]);
    const [competitorName, setCompetitorName] = useState<string>('');
    const [competitionName, setCompetitionName] = useState<string>('');
    const [competitionEvent, setCompetitionEvent] = useState<string>();
    const [results, setResults] = useState<ResultEntry[]>([]);
    const [selectError, setSelectError] = useState<boolean>(false);

    useEffect(() => {
        getAvailableEvents()
            .then(res => {
                setAvailableEvents(res)
                if (res.length > 0)
                    setCompetitionEvent(res[0].displayname);
            })
            .catch(console.error);
    }, []);

    const handleQuery = () => {
        if (competitionEvent === undefined) {
            setSelectError(true);
            return;
        }
        
        getResults(competitorName, competitionName, availableEvents.find(e => e.displayname === competitionEvent))
            .then(res => setResults(res))
            .catch(console.error);
    }

    const updateSolve = (newTime: string, resultsIdx: number, solveProp: string) => {
        const newResults = results.map((val, idx) => idx === resultsIdx ? {...val, [solveProp]: newTime} : {...val});
        setResults(newResults);
    };

    const handleTimeInputChange = (newValue: string, oldValue: string, resultsIdx: number, solveProp: string) => {
        if (results[resultsIdx].eventname === "FMC") {
            updateSolve(newValue, resultsIdx, solveProp);
            return;
        }
        
        // character deleted
        if (newValue.length + 1 === oldValue.length) {
            if (newValue.endsWith('N')) {
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
        const match = results[resultsIdx].format.match(/\d+$/)?.[0]
        const noOfSolves = match ? parseInt(match) : 1
        return ['solve1', 'solve2', 'solve3', 'solve4', 'solve5'].slice(0, noOfSolves);
    }

    const saveResult = (resultsIdx: number) => {
        sendResults(results[resultsIdx]);
    }

    return (
        <Stack spacing={4} sx={{marginTop: "2em", marginBottom: "2em"}}>
            <Typography level="h2">Edit results</Typography>
            <Card>
                <Stack spacing={2}>
                    <Typography level="h3" sx={{borderBottom: '1px solid #636D7433'}}>Query builder</Typography>
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
                    {competitionEvent && <FormControl>
                        <FormLabel>Event</FormLabel>
                        <Select
                            value={competitionEvent}
                            onChange={(e, val) => {setCompetitionEvent(val || ''); setSelectError(false);}}
                            required
                            renderValue={(event) => (
                                <Box sx={{ display: 'flex', gap: '0.25rem' }}>
                                    <Chip variant="soft" color="primary">
                                        <span className={`cubing-icon event-${event?.label}`} />&nbsp;
                                        {event?.value}
                                    </Chip>
                                </Box>
                            )}
                            color={selectError ? 'danger' : 'neutral'}
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
                        {selectError && <FormHelperText sx={{color: 'red'}}>This field is required. Please choose an event.</FormHelperText>}
                    </FormControl>}
                    <Button type="submit" onClick={handleQuery}>Query</Button>
                </Stack>
            </Card>
            <Card>
                <Stack spacing={2}>
                    <Typography level="h3" sx={{borderBottom: '1px solid #636D7433'}}>Results</Typography>
                    {results.map((result: ResultEntry, resultIdx: number) => {
                        console.log(result);
                        return (
                            <Card key={result.id}>
                                <Grid container>
                                    <Grid xs={6}>
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
                                                        <span className={`cubing-icon event-${result.iconcode}`}>
                                                            &nbsp;{result.eventname}
                                                        </span>
                                                    </Chip>
                                                </Typography>
                                            </div>
                                        </Stack>
                                    </Grid>
                                    <Grid xs={6}>
                                        <Stack spacing={1}>
                                            {getSolveProps(resultIdx).map((solveProp, solveIdx) => {
                                                return (
                                                    <FormControl>
                                                        <FormLabel>Solve {solveIdx+1}</FormLabel>
                                                        <Input
                                                            size="sm"
                                                            placeholder="Enter your time or solution..."
                                                            value={results[resultIdx][solveProp as keyof ResultEntry].toString()}
                                                            onChange={(e) => handleTimeInputChange(e.target.value, result[solveProp as keyof ResultEntry].toString(), resultIdx, solveProp)}
                                                        />
                                                    </FormControl>
                                                );
                                            })}
                                        </Stack>
                                    </Grid>
                                </Grid>
                                <Button type="submit" onClick={() => saveResult(resultIdx)}>Save</Button>
                            </Card>
                        );
                    })}
                </Stack>
            </Card>
        </Stack>
    )
}

export default ResultsEdit;
