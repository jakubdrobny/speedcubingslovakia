import "../../styles/cubing-icons.css";

import {
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
  getError,
  initialCompetitionState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import { useContext, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { CompetitionContext } from "../../context/CompetitionContext";
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
    warningModalOpen,
    setWarningModalOpen,
    competitionState,
    setCompetitionState,
    updateBasicInfo,
    loadingState,
    setLoadingState,
    resultsCompeteChoice,
    setResultsCompeteChoice,
  } = useContext(CompetitionContext) as CompetitionContextType;

  useEffect(() => {
    setLoadingState({ results: false, compinfo: true, error: {} });

    getCompetitionById(id)
      .then((info: CompetitionData | undefined) => {
        setLoadingState({ ...loadingState, compinfo: false });

        if (info === undefined) navigate("/not-found");
        else {
          updateBasicInfo(info);
        }
      })
      .catch((err) => {
        setLoadingState({
          ...loadingState,
          compinfo: false,
          error: getError(err),
        });
      });

    return () => setCompetitionState(initialCompetitionState);
  }, []);

  return (
    <Stack
      spacing={3}
      sx={{ display: "flex", alignItems: "center", mt: 4, mx: 2 }}
    >
      <Modal
        open={suspicousModalOpen || warningModalOpen}
        onClose={() => {
          if (suspicousModalOpen) setSuspicousModalOpen(false);
          else setWarningModalOpen(false);
        }}
      >
        <ModalDialog
          color={
            suspicousModalOpen
              ? "danger"
              : warningModalOpen
                ? "warning"
                : "neutral"
          }
          layout="center"
          size="lg"
          variant="soft"
          role="alertdialog"
        >
          <DialogTitle sx={{ display: "flex", alignItems: "center" }}>
            <Warning />
            {suspicousModalOpen
              ? "Suspicous result detected"
              : "Incorrect results entered."}
          </DialogTitle>
          <ModalClose />
          <Divider />
          <DialogContent>
            {suspicousModalOpen ? (
              <>
                Your results were identified as suspicous, which means you
                likely made a data entry error. If that's the case, please fix
                it, until then, these results won't show up on the leaderboard.
                <br />
                <br />
                If you achieved these results legitimately, please let us know
                in the comment box and your results will be approved.
              </>
            ) : (
              <>
                The results you entered had invalid format, eg. over 60 seconds
                or over 60 minutes, etc., which means you likely made a data
                entry error.
                <br />
                <br />
                The incorrect solves were saved as DNF, so please re-enter the
                result in the correct format.
              </>
            )}
          </DialogContent>
        </ModalDialog>
      </Modal>
      {!isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : loadingState.compinfo ? (
        <CircularProgress />
      ) : (
        <>
          <Typography level="h1" sx={{ textAlign: "center" }}>
            {competitionState.name}
          </Typography>
          <Typography sx={{ textAlign: "center" }}>
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
