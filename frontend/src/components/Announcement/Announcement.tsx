import {
  AnnouncementState,
  AuthContextType,
  LoadingState,
  initialAnnouncementState,
} from "../../Types";
import { Card, Chip, CircularProgress, Stack, Typography } from "@mui/joy";
import { Delete, Edit } from "@mui/icons-material";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  ReadAnnouncement,
  getAnnouncementById,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";
import { useContext, useEffect, useRef, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import MarkdownPreview from "@uiw/react-markdown-preview";
import { Paper } from "@mui/material";
import emoji from "remark-emoji";
import {SlackSelector} from 'react-reactions'

const Announcement: React.FC<{
  givenAnnouncementState?: AnnouncementState;
  onAnnouncementDelete?: (idx: number, title: string, id: number) => void;
  idx?: number;
}> = ({ givenAnnouncementState, onAnnouncementDelete, idx }) => {
  const given = !isObjectEmpty(givenAnnouncementState || {});
  const navigate = useNavigate();
  const [announcementState, setAnnouncementState] = useState<AnnouncementState>(
    given
      ? (givenAnnouncementState as AnnouncementState)
      : initialAnnouncementState
  );
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const targetRef = useRef(null);
  const { authStateRef } = useContext(AuthContext) as AuthContextType;

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (
          authStateRef.current.token &&
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
      return;
    }

    setLoadingState({ isLoading: true, error: {} });

    getAnnouncementById(announcementState.id.toString())
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
          <Stack
            direction="row"
            alignItems="center"
            justifyContent="space-between"
            sx={{ borderBottom: "1px solid #CDD7E1" }}
          >
            <Stack spacing={1} direction="row" alignItems="center">
              {!announcementState.read && (
                <Chip variant="soft" color="danger" sx={{ height: "24px" }}>
                  New
                </Chip>
              )}
              <Typography level="h2">{announcementState.title}</Typography>
            </Stack>
            {authStateRef.current.isadmin && (
              <Stack direction="row" gap="10px">
                <Edit
                  color="primary"
                  onClick={() =>
                    navigate(`/announcement/${announcementState.id}/edit`)
                  }
                  sx={{ cursor: "pointer" }}
                  className="profile-cubing-icon-mock"
                />
                <Delete
                  color="error"
                  className="profile-cubing-icon-mock"
                  sx={{ cursor: "pointer" }}
                  onClick={() => {
                    if (onAnnouncementDelete !== undefined)
                      onAnnouncementDelete(
                        idx || 0,
                        announcementState.title,
                        parseInt(announcementState.id.toString() || "0")
                      );
                  }}
                />
              </Stack>
            )}
          </Stack>
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
            <Stack spacing={1} direction="row" flexWrap="wrap" useFlexGap>
              {announcementState.tags.map((tag, idx) => (
                <Chip key={idx} color={tag.color} sx={{ padding: "0 12px" }}>
                  {tag.label}
                </Chip>
              ))}
            </Stack>
          </Stack>
          <Paper elevation={3} sx={{ padding: "0.5em" }}>
            <div data-color-mode="light">
              <div className="wmde-markdown-var"> </div>
              <MarkdownPreview
                source={announcementState.content}
                style={{ padding: 16 }}
                remarkPlugins={[emoji]}
              />
            </div>
          </Paper>
          <SlackSele
        </Card>
      )}
    </div>
  );
};

export default Announcement;
