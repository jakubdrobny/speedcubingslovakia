import "../../styles/cubing-icons.css";

import {
  Alert,
  CircularProgress,
  DialogContent,
  DialogTitle,
  Divider,
  Modal,
  ModalClose,
  ModalDialog,
  Stack,
  Typography,
} from "@mui/joy";
import { CompetitionContextType, CompetitionData } from "../../Types";
import { PriorityHigh, Warning } from "@mui/icons-material";
import {
  formatDate,
  getCompetitionById,
  initialCompetitionState,
} from "../../utils";
import { useContext, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { CompetitionContext } from "./CompetitionContext";
import CompetitorArea from "./CompetitorArea";
import { EventSelector } from "./EventSelector";

const Competition = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const {
    suspicousModalOpen,
    setSuspicousModalOpen,
    competitionState,
    setCompetitionState,
    updateBasicInfo,
  } = useContext(CompetitionContext) as CompetitionContextType;

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
      <Modal
        open={suspicousModalOpen}
        onClose={() => setSuspicousModalOpen(false)}
      >
        <ModalDialog
          color="danger"
          layout="center"
          size="lg"
          variant="soft"
          role="alertdialog"
        >
          <DialogTitle sx={{ display: "flex", alignItems: "center" }}>
            <Warning />
            Suspicous result detected
          </DialogTitle>
          <ModalClose />
          <Divider />
          <DialogContent>
            Your results were identified as suspicous, which means you likely
            made a data entry error. If that's the case, please fix it, until
            then, these results won't show up on the leaderboard.
            <br />
            <br />
            If you achieved these results legitimately, please let us know and
            your results will be approved.
          </DialogContent>
        </ModalDialog>
      </Modal>
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
