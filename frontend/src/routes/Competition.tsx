import { Button, ButtonGroup } from "@mui/joy";
import { CompetitionData, CompetitionEvent } from "../Types";
import React, { useEffect } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";

import { getCompetitionById } from "../utils";
import { useState } from "react";

const Competition = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const [competitionData, setCompetitionData] = useState<CompetitionData>();
    const [currentEvent, setCurrentEvent] = useState<number>();

    useEffect(() => {
        getCompetitionById(id)
            .then(res => {
                if (res === undefined) navigate('/not-found');
                else setCompetitionData(res);
            })
            .catch(console.error);
    }, []);

    return (
        <div style={{display: 'flex', flexDirection: 'column', alignItems: 'center'}}>
            <h1>{competitionData?.name}</h1>
            <p>{competitionData?.startdate.toLocaleString()} - {competitionData?.enddate.toLocaleString()}</p>
            <ButtonGroup>
                {competitionData?.events.map((e: CompetitionEvent, idx: number) => {
                    return (
                        <Button 
                            key={idx} 
                            onClick={() => setCurrentEvent(idx)}
                            variant={idx === currentEvent ? "solid" : "soft"}
                            color="primary"
                        >
                            <span className={`cubing-icon event-${e.iconcode}`}></span>
                            {e.displayname}
                        </Button>
                    )
                })}
            </ButtonGroup>
        </div>
    );
};

export default Competition;