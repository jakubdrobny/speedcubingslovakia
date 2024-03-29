import { Button, Card, Grid, Input } from "@mui/joy";
import { CompetitionContextType, ResultEntry } from "../../Types";
import { East, West } from "@mui/icons-material";

import { CompetitionContext } from "./CompetitionContext";
import Scramble from "./Scramble";
import TimeInput from "./TimeInput";
import { useContext } from "react";

const Compete = () => {
    const { competitionState, updateCurrentSolve, saveResults, updateSolve } = useContext(CompetitionContext) as CompetitionContextType

    return (
        <Card>
            <Grid container>
                <Grid xs={4}>
                    <Button
                        variant="outlined"
                        onClick={() => updateCurrentSolve((competitionState.currentSolveIdx - 1 + competitionState.noOfSolves) % competitionState.noOfSolves)}
                    >
                        <West />&nbsp;
                        Previous
                    </Button>
                </Grid>
                <Grid xs={4} sx={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                    Solve {competitionState.currentSolveIdx + 1}
                </Grid>
                <Grid xs={4} sx={{display: 'flex', justifyContent: 'flex-end'}}>
                    <Button
                        variant="outlined"
                        onClick={() => updateCurrentSolve((competitionState.currentSolveIdx + 1) % competitionState.noOfSolves)}
                    >
                        Next&nbsp;
                        <East />
                    </Button>
                </Grid>
            </Grid>
            <Scramble />
            <TimeInput />
            <Button color="primary" variant="solid" onClick={saveResults}>Save</Button>
        </Card>
    );
}

export default Compete;