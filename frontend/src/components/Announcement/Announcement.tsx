import {
  AddReactionToAnnouncement,
  ReadAnnouncement,
  getAnnouncementById,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";
import {
  AnnouncementReactResponse,
  AnnouncementState,
  AuthContextType,
  LoadingState,
  initialAnnouncementState,
} from "../../Types";
import { Card, Chip, CircularProgress, Stack, Typography } from "@mui/joy";
import { Delete, Edit } from "@mui/icons-material";
import { Link, useNavigate, useParams } from "react-router-dom";
import { SlackCounter, SlackSelector } from "@charkour/react-reactions";
import { useContext, useEffect, useRef } from "react";

import { AuthContext } from "../../context/AuthContext";
import MarkdownPreview from "@uiw/react-markdown-preview";
import { Paper } from "@mui/material";
import emoji from "remark-emoji";
import useState from "react-usestateref";

const Announcement: React.FC<{
  givenAnnouncementState?: AnnouncementState;
  onAnnouncementDelete?: (idx: number, title: string, id: number) => void;
  idx?: number;
}> = ({ givenAnnouncementState, onAnnouncementDelete, idx }) => {
  const given = !isObjectEmpty(givenAnnouncementState || {});
  const navigate = useNavigate();
  const [announcementState, setAnnouncementState, announcementStateRef] =
    useState<AnnouncementState>(
      given
        ? (givenAnnouncementState as AnnouncementState)
        : initialAnnouncementState
    );
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const { id } = useParams<{ id: string }>();
  const targetRef = useRef(null);
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const [emojiSelectorOpen, setEmojiSelectorOpen] = useState<boolean>(false);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    if (!given) {
      setAnnouncementState({ ...announcementState, id: parseInt(id || "0") });

      getAnnouncementById(announcementStateRef.current.id.toString())
        .then((res) => {
          setAnnouncementState(res);
          if (!res.read) return ReadAnnouncement(res);
        })
        .then((res) => setLoadingState({ isLoading: false, error: {} }))
        .catch((err) =>
          setLoadingState({ isLoading: false, error: getError(err) })
        );
    }

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

    setLoadingState({ isLoading: false, error: {} });

    return () => {
      if (targetRef.current) {
        observer.unobserve(targetRef.current);
      }
    };
  }, []);

  const handleOnReactionSelect = (emoji: string) => {
    const by = authStateRef.current.username;
    AddReactionToAnnouncement(announcementStateRef.current.id, emoji, by)
      .then((res: AnnouncementReactResponse) => {
        setAnnouncementState({
          ...announcementState,
          emojiCounters: res.set
            ? [...announcementStateRef.current.emojiCounters, { emoji, by }]
            : [...announcementStateRef.current.emojiCounters].filter(
                (entry) => !(entry.emoji === emoji && entry.by === by)
              ),
        });
      })
      .catch((err) =>
        setLoadingState({
          isLoading: loadingState.isLoading,
          error: getError(err),
        })
      );
  };

  return (
    <Stack style={{ margin: "0.5em", height: "100%" }} spacing={2}>
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
        announcementState.id !== 0 && (
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
            <SlackCounter
              counters={announcementStateRef.current.emojiCounters}
              onSelect={(emoji) => handleOnReactionSelect(emoji)}
              onAdd={() => {
                if (authStateRef.current.token) setEmojiSelectorOpen((p) => !p);
              }}
            />
            {emojiSelectorOpen && authStateRef.current.token && (
              <SlackSelector
                onSelect={(emoji) => handleOnReactionSelect(emoji)}
              />
            )}
          </Card>
        )
      )}
    </Stack>
  );
};

export default Announcement;
