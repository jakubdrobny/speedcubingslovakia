import {
  AnnouncementState,
  CompetitionEvent,
  LoadingState,
  Tag,
  initialAnnouncementState,
} from "../../Types";
import {
  Box,
  Button,
  Card,
  Chip,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  Option,
  Select,
  Stack,
  Typography,
} from "@mui/joy";
import {
  getAnnouncementById,
  getAvailableTags,
  getError,
  renderResponseError,
  updateAnnoncement,
} from "../../utils";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { CompetitionEditProps } from "../../Types";
import MDEditor from "@uiw/react-md-editor";

const AnnouncementEdit: React.FC<CompetitionEditProps> = ({ edit }) => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [availableTags, setAvailableTags] = useState<Tag[]>([]);
  const [announcementState, setAnnoucementState] = useState<AnnouncementState>(
    initialAnnouncementState
  );
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    if (edit) {
      getAnnouncementById(id)
        .then((newState: AnnouncementState | undefined) => {
          if (newState === undefined) {
            navigate("/not-found");
          } else {
            setAnnoucementState({ ...announcementState, ...newState });
          }
        })
        .catch((err) => {
          setLoadingState((ps) => ({ ...ps, error: getError(err) }));
        });
    }

    getAvailableTags()
      .then((res: Tag[]) => {
        setAvailableTags(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, []);

  const handleSelectedTagsChange = (selectedTagsLabels: string[]) => {
    const selectedTags = selectedTagsLabels.map(
      (tagLabel) =>
        availableTags.find((tag) => tag.tagLabel === tagLabel) as Tag
    );
    setAnnoucementState({ ...announcementState, tags: selectedTags });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !announcementState.title ||
      !announcementState.content ||
      !announcementState.tags
    )
      return;

    setLoadingState({ isLoading: true, error: {} });
    updateAnnoncement(announcementState, edit)
      .then((res: AnnouncementState) => {
        setAnnoucementState(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  return (
    <Stack style={{ margin: "1em" }} spacing={2}>
      {loadingState.error && renderResponseError(loadingState.error)}
      <Card>
        <Typography level="h3" sx={{ borderBottom: "1px solid #CDD7E1" }}>
          {edit ? `Edit ${announcementState.title}` : "Create"} announcement
        </Typography>
        <Stack spacing={2}>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Title:</Typography>
            </FormLabel>
            <Input
              placeholder="Enter announcement title..."
              value={announcementState.title}
              disabled={loadingState.isLoading}
              onChange={(e) =>
                setAnnoucementState({
                  ...announcementState,
                  title: e.target.value,
                })
              }
            />
            <FormHelperText>This field is required.</FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Tags:</Typography>
            </FormLabel>
            <Select
              multiple
              value={announcementState.tags.map((tag) => tag.tagLabel)}
              onChange={(e, val) => handleSelectedTagsChange(val)}
              disabled={loadingState.isLoading}
              renderValue={(selected) => (
                <Box sx={{ display: "flex", gap: "0.25rem" }}>
                  {selected.map((selectedOption, idx) => (
                    <Chip key={idx} variant="soft" color="primary">
                      {selectedOption.value}
                    </Chip>
                  ))}
                </Box>
              )}
            >
              {availableTags.map((tag: Tag, idx: number) => (
                <Option key={idx} value={tag.tagLabel} label={tag.tagLabel}>
                  {tag.tagLabel}
                </Option>
              ))}
            </Select>
            <FormHelperText>Choose at least 1 tag.</FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Content:</Typography>
            </FormLabel>
            <div data-color-mode="light">
              <div className="wmde-markdown-var"> </div>
              <MDEditor
                value={announcementState.content}
                onChange={(newContent) =>
                  setAnnoucementState({
                    ...announcementState,
                    content: newContent || "",
                  })
                }
              />
            </div>
            <FormHelperText>This field is required.</FormHelperText>
          </FormControl>
          <FormControl>
            <Button onClick={handleSubmit} loading={loadingState.isLoading}>
              {edit ? "Edit" : "Create"} announcement
            </Button>
          </FormControl>
        </Stack>
      </Card>
    </Stack>
  );
};

export default AnnouncementEdit;
