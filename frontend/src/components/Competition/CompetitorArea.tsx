import { Alert, Grid, Stack } from "@mui/joy";

import { AuthContext } from "../../context/AuthContext";
import { AuthContextType } from "../../Types";
import CommentBox from "./CommentBox";
import Compete from "./Compete";
import Guide from "./Guide";
import { saveCurrentLocation } from "../../utils/utils";
import { useContext } from "react";

const CompetitorArea: React.FC<{ loading: boolean }> = ({ loading }) => {
  const { authState } = useContext(AuthContext) as AuthContextType;

  return (
    <Grid container style={{ width: "100%" }} spacing={2}>
      {!authState.token ? (
        <Alert color="warning" sx={{ width: "100%", margin: "0 1em", gap: 0 }}>
          <Grid container direction="row" alignItems="center">
            You are not authenticated! Please{" "}
            <span style={{ fontSize: "1em" }}>&nbsp;</span>
            <a
              href={import.meta.env.VITE_WCA_GET_CODE_URL || ""}
              onClick={() => saveCurrentLocation(window.location.pathname)}
            >
              log in
            </a>
            <span style={{ fontSize: "1em" }}>&nbsp;</span> to compete.
          </Grid>
        </Alert>
      ) : (
        <>
          <Grid xs={0} sx={{ padding: 0, margin: 0 }}></Grid>
          <Grid xs={12} md={6}>
            <Compete />
          </Grid>
          <Grid xs={0} sx={{ padding: 0, margin: 0 }}></Grid>
          <Grid xs={0} sx={{ padding: 0, margin: 0 }}></Grid>
          <Grid xs={12} md={6}>
            <Stack spacing={2}>
              <Guide />
              <CommentBox disabled={loading} />
            </Stack>
          </Grid>
        </>
      )}
    </Grid>
  );
};

export default CompetitorArea;
