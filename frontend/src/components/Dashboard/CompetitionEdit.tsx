import { AuthContextType, CompetitionContextType, CompetitionData, CompetitionEvent } from "../../Types";
import { Box, Card, Chip, Input, Option, Select, Typography } from "@mui/joy";
import { useContext, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { AuthContext } from "../../context/AuthContext";
import { CompetitionContext } from "../Competition/CompetitionContext";
import { CompetitionEditProps } from "../../Types";
import { DateTimeField } from '@mui/x-date-pickers/DateTimeField';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import dayjs from 'dayjs';
import { getAvailableEvents } from "../../utils";
import { getCompetitionById } from "../../utils";

const CompetitionEdit: React.FC<CompetitionEditProps> = ({ edit }) => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const { competitionState, updateBasicInfo } = useContext(CompetitionContext) as CompetitionContextType
    const { authState } = useContext(AuthContext) as AuthContextType
    const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>([]);

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
                        updateBasicInfo(info);
                    }
                })
                .catch(console.error);
        }

        getAvailableEvents()
            .then((res: CompetitionEvent[]) => setAvailableEvents(res))
            .catch(console.error)
    }, []);
    
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
                            <Input
                                placeholder="Enter competition name..."
                                size="lg"
                            />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Starting date:</span>&nbsp;&nbsp;
                                <DateTimeField
                                    defaultValue={dayjs(Date.now())}
                                    format="LLL"
                                />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Ending date:</span>&nbsp;&nbsp;
                            <DateTimeField
                                defaultValue={dayjs(Date.now())}
                                format="LLL"
                            />
                        </div>
                        <div style={{display: 'flex', alignItems: 'center', marginTop: "1em", marginBottom: "1em"}}>
                            <span style={{fontWeight: 'bold'}}>Events:</span>&nbsp;&nbsp;
                            <Select
                                multiple
                                renderValue={(selected) => (
                                    <Box sx={{ display: 'flex', gap: '0.25rem' }}>
                                    {selected.map((selectedOption) => (
                                        <Chip variant="soft" color="primary">
                                            <span className={`cubing-icon event-${selectedOption.label}`} />&nbsp;
                                            {selectedOption.value.toString()}
                                        </Chip>
                                    ))}
                                    </Box>
                                )}
                                sx={{
                                    minWidth: '15rem',
                                }}
                                slotProps={{
                                    listbox: {
                                    sx: {
                                        width: '100%',
                                    },
                                    },
                                }}
                            >
                                {availableEvents.map((ev: CompetitionEvent) => {
                                    return (
                                        <Option key={ev.id} label={ev.iconcode} value={ev.displayname}>
                                            <span className={`cubing-icon event-${ev.iconcode}`} />
                                            {ev.displayname}
                                        </Option>
                                    )
                                })}
                            </Select>
                        </div>
                    </div>
                </Card>
            </LocalizationProvider>
        </div>
    )
};

export default CompetitionEdit;