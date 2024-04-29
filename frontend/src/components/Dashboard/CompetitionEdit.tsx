import {
  Alert,
  Box,
  Button,
  Card,
  Chip,
  CircularProgress,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  Option,
  Select,
  Skeleton,
  Stack,
  Typography,
} from "@mui/joy";
import {
  AuthContextType,
  CompetitionData,
  CompetitionEvent,
  CompetitionState,
} from "../../Types";
import {
  authorizeAdmin,
  formatCompetitionDateForInput,
  getAvailableEvents,
  getError,
  initialCompetitionState,
  updateCompetition,
} from "../../utils";
import { useContext, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { AuthContext } from "../../context/AuthContext";
import { CompetitionEditProps } from "../../Types";
import { getCompetitionById } from "../../utils";

const CompetitionEdit: React.FC<CompetitionEditProps> = ({ edit }) => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>(
    []
  );
  const [competitionState, setCompetitionState] = useState<CompetitionState>(
    initialCompetitionState
  );
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    setIsLoading(true);
    setError("");

    if (edit) {
      getCompetitionById(id)
        .then((info: CompetitionData | undefined) => {
          if (info === undefined) {
            navigate("/not-found");
          } else {
            setCompetitionState({ ...competitionState, ...info });
          }
        })
        .catch((err) => {
          setError(getError(err));
        });
    }

    getAvailableEvents()
      .then((res: CompetitionEvent[]) => {
        setAvailableEvents(res);
        setIsLoading(false);
      })
      .catch((err) => {
        setIsLoading(false);
        setError(getError(err));
      });
  }, []);

  const handleSelectedEventsChange = (selectedEventsNames: string[]) => {
    const selectedEvents = selectedEventsNames.map(
      (eName) =>
        availableEvents.find((e) => e.displayname === eName) as CompetitionEvent
    );
    setCompetitionState({ ...competitionState, events: selectedEvents });
  };

  const handleSubmit = () => {
    if (
      !competitionState.name ||
      !competitionState.startdate ||
      !competitionState.enddate ||
      !competitionState.events
    )
      return;

    setIsLoading(true);
    updateCompetition(competitionState, edit)
      .then((res: CompetitionState) => {
        setCompetitionState(res);
        setIsLoading(false);
        setError("");
      })
      .catch((err) => {
        setIsLoading(false);
        setError(getError(err));
      });
  };

  return (
    <Stack style={{ marginTop: "2em" }} spacing={2}>
      {error && <Alert color="danger">{error}</Alert>}
      <Card>
        <Typography level="h3" sx={{ borderBottom: "1px solid #CDD7E1" }}>
          {edit ? `Edit ${competitionState.name}` : "Create"} competition
        </Typography>
        <Stack spacing={2}>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Name:</Typography>
            </FormLabel>
            <Input
              placeholder="Enter competition name..."
              value={competitionState.name}
              disabled={isLoading}
              onChange={(e) =>
                setCompetitionState({
                  ...competitionState,
                  name: e.target.value,
                })
              }
            />
          </FormControl>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Starting date:</Typography>
            </FormLabel>
            <Input
              type="datetime-local"
              value={formatCompetitionDateForInput(competitionState.startdate)}
              disabled={isLoading}
              onChange={(e) =>
                setCompetitionState({
                  ...competitionState,
                  startdate: !isNaN(Date.parse(e.target.value))
                    ? e.target.value
                    : competitionState.startdate,
                })
              }
            />
            <FormHelperText>(format is mm/dd/yyyy, hh/mm)</FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Ending date:</Typography>
            </FormLabel>
            <Input
              type="datetime-local"
              value={formatCompetitionDateForInput(competitionState.enddate)}
              disabled={isLoading}
              onChange={(e) =>
                setCompetitionState({
                  ...competitionState,
                  enddate: !isNaN(Date.parse(e.target.value))
                    ? e.target.value
                    : competitionState.enddate,
                })
              }
            />
            <FormHelperText>(format is mm/dd/yyyy, hh/mm)</FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>
              <Typography level="h4">Events:</Typography>
            </FormLabel>
            <Select
              multiple
              value={competitionState.events.map((e) => e.displayname)}
              onChange={(e, val) => handleSelectedEventsChange(val)}
              disabled={isLoading}
              renderValue={(selected) => (
                <Box sx={{ display: "flex", gap: "0.25rem" }}>
                  {selected.map((selectedOption, idx) => (
                    <Chip key={idx} variant="soft" color="primary">
                      <span
                        className={`cubing-icon event-${selectedOption.label}`}
                      />
                      &nbsp;{selectedOption.value}
                    </Chip>
                  ))}
                </Box>
              )}
            >
              {availableEvents.map((ev: CompetitionEvent) => (
                <Option key={ev.id} value={ev.displayname} label={ev.iconcode}>
                  <span className={`cubing-icon event-${ev.iconcode}`} />
                  {ev.displayname}
                </Option>
              ))}
            </Select>
          </FormControl>
          <FormControl>
            <Button onClick={handleSubmit} loading={isLoading}>
              {edit ? "Edit" : "Create"} competition
            </Button>
          </FormControl>
        </Stack>
      </Card>
    </Stack>
  );
};

export default CompetitionEdit;
