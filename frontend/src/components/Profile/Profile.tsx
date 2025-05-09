import { CircularProgress, Grid, Stack, Typography } from "@mui/joy";
import { LoadingState, ProfileType } from "../../Types";
import {
  defaultProfile,
  getError,
  getProfile,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import { useEffect, useState } from "react";

import MedalRecordColletion from "./MedalRecordColletion";
import ProfileBasics from "./ProfileBasics";
import ProfilePersonalBests from "./ProfilePersonalBests";
import ProfileResultsHistory from "./ProfileResultsHistory";
import { WIN_SMALL } from "../../constants";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { WindowSizeContextType } from "../../Types";
import { useContext } from "react";
import { useParams } from "react-router-dom";

const Profile = () => {
  const { id = "trolko" } = useParams<{ id: string }>();
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const [profile, setProfile] = useState<ProfileType>(defaultProfile);
  const { windowSize } = useContext(
    WindowSizeContext
  ) as WindowSizeContextType;

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} })

    getProfile(id)
      .then((p: ProfileType) => {
        setProfile(p);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, [id]);

  return (
    <div style={{ margin: "2em" }}>
      {loadingState.isLoading ? (
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <CircularProgress />
          &nbsp; &nbsp; <Typography level="h3">Loading profile...</Typography>
        </div>
      ) : !isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <Stack spacing={3}>
          <ProfileBasics basics={profile.basics} />
          <ProfilePersonalBests pbs={profile.personalBests} />
          <Grid container>
            <Grid
              xs={12}
              md={6}
              style={
                windowSize.width >= WIN_SMALL ? { paddingRight: "0.5em" } : {}
              }
            >
              <MedalRecordColletion
                title="Medal Collection"
                headers={["Gold", "Silver", "Bronze"]}
                values={[
                  profile.medalCollection.gold,
                  profile.medalCollection.silver,
                  profile.medalCollection.bronze,
                ]}
              />
            </Grid>
            <Grid
              xs={12}
              md={6}
              style={
                windowSize.width >= WIN_SMALL
                  ? { paddingLeft: "0.5em" }
                  : { paddingTop: "1em" }
              }
            >
              <MedalRecordColletion
                title="Record Collection"
                headers={["WR", "CR", "NR"]}
                values={[
                  profile.recordCollection.wr,
                  profile.recordCollection.cr,
                  profile.recordCollection.nr,
                ]}
              />
            </Grid>
          </Grid>
          {profile.resultsHistory && profile.resultsHistory.length > 0 && (
            <ProfileResultsHistory resultsHistory={profile.resultsHistory} />
          )}
        </Stack>
      )}
    </div>
  );
};

export default Profile;
