import { Button, Card, Grid, Typography } from "@mui/joy";
import { CompetitionContextType, InputMethod } from "../../Types";
import { East, Keyboard, Timer, West } from "@mui/icons-material";

import { CompetitionContext } from "./CompetitionContext";
import ManualInput from "./ManualInput";
import Scramble from "./Scramble";
import TimerInput from "./TimerInput";
import { useContext } from "react";

const Compete = () => {
    const { competitionState, updateCurrentSolve, toggleInputMethod, saveResults } = useContext(CompetitionContext) as CompetitionContextType

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
            <Grid container>
                <Grid xs={12} sx={{display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                    <h3
                        onClick={toggleInputMethod}
                        style={{display: 'flex', alignItems: 'center'}}
                    >
                        {competitionState.inputMethod === InputMethod.Manual ? 
                            <>
                                Manual&nbsp;
                                <Keyboard />
                            </>
                        :
                            <>
                                Timer&nbsp;
                                <Timer />
                            </>
                        }
                    </h3>
                </Grid>
                <Grid xs={12}>
                    {competitionState.inputMethod === InputMethod.Manual ? <ManualInput /> : <TimerInput />}
                </Grid>
                <Grid xs={12} sx={{marginTop: 2}}>
                    <Button color="primary" variant="solid" onClick={saveResults} sx={{width: '100%'}}>Save</Button>
                </Grid>
            </Grid>
        </Card>
    );
}

export default Compete;