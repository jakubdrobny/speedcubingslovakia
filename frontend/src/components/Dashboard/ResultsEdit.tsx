import { Box, Button, Card, Chip, FormControl, FormHelperText, FormLabel, IconButton, Input, Option, Select, Stack } from "@mui/joy";
import { ChangeEvent, SyntheticEvent, useEffect, useState } from "react";

import { Close } from "@mui/icons-material";
import { CompetitionEvent } from "../../Types";
import { Form } from "react-router-dom";
import { SelectChangeEvent } from "@mui/material";
import { getAvailableEvents } from "../../utils";

const ResultsEdit = () => {
    const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>([]);
    const [competitorName, setCompetitorName] = useState<string>('');
    const [competitionName, setCompetitionName] = useState<string>('');
    const [competitionEvent, setCompetitionEvent] = useState<string>();

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
        
    }

    return (
        <div>
            <Card sx={{marginTop: "2em"}}>
                <form noValidate>
                    <Stack spacing={2}>
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
                                onChange={(e, val) => setCompetitionEvent(val || '')}
                                required
                                renderValue={(event) => (
                                    <Box sx={{ display: 'flex', gap: '0.25rem' }}>
                                        <Chip variant="soft" color="primary">
                                            <span className={`cubing-icon event-${event?.label}`} />&nbsp;
                                            {event?.value}
                                        </Chip>
                                    </Box>
                                  )}
                                  
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
                        </FormControl>}
                        <Button type="submit">Query</Button>
                    </Stack>
                </form>
            </Card>
        </div>
    )
}

export default ResultsEdit;