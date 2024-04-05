import { Alert, Grid, Stack } from "@mui/joy";

import { AuthContext } from "../../context/AuthContext";
import { AuthContextType } from "../../Types";
import CommentBox from "./CommentBox";
import Compete from "./Compete";
import Guide from "./Guide";
import { useContext } from "react";

const CompetitorArea = () => {
  const { authState } = useContext(AuthContext) as AuthContextType;

  return (
    <Grid container style={{ width: "100%" }} spacing={2}>
      {!authState.authenticated ? (
        <Alert color="warning" sx={{ width: "100%", margin: "0 1em" }}>
          You are not authenticated! Please log in to compete.
        </Alert>
      ) : (
        <>
          <Grid xs={1} md={0} sx={{ padding: 0, margin: 0 }}></Grid>
          <Grid xs={10} md={6}>
            <Compete />
          </Grid>
          <Grid xs={1} md={0} sx={{ padding: 0, margin: 0 }}></Grid>
          <Grid xs={1} md={0} sx={{ padding: 0, margin: 0 }}></Grid>
          <Grid xs={10} md={6}>
            <Stack spacing={2}>
              <Guide />
              <CommentBox />
            </Stack>
          </Grid>
        </>
      )}
    </Grid>
  );
};

export default CompetitorArea;
