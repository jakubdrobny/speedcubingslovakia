import { Button, Card, Grid, Input } from "@mui/joy";
import { East, West } from "@mui/icons-material";

import { CompetitionContext } from "./CompetitionContext";
import { CompetitionContextType } from "../../Types";
import { useContext } from "react";

const Compete = () => {
    const { competitionState, updateCurrentSolve } = useContext(CompetitionContext) as CompetitionContextType

    return (
        <Card>
            <Grid container>
                <Grid xs={4}>
                    <Button
                        variant="outlined"
                        disabled={competitionState.currentSolveIdx === 0}
                        onClick={() => updateCurrentSolve(competitionState.currentSolveIdx - 1)}
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
                        disabled={competitionState.currentSolveIdx === competitionState.noOfSolves - 1}
                        onClick={() => updateCurrentSolve(competitionState.currentSolveIdx + 1)}
                    >
                        Next&nbsp;
                        <East />
                    </Button>
                </Grid>
            </Grid>
            <Input size="lg" placeholder="Large" />
        </Card>
    );
}

export default Compete;