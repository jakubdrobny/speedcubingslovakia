import { CompetitionContextType, InputMethod } from "../../Types";

import { Card } from "@mui/joy";
import { CompetitionContext } from "./CompetitionContext";
import { useContext } from "react";

const Guide = () => {
    const { competitionState } = useContext(CompetitionContext) as CompetitionContextType

    return (
        <Card>
            <h3 style={{ textAlign: 'center', marginBottom: 0 }}>How to submit results?</h3>
            {competitionState?.events[competitionState?.currentEventIdx]?.iconcode === "fmc" ? (
                <div>
                    <p>
                        For FMC enter your solutions instead of times. They will be evaluated automatically.
                    </p>
                    <p>
                        You can find the list of allowed moves <a href="https://www.worldcubeassociation.org/regulations/#12a">here</a>.
                        (TLDR: basically anthing except slice moves)
                    </p>
                </div>
            ) : competitionState.inputMethod === InputMethod.Manual ? (
                <div>
                    <p>
                        To enter your times, type just the numbers. For example, to enter 12 seconds and 55 hundreths, type "1255".
                    </p>
                    <p>
                        Penalties:
                    </p>
                    <ul>
                        <li key={"3"}>
                            If you get a +2, enter the final result. For example, if you finished
                            in 12 second and 55 hundreths, with a +4 penalty, type "1655".
                        </li>
                        <li key={"4"}>
                        If you get a DNF, type "d" and if you get a DNS, type "s".
                        </li>
                    </ul>
                </div>
            ) : 
                <div>
                    <p>
                        The timer is controlled using Spacebar. To start the solve, hold for 1 second.
                    </p>
                    <p>
                        After the solve, to add penalties, just click the corresponding buttons.
                    </p>
                    <ul>
                        <li key={"5"}>
                            You can add up to +16, after that, it will cycle back to no penalty.
                        </li>
                        <li key={"6"}>
                            DNF can be removed by clicking the DNF button again.
                        </li>
                    </ul>
                </div>
            }
            <p>After you are done, don't forget to save your results!</p>
        </Card>
    );
}

export default Guide;