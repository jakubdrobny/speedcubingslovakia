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
import {
  CompetitionContextType,
  CompetitionData,
  ResultsCompeteChoiceEnum,
} from "../../Types";
import {
  formatDate,
  getCompetitionById,
  initialCompetitionState,
} from "../../utils";
import { useContext, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { CompetitionContext } from "./CompetitionContext";
import CompetitionResults from "./CompetitionResults";
import CompetitorArea from "./CompetitorArea";
import { EventSelector } from "./EventSelector";
import ResultsCompeteChoice from "./ResultsCompeteChoice";
import { Warning } from "@mui/icons-material";

const Competition = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const {
    suspicousModalOpen,
    setSuspicousModalOpen,
    competitionState,
    setCompetitionState,
    updateBasicInfo,
    loadingState,
    setLoadingState,
    resultsCompeteChoice,
    setResultsCompeteChoice,
  } = useContext(CompetitionContext) as CompetitionContextType;

  useEffect(() => {
    setLoadingState({ results: false, compinfo: true, error: "" });

    getCompetitionById(id)
      .then((info: CompetitionData | undefined) => {
        setLoadingState({ ...loadingState, compinfo: false });

        if (info === undefined) navigate("/not-found");
        else {
          updateBasicInfo(info);
        }
      })
      .catch((err) =>
        setLoadingState({
          ...loadingState,
          compinfo: false,
          error: err.message,
        })
      );

    return () => setCompetitionState(initialCompetitionState);
  }, []);

  return (
    <Stack spacing={3} sx={{ display: "flex", alignItems: "center", mt: 2 }}>
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
            If you achieved these results legitimately, please let us know in
            the comment box and your results will be approved.
          </DialogContent>
        </ModalDialog>
      </Modal>
      {loadingState.error ? (
        <Alert color="danger">{loadingState.error}</Alert>
      ) : loadingState.compinfo ? (
        <CircularProgress />
      ) : (
        <>
          <Typography level="h1">{competitionState.name}</Typography>
          <Typography>
            {formatDate(competitionState.startdate)} -{" "}
            {formatDate(competitionState.enddate)}
          </Typography>
          <EventSelector />
          <ResultsCompeteChoice
            resultsCompeteChoice={resultsCompeteChoice}
            setResultsCompeteChoice={setResultsCompeteChoice}
            loading={loadingState.results}
          />
          {resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete ? (
            <CompetitorArea loading={loadingState.results} />
          ) : (
            <CompetitionResults />
          )}
        </>
      )}
    </Stack>
  );
};

export default Competition;
