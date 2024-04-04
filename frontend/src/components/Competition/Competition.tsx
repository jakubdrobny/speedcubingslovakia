import "../../styles/cubing-icons.css";

import { Alert, CircularProgress, Stack, Typography } from "@mui/joy";
import {
  AuthContextType,
  CompetitionContextType,
  CompetitionData,
  ResultEntry,
} from "../../Types";
import {
  formatDate,
  getCompetitionById,
  getResultsFromCompetitionAndEvent,
  initialCompetitionState,
} from "../../utils";
import { useContext, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { AuthContext } from "../../context/AuthContext";
import { CompetitionContext } from "./CompetitionContext";
import CompetitorArea from "./CompetitorArea";
import { EventSelector } from "./EventSelector";

const Competition = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { competitionState, setCompetitionState, updateBasicInfo } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const { authState } = useContext(AuthContext) as AuthContextType;

  useEffect(() => {
    setCompetitionState({
      ...competitionState,
      loadingState: { compinfo: true, error: "" },
    });

    getCompetitionById(id)
      .then((info: CompetitionData | undefined) => {
        setCompetitionState({
          ...competitionState,
          loadingState: { ...competitionState.loadingState, compinfo: false },
        });

        if (info === undefined) navigate("/not-found");
        else {
          updateBasicInfo(info);
        }
      })
      .catch((err) =>
        setCompetitionState({
          ...competitionState,
          loadingState: { compinfo: false, error: err.message },
        })
      );

    return () => setCompetitionState(initialCompetitionState);
  }, []);

  return (
    <Stack
      spacing={3}
      sx={{ display: "flex", alignItems: "center", margin: "2em 0" }}
    >
      {competitionState.loadingState.error ? (
        <Alert color="danger">{competitionState.loadingState.error}</Alert>
      ) : competitionState.loadingState.compinfo ? (
        <CircularProgress />
      ) : (
        <>
          <Typography level="h1">{competitionState.name}</Typography>
          <Typography>
            {formatDate(competitionState.startdate)} -{" "}
            {formatDate(competitionState.enddate)}
          </Typography>
          <EventSelector />
          <br />
          <CompetitorArea />
        </>
      )}
    </Stack>
  );
};

export default Competition;
