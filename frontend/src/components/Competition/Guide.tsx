import { Card } from "@mui/joy";
import { CompetitionContext } from "./CompetitionContext";
import { CompetitionContextType } from "../../Types";
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
            ) : (
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
                        If you get a DNF, type "D" and if you get a DNS, type "S".
                        </li>
                    </ul>
                </div>
            )}
            <p>After you are done, don't forget to save your results!</p>
        </Card>
    );
}

export default Guide;