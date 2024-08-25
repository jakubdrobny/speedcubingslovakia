import { AnnouncementState, LoadingState } from "../../Types";
import {
  Button,
  ButtonGroup,
  CircularProgress,
  DialogContent,
  DialogTitle,
  Divider,
  Modal,
  ModalClose,
  ModalDialog,
  Stack,
  Typography,
} from "@mui/joy";
import {
  DeleteAnnouncement,
  GetAnnouncements,
  getError,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";

import Announcement from "./Announcement";
import { Warning } from "@mui/icons-material";
import { useEffect } from "react";
import useState from "react-usestateref";

const Announcements = () => {
  const [announcements, setAnnouncements, announcementsRef] = useState<
    AnnouncementState[]
  >([]);
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const [deletingAnnouncement, setDeletingAnnouncement] =
    useState<boolean>(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState<boolean>(false);
  const [deleteCandidate, setDeleteCandidate, deleteCandidateRef] = useState<{
    title: string;
    idx: number;
    id: number;
  }>({ title: "", idx: 0, id: 0 });

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

  const onAnnouncementDelete = (idx: number, title: string, id: number) => {
    setDeleteCandidate({ title, idx, id });
    setDeleteModalOpen(true);
  };

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
                    You are about to delete the <b>{deleteCandidate.title}</b>{" "}
                    announcement. This action cannot be undone. Are you sure?
                  </div>
                  <Stack direction="row" spacing={1}>
                    <ButtonGroup color="danger" variant="outlined">
                      <Button
                        disabled={deletingAnnouncement}
                        onClick={() => {
                          setDeletingAnnouncement(true);

                          console.log(deleteCandidateRef.current.idx);
                          const newAnnouncements = announcements.slice();
                          newAnnouncements.splice(
                            deleteCandidateRef.current.idx,
                            1
                          );
                          setAnnouncements(newAnnouncements);

                          console.log(announcementsRef.current);
                          DeleteAnnouncement(deleteCandidateRef.current.id)
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
          <Typography
            level="h1"
            sx={{ margin: "0.5em 0.5em", borderBottom: "1px solid #CDD7E1" }}
          >
            Announcements
          </Typography>
          <Stack spacing={2} direction="column">
            {announcements.map((announcement, idx) => (
              <Announcement
                key={idx.toString() + announcement.id.toString()}
                idx={idx}
                givenAnnouncementState={announcement}
                onAnnouncementDelete={onAnnouncementDelete}
              />
            ))}
          </Stack>
        </div>
      )}
    </div>
  );
};

export default Announcements;
