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
  CompetitionEvent,
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
import { useNavigate, useParams, useSearchParams } from "react-router-dom";

import { CompetitionContext } from "../../context/CompetitionContext";
import CompetitionResults from "./CompetitionResults";
import CompetitorArea from "./CompetitorArea";
import { EventSelector } from "./EventSelector";
import ResultsCompeteChoice from "./ResultsCompeteChoice";
import { Warning } from "@mui/icons-material";
import { RESULTS_COMPETE_CHOICE_QUERY_PARAM_NAME } from "../../constants";

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
  const [searchParams, setSearchParams] = useSearchParams();

  // returns the index of the event from the events
  const locateEventFromURLQuery = (
    eventQueryParam: string,
    events: CompetitionEvent[],
  ): number => {
    for (let i = 0; i < events.length; i++) {
      if (eventQueryParam === events[i].iconcode) {
        return i;
      }
    }

    searchParams.set("event", "overall");
    setSearchParams(searchParams);
    return -1;
  };

  const handleResultsCompeteChoiceQueryParam = () => {
    let currentParamValue =
      searchParams.get(RESULTS_COMPETE_CHOICE_QUERY_PARAM_NAME) || "";
    if (!["compete", "results"].includes(currentParamValue)) {
      currentParamValue = "results";
      searchParams.set(
        RESULTS_COMPETE_CHOICE_QUERY_PARAM_NAME,
        currentParamValue,
      );
      setSearchParams(searchParams);
    }

    setResultsCompeteChoice(
      currentParamValue === "results"
        ? ResultsCompeteChoiceEnum.Results
        : ResultsCompeteChoiceEnum.Compete,
    );
  };

  useEffect(() => {
    setLoadingState({ results: false, compinfo: true, error: {} });

    getCompetitionById(id)
      .then((info: CompetitionData | undefined) => {
        setLoadingState({ ...loadingState, compinfo: false });

        if (info === undefined) navigate("/not-found");
        else {
          let eventIdx = locateEventFromURLQuery(
            searchParams.get("event") || "",
            info.events,
          );
          eventIdx =
            eventIdx < 0 || eventIdx >= info.events.length
              ? info.events.length - 1
              : eventIdx;
          updateBasicInfo(info, eventIdx);
          handleResultsCompeteChoiceQueryParam();
        }
      })
      .catch((err) => {
        setLoadingState({
          ...loadingState,
          compinfo: false,
          error: getError(err),
        });
      });

    return () => setCompetitionState({ ...initialCompetitionState });
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
