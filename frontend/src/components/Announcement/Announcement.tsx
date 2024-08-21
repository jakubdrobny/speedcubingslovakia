import {
  AnnouncementState,
  LoadingState,
  initialAnnouncementState,
} from "../../Types";
import { Box, Card, Chip, CircularProgress, Stack, Typography } from "@mui/joy";
import { Link, useParams } from "react-router-dom";
import { getAnnouncementById, getError } from "../../utils";
import { useEffect, useState } from "react";

import Markdown from "react-markdown";
import { Paper } from "@mui/material";

const Announcement = () => {
  const { id } = useParams<{ id: string }>();
  const [announcementState, setAnnouncementState] = useState<AnnouncementState>(
    initialAnnouncementState
  );
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    getAnnouncementById(id)
      .then((res) => {
        setAnnouncementState(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) =>
        setLoadingState({ isLoading: false, error: getError(err) })
      );
  }, []);

  return (
    <div style={{ margin: "2em 0.5em" }}>
      {loadingState.isLoading ? (
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <CircularProgress />
          &nbsp; &nbsp;{" "}
          <Typography level="h3">Loading announcement...</Typography>
        </div>
      ) : (
        <Card>
          <Typography level="h2" sx={{ borderBottom: "1px solid #CDD7E1" }}>
            {announcementState.title}
          </Typography>
          <Stack spacing={1} direction="row">
            <div>author:</div>
            <Link
              to={`/profile/${announcementState.authorWcaId}`}
              style={{
                color: "#0B6BCB",
                textDecoration: "none",
                fontWeight: 555,
              }}
            >
              {announcementState.authorUsername}
            </Link>
          </Stack>
          <Stack spacing={1} direction="row">
            <div>tags:</div>
            <Stack spacing={1}>
              {announcementState.tags.map((tag) => (
                <Chip color={tag.color} sx={{ padding: "0 12px" }}>
                  {tag.label}
                </Chip>
              ))}
            </Stack>
          </Stack>
          <Paper elevation={3} sx={{ padding: "0.5em" }}>
            <Markdown>{announcementState.content}</Markdown>
          </Paper>
        </Card>
      )}
    </div>
  );
};

export default Announcement;
