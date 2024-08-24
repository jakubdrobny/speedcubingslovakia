import {
  AnnouncementState,
  AuthContextType,
  LoadingState,
  initialAnnouncementState,
} from "../../Types";
import {
  Button,
  ButtonGroup,
  Card,
  Chip,
  CircularProgress,
  DialogTitle,
  Divider,
  Modal,
  ModalClose,
  ModalDialog,
  Stack,
  Typography,
} from "@mui/joy";
import { Delete, Edit, Warning } from "@mui/icons-material";
import {
  DeleteAnnouncement,
  ReadAnnouncement,
  getAnnouncementById,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";
import { DialogContent, Paper } from "@mui/material";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useContext, useEffect, useRef, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import MarkdownPreview from "@uiw/react-markdown-preview";

const Announcement: React.FC<{
  givenAnnouncementState?: AnnouncementState;
}> = ({ givenAnnouncementState }) => {
  const given = !isObjectEmpty(givenAnnouncementState || {});
  const navigate = useNavigate();
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
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const [deletingAnnouncement, setDeletingAnnouncement] =
    useState<boolean>(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState<boolean>(false);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (
          authStateRef.current.token &&
          entry.isIntersecting &&
          !announcementState.read
        ) {
          ReadAnnouncement(announcementState)
            .then((res) => {})
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

  console.log(announcementState.title, announcementState.read);

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
          <Modal
            open={deleteModalOpen}
            onClose={() => setDeleteModalOpen(false)}
          >
            <ModalDialog
              color="danger"
              layout="center"
              size="lg"
              variant="soft"
              role="alertdialog"
            >
              <DialogTitle sx={{ display: "flex", alignItems: "center" }}>
                <Warning />
                Deleting announcement
              </DialogTitle>
              <ModalClose />
              <Divider />
              <DialogContent>
                <Stack direction="column" spacing={2}>
                  <div>
                    You are about to delete the <b>{announcementState.title}</b>{" "}
                    announcement. This action cannot be undone. Are you sure?
                  </div>
                  <Stack direction="row" spacing={1}>
                    <ButtonGroup color="danger" variant="outlined">
                      <Button
                        disabled={deletingAnnouncement}
                        onClick={() => {
                          setDeletingAnnouncement(true);
                          DeleteAnnouncement(announcementState.id)
                            .then((res) => {
                              setDeletingAnnouncement(false);
                              setDeleteModalOpen(false);
                            })
                            .catch((err) => {
                              setLoadingState({
                                error: getError(err),
                                isLoading: loadingState.isLoading,
                              });
                            });
                        }}
                      >
                        Yes
                      </Button>
                      <Button
                        disabled={deletingAnnouncement}
                        onClick={() => setDeleteModalOpen(false)}
                      >
                        No
                      </Button>
                    </ButtonGroup>
                  </Stack>
                </Stack>
              </DialogContent>
            </ModalDialog>
          </Modal>
          <Typography level="h2" sx={{ borderBottom: "1px solid #CDD7E1" }}>
            <Stack
              direction="row"
              alignItems="center"
              justifyContent="space-between"
            >
              <div>{announcementState.title}</div>
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
                    onClick={() => setDeleteModalOpen(true)}
                  />
                </Stack>
              )}
            </Stack>
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
              />
            </div>
          </Paper>
        </Card>
      )}
    </div>
  );
};

export default Announcement;
