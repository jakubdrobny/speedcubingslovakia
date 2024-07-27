import {
  Alert,
  Button,
  ButtonGroup,
  Divider,
  Option,
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
import {
  getAvailableEvents,
  getError,
  getRankings,
  getRegionGroups,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";

import RankingsTable from "./RankingsTable";
import { useEffect } from "react";
import useState from "react-usestateref";

const defaultRegionGroup = "World+World";
const defaultQueryTypeValue = "100+Persons";

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
  const isfmc = events[currentEventIdx]?.iconcode === "333fm";
  const ismbld = events[currentEventIdx]?.iconcode === "333mbf";
  const [queryTypeValue, setQueryTypeValue, queryTypeValueRef] =
    useState<string>(defaultQueryTypeValue);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });
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
    if (!singleRef.current && ismbld) return;

    setLoadingState({ isLoading: true, error: {} });
    getRankings(
      eventsRef.current[currentEventIdxRef.current].id,
      singleRef.current,
      regionValueRef.current.split("+")[0],
      regionValueRef.current.split("+")[1],
      queryTypeValueRef.current
    )
      .then((res: RankingsEntry[]) => {
        setRankings(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  return (
    <Stack sx={{ margin: "1em" }} spacing={2}>
      <Typography level="h2" className="bottom-divider">
        Rankings
      </Typography>
      <Stack direction="row" spacing={1}>
        <Typography level="h3">Event:</Typography>
        <div>
          {events.map((event: CompetitionEvent, idx: number) => (
            <span
              key={idx}
              className={`cubing-icon event-${event.iconcode} profile-cubing-icon-mock`}
              onClick={() => {
                if (!loadingState.isLoading) {
                  if (
                    eventsRef &&
                    eventsRef.current &&
                    idx < eventsRef.current.length &&
                    eventsRef.current[idx].displayname === "MBLD"
                  )
                    setSingle(true);
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
      <Stack direction="row" spacing={2} flexWrap="nowrap">
        <Typography level="h3">Type:</Typography>
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
          {!ismbld && (
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
          )}
        </ButtonGroup>
        <select
          value={queryTypeValue}
          defaultValue={defaultQueryTypeValue}
          onChange={(e) => {
            setQueryTypeValue(e.target.value);
            fetchRankings();
          }}
          color="primary"
        >
          <option value="100+Persons">100 Persons</option>
          <option value="1000+Persons">1000 Persons</option>
          <option value="100+Results">100 Results</option>
          <option value="1000+Results">1000 Results</option>
        </select>
      </Stack>
      <Stack direction="row" spacing={2}>
        <Typography level="h3">Region:</Typography>
        <select
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
      <Divider />
      {!isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <RankingsTable
          rankings={rankings}
          single={single}
          loading={loadingState.isLoading}
          isfmc={isfmc}
          ismbld={ismbld}
        />
      )}
    </Stack>
  );
};

export default Rankings;
