import { Button, ButtonGroup } from "@mui/joy";
import {
  CompetitionContextType,
  InputMethod,
  TimerInputContextType,
  TimerInputCurrentState,
} from "../../Types";

import { CompetitionContext } from "../../context/CompetitionContext";
import { TimerInputContext } from "../../context/TimerInputContext";
import { competitionOnGoing } from "../../utils";
import { useContext } from "react";

const Penalties = () => {
  const { timerInputState } = useContext(
    TimerInputContext
  ) as TimerInputContextType;
  const { competitionState, addPenalty } = useContext(
    CompetitionContext
  ) as CompetitionContextType;

  return (
    <div
      style={
        timerInputState.currentState === TimerInputCurrentState.NotSolving &&
        competitionState.inputMethod === InputMethod.Timer
          ? {
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              width: "100%",
              margin: "1em 0 1em 0",
            }
          : { display: "none" }
      }
    >
      <ButtonGroup disabled={!competitionOnGoing(competitionState)}>
        <Button
          color="primary"
          variant="outlined"
          onClick={() => addPenalty("2")}
        >
          +2
        </Button>
        <Button
          color="primary"
          variant="outlined"
          onClick={() => addPenalty("DNF")}
        >
          DNF
        </Button>
      </ButtonGroup>
    </div>
  );
};

export default Penalties;
