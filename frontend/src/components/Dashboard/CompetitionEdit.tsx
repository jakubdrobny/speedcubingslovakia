import { AuthContextType, CompetitionData, CompetitionEvent, CompetitionState, InputMethod } from "../../Types";
import { Box, Button, Card, Chip, Input, Option, Select, Typography } from "@mui/joy";
import { ChangeEvent, ReactNode, useContext, useEffect, useState } from "react";
import dayjs, { Dayjs } from 'dayjs';
import { getAvailableEvents, updateCompetition } from "../../utils";
import { useNavigate, useParams } from "react-router-dom";

import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { AuthContext } from "../../context/AuthContext";
import { CompetitionEditProps } from "../../Types";
import { DateTimeField } from '@mui/x-date-pickers/DateTimeField';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { SelectChangeEvent } from "@mui/material/Select";
import { getCompetitionById } from "../../utils";

const initialState: CompetitionState = {
    id: "",
    name: "",
    startdate: new Date(),
    enddate: new Date(),
    events: [],
    currentEventIdx: 0,
    noOfSolves: 1,
    currentSolveIdx: 0,
    scrambles: [],
    inputMethod: InputMethod.Manual,
    results: {
        id: 0,
        userid: 0,
        solve1: '',
        solve2: '',
        solve3: '',
        solve4: '',
        solve5: '',
        comment: '',
        statusid: 0,
    },
    penalties: Array(5).fill('0')
};

const CompetitionEdit: React.FC<CompetitionEditProps> = ({ edit }) => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const { authState } = useContext(AuthContext) as AuthContextType
    const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>([]);
    const [competitionState, setCompetitionState] = useState<CompetitionState>(initialState);

    useEffect(() => {
        if (!authState.authenticated || !authState.admin) {
            navigate("/");
        }

        if (edit) {
            getCompetitionById(id)
                .then((info: CompetitionData | undefined) => {
                    if (info === undefined) {
                        navigate('/not-found')
                    } else {
                        setCompetitionState({...competitionState, ...info});
                    }
                })
                .catch(console.error);
        }

        getAvailableEvents()
            .then((res: CompetitionEvent[]) => setAvailableEvents(res))
            .catch(console.error)
    }, []);

    const handleSelectedEventsChange = (event: ChangeEvent<HTMLSelectElement>) => {
        const selectedOptions = Array.from(event.target.selectedOptions);
        const selectedEvents = selectedOptions.map(o => availableEvents.find(e => e.displayname === o.value) as CompetitionEvent);
        console.log(selectedEvents)
        setCompetitionState({...competitionState, events: selectedEvents});
    };

    const handleSubmit = () => {
        updateCompetition(competitionState, edit);
    }

    return (
        <div style={{marginTop: "2em"}}>
            <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale="sk">
                <Card>
                    <div style={{borderBottom: "1px solid #CDD7E1"}}>
                        <Typography sx={{fontWeight: 'bold'}}>{edit ? `Edit ${competitionState.name}` : 'Create'} competition</Typography>
                    </div>
                    <div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Name:</span>&nbsp;&nbsp;
                            <input
                                placeholder="Enter competition name..."
                                value={competitionState.name}
                                onChange={(e) => setCompetitionState({...competitionState, name: e.target.value})}
                            />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Starting date:</span>&nbsp;&nbsp;
                            <input
                                type="datetime-local"
                                value={competitionState.startdate.toISOString().slice(0, 16)}
                                onChange={(e) => setCompetitionState({...competitionState, startdate: !isNaN(Date.parse(e.target.value)) ? new Date(e.target.value) : competitionState.startdate})}
                            />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Ending date:</span>&nbsp;&nbsp;
                            <input
                                type="datetime-local"
                                value={competitionState.enddate.toISOString().slice(0, 16)}
                                onChange={(e) => setCompetitionState({...competitionState, enddate: !isNaN(Date.parse(e.target.value)) ? new Date(e.target.value) : competitionState.enddate})}
                            />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Events:</span>&nbsp;&nbsp;
                            <select
                                multiple
                                value={competitionState.events.map(e => e.displayname)}
                                onChange={handleSelectedEventsChange}
                            >
                                {availableEvents.map((ev: CompetitionEvent) => (
                                    <option key={ev.id} value={ev.displayname} className={`cubing-icon event-${ev.iconcode}`}>
                                        {ev.displayname}
                                    </option>
                                ))}
                            </select>
                            <span>&nbsp;(hold shift or ctrl to select multiple events)</span>
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <Button type="submit" onClick={handleSubmit}>{edit ? "Edit" : "Create"} competition</Button>
                        </div>
                    </div>
                </Card>
            </LocalizationProvider>
        </div>
    )
};

export default CompetitionEdit;