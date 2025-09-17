import { Button, ButtonGroup, Grid } from "@mui/joy";
import {
  CompetitionContextType,
  CompetitionEvent,
  ResultsCompeteChoiceEnum,
} from "../../Types";

import { CompetitionContext } from "../../context/CompetitionContext";
import { EmojiEvents } from "@mui/icons-material";
import { getCubingIconClassName } from "../../utils/utils";
import { useContext } from "react";
import { useSearchParams } from "react-router-dom";

export const EventSelector = () => {
  const {
    competitionState,
    updateCurrentEvent,
    loadingState,
    resultsCompeteChoice,
  } = useContext(CompetitionContext) as CompetitionContextType;
  const [searchParams, setSearchParams] = useSearchParams();

  let events = competitionState?.events;
  const shouldNotHaveOverall =
    events &&
    events.length > 0 &&
    events[0].id === -1 &&
    resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete;
  events = events.slice(+shouldNotHaveOverall);

  return (
    <Grid container>
      <ButtonGroup sx={{ py: 1, flexWrap: "wrap" }}>
        {events.map((e: CompetitionEvent, idx: number) => {
          const realIdx = idx + +shouldNotHaveOverall;
          return (
            <Button
              key={realIdx.toString() + e.iconcode}
              onClick={() => {
                updateCurrentEvent(realIdx);
                searchParams.set("event", e.iconcode);
                setSearchParams(searchParams);
              }}
              variant={
                realIdx === competitionState.currentEventIdx ? "solid" : "soft"
              }
              color="primary"
              disabled={loadingState.results}
            >
              {resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete ||
              e.id !== -1 ? (
                <span className={getCubingIconClassName(e.iconcode)}>
                  &nbsp;
                </span>
              ) : (
                <EmojiEvents />
              )}
              &nbsp;
              {e.displayname}
            </Button>
          );
        })}
      </ButtonGroup>
    </Grid>
  );
};
