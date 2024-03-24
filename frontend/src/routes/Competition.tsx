import React, { useEffect } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";

import { CompetitionData } from "../Types";
import { Typography } from "@mui/joy";
import { getCompetitionById } from "../utils";
import { useState } from "react";

const Competition = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const [competitionData, setCompetitionData] = useState<CompetitionData>();

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
        </div>
    );
};

export default Competition;