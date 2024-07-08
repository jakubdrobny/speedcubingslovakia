import { Button, ButtonGroup, Grid } from "@mui/joy";
import {
  CompetitionContextType,
  CompetitionEvent,
  ResultsCompeteChoiceEnum,
} from "../../Types";

import { CompetitionContext } from "./CompetitionContext";
import { EmojiEvents } from "@mui/icons-material";
import { useContext } from "react";

export const EventSelector = () => {
  const {
    competitionState,
    updateCurrentEvent,
    loadingState,
    resultsCompeteChoice,
  } = useContext(CompetitionContext) as CompetitionContextType;

  let events = competitionState?.events;
  const shouldNotHaveOverall =
    events &&
    events.length > 0 &&
    events[events.length - 1].id === -1 &&
    resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete;
  if (shouldNotHaveOverall) events = events.slice(0, events.length - 1);

  return (
    <Grid container>
      <ButtonGroup sx={{ py: 1, flexWrap: "wrap" }}>
        {events.map((e: CompetitionEvent, idx: number) => {
          return (
            <Button
              key={idx}
              onClick={() => updateCurrentEvent(idx)}
              variant={
                idx === competitionState.currentEventIdx ? "solid" : "soft"
              }
              color="primary"
              loading={loadingState.results}
            >
              {resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete ||
              idx < events.length - 1 ? (
                <span className={`cubing-icon event-${e.iconcode}`}>
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
