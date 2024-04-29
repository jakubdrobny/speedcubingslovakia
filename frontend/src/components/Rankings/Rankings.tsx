import {
  Button,
  ButtonGroup,
  List,
  ListItem,
  Select,
  Stack,
  Typography,
} from "@mui/joy";
import {
  CompetitionEvent,
  LoadingState,
  RankingsEntry,
  RegionSelectGroup,
} from "../../Types";
import Option, { optionClasses } from "@mui/joy/Option";
import {
  getAvailableEvents,
  getError,
  getRankings,
  getRegionGroups,
  initialLoadingState,
} from "../../utils";

import RankingsTable from "./RankingsTable";
import { useEffect } from "react";
import useState from "react-usestateref";

const defaultRegionGroup = "World+World";

const Rankings = () => {
  const [events, setEvents, eventsRef] = useState<CompetitionEvent[]>([]);
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [currentEventIdx, setCurrentEventIdx, currentEventIdxRef] =
    useState<number>(0);
  const [single, setSingle, singleRef] = useState<boolean>(true);
  const [regionGroups, setRegionGroups] = useState<RegionSelectGroup[]>([]);
  const [regionValue, setRegionValue, regionValueRef] =
    useState<string>(defaultRegionGroup);
  const [rankings, setRankings] = useState<RankingsEntry[]>([]);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: "" });
    getAvailableEvents()
      .then((res: CompetitionEvent[]) => {
        setEvents(res);
        return getRegionGroups();
      })
      .then((res: RegionSelectGroup[]) => {
        setRegionGroups(res);
        fetchRankings();
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, []);

  const fetchRankings = () => {
    setLoadingState({ isLoading: true, error: "" });
    getRankings(
      eventsRef.current[currentEventIdxRef.current].id,
      singleRef.current,
      regionValueRef.current
    )
      .then((res: RankingsEntry[]) => {
        setRankings(res);
        setLoadingState({ isLoading: false, error: "" });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  return (
    <Stack sx={{ margin: "1em" }} spacing={2}>
      <Typography level="h2">Rankings</Typography>
      <Stack direction="row" spacing={1}>
        <Typography level="h3">Event:</Typography>
        <div>
          {events.map((event: CompetitionEvent, idx: number) => (
            <span
              key={idx}
              className={`cubing-icon event-${event.iconcode} profile-cubing-icon-mock`}
              onClick={() => {
                if (!loadingState.isLoading) {
                  setCurrentEventIdx(idx);
                  fetchRankings();
                }
              }}
              style={{
                padding: "0 0.25em",
                fontSize: "1.75em",
                color: currentEventIdx === idx ? "#0B6BCB" : "",
                cursor: "pointer",
              }}
            />
          ))}
        </div>
      </Stack>
      <Stack direction="row" spacing={2}>
        <ButtonGroup>
          <Button
            variant={single ? "solid" : "outlined"}
            color="primary"
            disabled={loadingState.isLoading}
            onClick={() => {
              setSingle(true);
              fetchRankings();
            }}
          >
            Single
          </Button>
          <Button
            variant={!single ? "solid" : "outlined"}
            color="primary"
            disabled={loadingState.isLoading}
            onClick={() => {
              setSingle(false);
              fetchRankings();
            }}
          >
            Average
          </Button>
        </ButtonGroup>
        <select
          defaultValue={defaultRegionGroup}
          value={regionValue}
          onChange={(e) => {
            setRegionValue(e.target.value);
            fetchRankings();
          }}
        >
          {regionGroups.map((regionGroup: RegionSelectGroup, idx: number) => (
            <optgroup key={idx} label={regionGroup.groupName}>
              {regionGroup.groupMembers.map((groupMember: string, idx2) => (
                <option
                  key={idx2}
                  value={regionGroup.groupName + "+" + groupMember}
                  label={groupMember}
                >
                  {groupMember}
                </option>
              ))}
            </optgroup>
          ))}
        </select>
      </Stack>
      <RankingsTable rankings={rankings} single={single} />
    </Stack>
  );
};

export default Rankings;
