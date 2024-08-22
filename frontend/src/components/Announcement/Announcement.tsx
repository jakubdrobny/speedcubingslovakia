import {
  AnnouncementState,
  AuthContextType,
  LoadingState,
  initialAnnouncementState,
} from "../../Types";
import { Card, Chip, CircularProgress, Stack, Typography } from "@mui/joy";
import { Link, useParams } from "react-router-dom";
import {
  ReadAnnouncement,
  getAnnouncementById,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";
import { useContext, useEffect, useRef, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import Markdown from "react-markdown";
import { Paper } from "@mui/material";

const Announcement: React.FC<{
  givenAnnouncementState?: AnnouncementState;
}> = ({ givenAnnouncementState }) => {
  const given = !isObjectEmpty(givenAnnouncementState || {});
  const [announcementState, setAnnouncementState] = useState<AnnouncementState>(
    given
      ? (givenAnnouncementState as AnnouncementState)
      : initialAnnouncementState
  );
  let { id } = useParams<{ id: string }>();
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const targetRef = useRef(null);
  const { authState } = useContext(AuthContext) as AuthContextType;

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        console.log(
          authState.token,
          entry.isIntersecting,
          announcementState.read
        );
        if (
          authState.token &&
          entry.isIntersecting &&
          !announcementState.read
        ) {
          ReadAnnouncement(announcementState)
            .then((res) =>
              setAnnouncementState({ ...announcementState, read: true })
            )
            .catch((err) =>
              setLoadingState({
                isLoading: loadingState.isLoading,
                error: getError(err),
              })
            );
        }
      },
      {
        root: null,
        rootMargin: "0px",
        threshold: 0.5,
      }
    );

    if (targetRef.current) {
      observer.observe(targetRef.current);
    }

    if (given) {
      id = announcementState.id.toString();
      return;
    }

    setLoadingState({ isLoading: true, error: {} });

    getAnnouncementById(id)
      .then((res) => {
        setAnnouncementState(res);
        if (!res.read) return ReadAnnouncement(res);
      })
      .then((res) => setLoadingState({ isLoading: false, error: {} }))
      .catch((err) =>
        setLoadingState({ isLoading: false, error: getError(err) })
      );

    return () => {
      if (targetRef.current) {
        observer.unobserve(targetRef.current);
      }
    };
  }, []);

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
          <CircularProgress /> &nbsp; Loading announcement...
        </Typography>
      ) : (
        <Card ref={targetRef}>
          <Typography level="h2" sx={{ borderBottom: "1px solid #CDD7E1" }}>
            {announcementState.title}
          </Typography>
          {!announcementState.read && (
            <Chip variant="soft" color="danger">
              New
            </Chip>
          )}
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
              {announcementState.tags.map((tag, idx) => (
                <Chip key={idx} color={tag.color} sx={{ padding: "0 12px" }}>
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
