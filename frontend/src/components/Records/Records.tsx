import {
  Button,
  Card,
  CircularProgress,
  Divider,
  Stack,
  Typography,
} from "@mui/joy";
import {
  CompetitionEvent,
  LoadingState,
  RecordsItem,
  RegionSelectGroup,
} from "../../Types";
import {
  getAvailableEvents,
  getError,
  getRecords,
  getRegionGroups,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";

import LanguageIcon from "@mui/icons-material/Language";
import RecordsTable from "./RecordsTable";
import { useEffect } from "react";
import useState from "react-usestateref";

const defaultRegionGroup = "World+World";
const ALL_EVENT = -1;

const Records = () => {
  const [events, setEvents, eventsRef] = useState<CompetitionEvent[]>([]);
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [currentEventIdx, setCurrentEventIdx, currentEventIdxRef] =
    useState<number>(-1);
  const [regionGroups, setRegionGroups] = useState<RegionSelectGroup[]>([]);
  const [regionValue, setRegionValue, regionValueRef] =
    useState<string>(defaultRegionGroup);
  const [records, setRecords] = useState<RecordsItem[]>([]);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });
    getAvailableEvents()
      .then((res: CompetitionEvent[]) => {
        setEvents(res);
        return getRegionGroups();
      })
      .then((res: RegionSelectGroup[]) => {
        setRegionGroups(res);
        fetchRecords();
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, []);

  const fetchRecords = () => {
    setLoadingState({ isLoading: true, error: {} });
    getRecords(
      currentEventIdxRef.current === -1
        ? -1
        : eventsRef.current[currentEventIdxRef.current].id,
      regionValueRef.current.split("+")[0],
      regionValueRef.current.split("+")[1]
    )
      .then((res: RecordsItem[]) => {
        setRecords(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  return (
    <Stack sx={{ margin: "1em" }} spacing={2}>
      <Typography level="h2" className="bottom-divider">
        Records
      </Typography>
      <Stack direction="row" spacing={1} alignItems="center">
        <Typography level="h3">Event:</Typography>
        <div>
          <LanguageIcon
            className={`cubing-icon profile-cubing-icon-mock`}
            onClick={() => {
              if (!loadingState.isLoading) {
                setCurrentEventIdx(ALL_EVENT);
                fetchRecords();
              }
            }}
            style={{
              color: currentEventIdx === ALL_EVENT ? "#0B6BCB" : "black",
              cursor: "pointer",
              transform: "scale(1.25)",
              padding: "0 9px 0 10px",
            }}
          />
          {events.map((event: CompetitionEvent, idx: number) => (
            <span
              key={idx}
              className={`cubing-icon event-${event.iconcode} profile-cubing-icon-mock`}
              onClick={() => {
                if (!loadingState.isLoading) {
                  setCurrentEventIdx(idx);
                  fetchRecords();
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
        <Typography level="h3">Region:</Typography>
        <select
          value={regionValue}
          onChange={(e) => {
            setRegionValue(e.target.value);
            fetchRecords();
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
      ) : loadingState.isLoading ? (
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <CircularProgress />
          &nbsp; &nbsp; <Typography level="h3">Fetching results...</Typography>
        </div>
      ) : (
        <RecordsTable recordItems={records} />
      )}
    </Stack>
  );
};

export default Records;
