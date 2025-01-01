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
  CompetitionData,
  CompetitionEvent,
  CompetitionState,
  ResponseError,
} from "../../Types";
import {
  formatCompetitionDateForInput,
  getAvailableEvents,
  getCubingIconClassName,
  getError,
  initialCompetitionState,
  renderResponseError,
  updateCompetition,
} from "../../utils/utils";
import { useNavigate, useParams } from "react-router-dom";

import { CompetitionEditProps } from "../../Types";
import { getCompetitionById } from "../../utils/utils";
import { useEffect } from "react";
import useState from "react-usestateref";

const CompetitionEdit: React.FC<CompetitionEditProps> = ({ edit }) => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [availableEvents, setAvailableEvents] = useState<CompetitionEvent[]>(
    [],
  );
  const [competitionState, setCompetitionState] = useState<CompetitionState>(
    initialCompetitionState,
  );
  const [_, setIsLoading, isLoadingRef] = useState<boolean>(false);
  const [error, setError] = useState<ResponseError>();

  useEffect(() => {
    setIsLoading(true);
    setError({});

    const initPage = async () => {
      try {
        let events: CompetitionEvent[] = await getAvailableEvents();
        events = events.filter((e) => e.displayname !== "Overall");
        setAvailableEvents(events);

        if (!edit) {
          setCompetitionState({ ...competitionState, events });
          setIsLoading(false);
          return;
        }

        const info: CompetitionData | undefined = await getCompetitionById(id);
        if (info === undefined) {
          navigate("/not-found");
        } else {
          setCompetitionState({ ...competitionState, ...info });
        }

        setIsLoading(false);
      } catch (err: any) {
        setIsLoading(false);
        setError(getError(err));
      }
    };

    initPage();
  }, []);

  const handleSelectedEventsChange = (selectedEventsNames: string[]) => {
    const selectedEvents = selectedEventsNames.map(
      (eName) =>
        availableEvents.find(
          (e) => e.displayname === eName,
        ) as CompetitionEvent,
    );
    setCompetitionState({ ...competitionState, events: selectedEvents });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
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
        setError({});
      })
      .catch((err) => {
        setIsLoading(false);
        setError(getError(err));
      });
  };

  return (
    <Stack style={{ marginTop: "2em" }} spacing={2}>
      {error && renderResponseError(error)}
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
              disabled={isLoadingRef.current}
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
              disabled={isLoadingRef.current}
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
              disabled={isLoadingRef.current}
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
              onChange={(_, val) => handleSelectedEventsChange(val)}
              disabled={isLoadingRef.current}
              renderValue={(selected) => (
                <Box sx={{ display: "flex", gap: "0.25rem" }}>
                  {selected.map((selectedOption, idx) => (
                    <Chip key={idx} variant="soft" color="primary">
                      <span
                        className={getCubingIconClassName(selectedOption.label)}
                      />
                      &nbsp;{selectedOption.value}
                    </Chip>
                  ))}
                </Box>
              )}
            >
              {availableEvents.map((ev: CompetitionEvent) => (
                <Option key={ev.id} value={ev.displayname} label={ev.iconcode}>
                  <span className={getCubingIconClassName(ev.iconcode)} />
                  {ev.displayname}
                </Option>
              ))}
            </Select>
          </FormControl>
          <FormControl>
            <Button onClick={handleSubmit} loading={isLoadingRef.current}>
              {edit ? "Edit" : "Create"} competition
            </Button>
          </FormControl>
        </Stack>
      </Card>
    </Stack>
  );
};

export default CompetitionEdit;
