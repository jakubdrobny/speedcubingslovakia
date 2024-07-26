import {
  CompetitionEvent,
  LoadingState,
  RecordsItem,
  RegionSelectGroup,
} from "../../Types";
import { Stack, Typography } from "@mui/joy";
import {
  getAvailableEvents,
  getError,
  getRecords,
  getRegionGroups,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";

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
    useState<number>(0);
  const [regionGroups, setRegionGroups] = useState<RegionSelectGroup[]>([]);
  const [regionValue, setRegionValue, regionValueRef] =
    useState<string>(defaultRegionGroup);
  const [records, setRecords] = useState<RecordsItem[]>([]);
  const isfmc = events[currentEventIdx]?.iconcode === "333fm";
  const ismbld = events[currentEventIdx]?.iconcode === "333mbf";

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
      eventsRef.current[currentEventIdxRef.current].id,
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
      <Typography level="h2">Records</Typography>
      <Stack direction="row" spacing={1}>
        <Typography level="h3">Event:</Typography>
        <div>
          <span
            className={`cubing-icon profile-cubing-icon-mock`}
            onClick={() => {
              if (!loadingState.isLoading) {
                setCurrentEventIdx(ALL_EVENT);
                fetchRecords();
              }
            }}
            style={{
              padding: "0 0.25em",
              fontSize: "1.75em",
              color: currentEventIdx === ALL_EVENT ? "#0B6BCB" : "",
              cursor: "pointer",
            }}
          ></span>
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
      <Stack direction="row" spacing={2} flexWrap="nowrap">
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
      {!isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <RecordsTable
          recordItems={records}
          loading={loadingState.isLoading}
          isfmc={isfmc}
          ismbld={ismbld}
        />
      )}
    </Stack>
  );
};

export default Records;
