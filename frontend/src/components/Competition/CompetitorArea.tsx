import Compete from "./Compete";
import { Grid } from "@mui/joy";
import Guide from "./Guide";

const CompetitorArea = () => {
    return (
        <Grid container style={{ width: '100%' }} spacing={2}>
            <Grid xs={1} md={0} sx={{padding: 0, margin: 0}}></Grid>
            <Grid xs={10} md={6}>
               <Compete />
            </Grid>
            <Grid xs={1} md={0} sx={{padding: 0, margin: 0}}></Grid>
            <Grid xs={1} md={0} sx={{padding: 0, margin: 0}}></Grid>
            <Grid xs={10} md={6}>
                <Guide />
            </Grid>
        </Grid>
    );
}

export default CompetitorArea;