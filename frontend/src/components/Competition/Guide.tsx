import { Alert, Card, Divider, Typography } from "@mui/joy";
import { CompetitionContextType, InputMethod } from "../../Types";
import { Keyboard, Timer } from "@mui/icons-material";

import { CompetitionContext } from "./CompetitionContext";
import { useContext } from "react";

const Guide = () => {
  const { competitionState } = useContext(
    CompetitionContext
  ) as CompetitionContextType;

  return (
    <Card style={{ paddingBottom: "1em" }}>
      <h3 style={{ textAlign: "center", marginBottom: 0 }}>
        How to submit results?
      </h3>
      {competitionState?.events[competitionState?.currentEventIdx]?.iconcode ===
      "fmc" ? (
        <div>
          <Typography>
            For FMC enter your solutions instead of times. They will be
            evaluated automatically.
          </Typography>
          <Typography>
            You can find the list of allowed moves{" "}
            <a href="https://www.worldcubeassociation.org/regulations/#12a">
              here
            </a>
            . (TLDR: basically anthing except slice moves)
          </Typography>
        </div>
      ) : competitionState.inputMethod === InputMethod.Manual ? (
        <div>
          <Typography>
            To enter your times, type just the numbers. For example, to enter 12
            seconds and 55 hundreths, type "1255".
          </Typography>
          <Typography>Penalties:</Typography>
          <ul>
            <li key={"3"}>
              If you get a +2, enter the final result. For example, if you
              finished in 12 second and 55 hundreths, with a +4 penalty, type
              "1655".
            </li>
            <li key={"4"}>
              If you get a DNF, type "d" and if you get a DNS, type "s".
            </li>
          </ul>
        </div>
      ) : (
        <div>
          <Typography>
            The timer is controlled using Spacebar. To start the solve, hold for
            1 second.
          </Typography>
          <Typography>
            After the solve, to add penalties, just click the corresponding
            buttons.
          </Typography>
          <ul>
            <li key={"5"}>
              You can add up to +16, after that, it will cycle back to no
              penalty.
            </li>
            <li key={"6"}>
              DNF can be removed by clicking the DNF button again.
            </li>
          </ul>
        </div>
      )}
      <Typography sx={{ paddingBottom: "1em" }}>
        After you are done, don't forget to save your results!
      </Typography>
      <Divider />
      <Alert color="warning">
        <Typography
          sx={{ display: "flex", alignItems: "center", fontWeight: "bold" }}
        >
          You can switch timing methods by clicking on Manual&nbsp;
          <Keyboard />
          &nbsp;/ Timer&nbsp; <Timer />.
        </Typography>
      </Alert>
    </Card>
  );
};

export default Guide;
