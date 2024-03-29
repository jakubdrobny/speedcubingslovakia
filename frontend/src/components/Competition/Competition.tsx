import '../../styles/cubing-icons.css'

import { AuthContextType, CompetitionContextType, CompetitionData, ResultEntry } from "../../Types";
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
        <div style={{display: 'flex', flexDirection: 'column', alignItems: 'center'}}>
            <h1>{competitionState.name}</h1>
            <p>{competitionState.startdate.toLocaleString()} - {competitionState.enddate.toLocaleString()}</p>
            <EventSelector />
            <br/>
            <CompetitorArea />
        </div>
    );
};

export default Competition;