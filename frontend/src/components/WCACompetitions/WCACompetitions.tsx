import {
  Button,
  Card,
  Chip,
  Divider,
  Option,
  Select,
  Stack,
  Typography,
} from "@mui/joy";
import { Box } from "@mui/system";
import { useEffect } from "react";
import useState from "react-usestateref";
import {
  CompetitionEvent,
  LoadingState,
  RegionSelectGroup,
  WCACompetitionType,
} from "../../Types";
import {
  getCubingIconClassName,
  getError,
  getRegionGroups,
  GetWCACompetitions,
  renderResponseError,
  renderUpcomingWCACompetitionDateRange,
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";
import { Link } from "react-router-dom";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

const defaultRegionGroup = "Country+Slovakia";

const WCACompetitions = () => {
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const [competitions, setCompetitions] = useState<WCACompetitionType[]>([]);
  const [regionGroups, setRegionGroups] = useState<RegionSelectGroup[]>([]);
  const [regionValue, setRegionValue, regionValueRef] =
    useState<string>(defaultRegionGroup);

  useEffect(() => {
    getRegionGroups()
      .then((res: RegionSelectGroup[]) => {
        res = res.filter((g: RegionSelectGroup) => g.groupName === "Country");
        setRegionGroups(res);
        fetchWCACompetitions();
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, []);

  const fetchWCACompetitions = () => {
    setLoadingState({ isLoading: true, error: {} });

    const _regionValueSplit = regionValueRef.current.split("+");
    const regionPrecise = _regionValueSplit[_regionValueSplit.length - 1];
    GetWCACompetitions(regionPrecise)
      .then((res: WCACompetitionType[]) => {
        setCompetitions(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  return (
    <Stack spacing={3} sx={{ mt: 3 }}>
      <Typography
        level="h2"
        sx={{ pl: 1, borderBottom: "1px solid #636d7433" }}
      >
        Upcoming WCA Competitions
      </Typography>
      <Stack
        direction="row"
        spacing={1}
        flexWrap="wrap"
        gap="10px"
        sx={{ pl: 2 }}
      >
        <Typography level="h3">Region:</Typography>
        <Select
          value={regionValue}
          onChange={(_, val) => {
            setRegionValue(val || "");
            fetchWCACompetitions();
          }}
          renderValue={(sel) => <Box sx={{ pl: 1 }}>{sel?.label}</Box>}
          sx={{ minWidth: "200px" }}
          disabled={loadingState.isLoading}
        >
          {regionGroups.map((regionGroup: RegionSelectGroup, idx: number) => (
            <div key={idx}>
              <Option value={regionGroup.groupName} disabled sx={{ pl: 2 }}>
                <b style={{ color: "black" }}>{regionGroup.groupName}</b>
              </Option>
              {regionGroup.groupMembers.map(
                (groupMember: string, idx2: number) => (
                  <Option
                    key={idx2}
                    value={regionGroup.groupName + "+" + groupMember}
                    label={groupMember}
                    sx={{ pl: 4 }}
                    color="neutral"
                  >
                    {groupMember}
                  </Option>
                ),
              )}
            </div>
          ))}
        </Select>
      </Stack>
      <Divider />
      {loadingState.error && renderResponseError(loadingState.error)}
      {loadingState.isLoading ? (
        <LoadingComponent title="Loading upcoming WCA competitions..." />
      ) : (
        <Stack spacing={1}>
          {competitions.map((comp: WCACompetitionType, idx1: number) => (
            <Stack component={Card} key={idx1} direction="column">
              <Typography level="h3">{comp.name}</Typography>
              <Divider />
              <Typography>
                <b>Place:</b>&nbsp;{comp.venueAddress}
              </Typography>
              <Typography>
                <b>Date:</b>&nbsp;
                {renderUpcomingWCACompetitionDateRange(
                  comp.startdate,
                  comp.enddate,
                )}
              </Typography>
              {dayjs().isBefore(dayjs(comp.registrationOpen)) ? (
                <Stack spacing={1} direction="row">
                  <Typography>
                    <b>Registration opens:</b>
                  </Typography>
                  <Typography>
                    {new Date(comp.registrationOpen).toLocaleDateString() +
                      " " +
                      new Date(comp.registrationOpen).toLocaleTimeString()}
                  </Typography>
                  <Chip color="warning">
                    {dayjs(comp.registrationOpen).fromNow()}
                  </Chip>
                </Stack>
              ) : (
                <Stack spacing={1} direction="row">
                  <Typography>
                    <b>Competitors:</b>
                  </Typography>
                  <Typography>
                    {comp.registered + "/" + comp.competitorLimit}
                  </Typography>
                  <Chip
                    color={
                      comp.registered === comp.competitorLimit
                        ? "danger"
                        : "success"
                    }
                  >
                    {comp.registered === comp.competitorLimit
                      ? "Full"
                      : (comp.competitorLimit - comp.registered).toString() +
                        " spot" +
                        (comp.competitorLimit - comp.registered > 1
                          ? "s"
                          : "") +
                        " remaining"}
                  </Chip>
                </Stack>
              )}
              <Stack direction="row" alignItems="center" spacing={1}>
                <Typography>
                  <b>Events:</b>
                </Typography>
                <Stack spacing={1} direction="row">
                  {comp.events.map((event: CompetitionEvent, idx2: number) => (
                    <span
                      key={idx2 + 100000}
                      className={`${getCubingIconClassName(
                        event.iconcode,
                      )} profile-cubing-icon-mock`}
                    />
                  ))}
                </Stack>
              </Stack>
              <Divider />
              <Button
                sx={{ float: "right" }}
                variant="outlined"
                component={Link}
                to={comp.url}
              >
                More info!
              </Button>
            </Stack>
          ))}
        </Stack>
      )}
    </Stack>
  );
};

export default WCACompetitions;
