import { Button, ButtonGroup, Grid } from "@mui/joy";
import { CompetitionContextType, CompetitionEvent } from "../../Types";

import { CompetitionContext } from "./CompetitionContext";
import { useContext } from "react";

export const EventSelector = () => {
  const { competitionState, updateCurrentEvent } = useContext(
    CompetitionContext
  ) as CompetitionContextType;

  return (
    <Grid container>
      <ButtonGroup style={{ padding: "1em", flexWrap: "wrap" }}>
        {competitionState.events.map((e: CompetitionEvent, idx: number) => {
          return (
            <Button
              key={idx}
              onClick={() => updateCurrentEvent(idx)}
              variant={
                idx === competitionState.currentEventIdx ? "solid" : "soft"
              }
              color="primary"
              loading={competitionState.loadingState.results}
            >
              <span className={`cubing-icon event-${e.iconcode}`}>&ensp;</span>
              {e.displayname}
            </Button>
          );
        })}
      </ButtonGroup>
    </Grid>
  );
};
