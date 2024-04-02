import { Alert, Button, Card, CircularProgress, Typography } from "@mui/joy";
import { AuthContextType, CompetitionData, CompetitionEvent, CompetitionState } from "../../Types";
import { ChangeEvent, useContext, useEffect, useState } from "react";
import { getAvailableEvents, initialCompetitionState, updateCompetition } from "../../utils";
import { useNavigate, useParams } from "react-router-dom";

import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { AuthContext } from "../../context/AuthContext";
import { CompetitionEditProps } from "../../Types";
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { getCompetitionById } from "../../utils";

const CompetitionEdit: React.FC<CompetitionEditProps> = ({ edit }) => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const { authState } = useContext(AuthContext) as AuthContextType
    const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>([]);
    const [competitionState, setCompetitionState] = useState<CompetitionState>(initialCompetitionState);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>('');

    useEffect(() => {
        if (!authState.authenticated || !authState.admin) {
            navigate("/");
        }

        setIsLoading(true);
        setError('');

        if (edit) {
            getCompetitionById(id)
                .then((info: CompetitionData | undefined) => {
                    if (info === undefined) {
                        navigate('/not-found')
                    } else {
                        setCompetitionState({...competitionState, ...info});
                        setIsLoading(false);
                    }
                })
                .catch(err => {
                    setIsLoading(false);
                    setError(err.message);
                });
        }

        getAvailableEvents()
            .then((res: CompetitionEvent[]) => {
                setAvailableEvents(res)
                setIsLoading(false);
            })
            .catch(err => {
                setIsLoading(false);
                setError(err.message);
            })
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
            {error ? <Alert color="danger">{error}</Alert> : isLoading ? <CircularProgress /> : 
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
                                value={competitionState.startdate}
                                onChange={(e) => setCompetitionState({...competitionState, startdate: !isNaN(Date.parse(e.target.value)) ? e.target.value : competitionState.startdate})}
                            />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Ending date:</span>&nbsp;&nbsp;
                            <input
                                type="datetime-local"
                                value={competitionState.enddate}
                                onChange={(e) => setCompetitionState({...competitionState, enddate: !isNaN(Date.parse(e.target.value)) ? e.target.value : competitionState.enddate})}
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
            </LocalizationProvider>}
        </div>
    )
};

export default CompetitionEdit;