import '../../styles/cubing-icons.css'

import { useContext, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { CompetitionContext } from './CompetitionContext';
import { CompetitionContextType } from "../../Types";
import CompetitorArea from './CompetitorArea';
import { EventSelector } from './EventSelector';
import { getCompetitionById } from "../../utils";

const Competition = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const { competitionState, updateBasicInfo } = useContext(CompetitionContext) as CompetitionContextType

    useEffect(() => {
        getCompetitionById(id)
            .then(res => {
                if (res === undefined) navigate('/not-found');
                else updateBasicInfo(res);
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