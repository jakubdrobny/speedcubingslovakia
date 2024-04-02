import '../../styles/cubing-icons.css'

import { AuthContextType, CompetitionContextType, CompetitionData, ResultEntry } from "../../Types";
import { Stack, Typography } from '@mui/joy';
import { getCompetitionById, getResultsFromCompetitionAndEvent } from "../../utils";
import { useContext, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { AuthContext } from '../../context/AuthContext';
import { CompetitionContext } from './CompetitionContext';
import CompetitorArea from './CompetitorArea';
import { EventSelector } from './EventSelector';

const Competition = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const { competitionState, updateBasicInfo } = useContext(CompetitionContext) as CompetitionContextType
    const { authState } = useContext(AuthContext) as AuthContextType

    useEffect(() => {
        getCompetitionById(id)
            .then((info: CompetitionData | undefined) => {
                if (info === undefined) navigate('/not-found');
                else {
                    getResultsFromCompetitionAndEvent(authState.token, id, info.events[0])
                        .then((resultEntry: ResultEntry) => updateBasicInfo({...info, results: resultEntry}))
                        .catch(console.error)
                }
            })
            .catch(console.error);
    }, []);

    return (
        <Stack spacing={3} sx={{display: 'flex', alignItems: 'center', margin: "2em 0"}}>
            <Typography level="h1">{competitionState.name}</Typography>
            <Typography>{competitionState.startdate.toLocaleString()} - {competitionState.enddate.toLocaleString()}</Typography>
            <EventSelector />
            <br/>
            <CompetitorArea />
        </Stack>
    );
};

export default Competition;