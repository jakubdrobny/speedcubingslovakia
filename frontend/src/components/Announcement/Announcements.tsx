import { AnnouncementState, LoadingState } from "../../Types";
import { CircularProgress, Stack, Typography } from "@mui/joy";
import {
  GetAnnouncements,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";
import { useEffect, useState } from "react";

import Announcement from "./Announcement";

const Announcements = () => {
  const [announcements, setAnnouncements] = useState<AnnouncementState[]>([]);
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    GetAnnouncements()
      .then((res) => {
        setAnnouncements(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) =>
        setLoadingState({ isLoading: false, error: getError(err) })
      );
  }, []);

  useEffect(() => {}, []);

  return (
    <div style={{ margin: "0.5em", height: "100%" }}>
      {!isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : loadingState.isLoading ? (
        <Typography
          level="h3"
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "100%",
          }}
        >
          <CircularProgress /> &nbsp; Loading announcements...
        </Typography>
      ) : (
        <div>
          <Typography
            level="h1"
            sx={{ margin: "0.5em 0.5em", borderBottom: "1px solid #CDD7E1" }}
          >
            Announcements
          </Typography>
          <Stack spacing={2} direction="column">
            {announcements.map((announcement, idx) => (
              <Announcement key={idx} givenAnnouncementState={announcement} />
            ))}
          </Stack>
        </div>
      )}
    </div>
  );
};

export default Announcements;
