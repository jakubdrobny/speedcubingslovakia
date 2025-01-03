import {
  Button,
  Card,
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
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";
import { Link } from "react-router-dom";

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

    GetWCACompetitions(regionValueRef?.current)
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
      {loadingState.error && renderResponseError(loadingState.error)}
      {loadingState.isLoading ? (
        <LoadingComponent title="Loading upcoming WCA competitions..." />
      ) : (
        <Stack spacing={1}>
          {competitions.map((comp: WCACompetitionType, idx1: number) => (
            <Stack component={Card} key={idx1} spacing={1} direction="column">
              <Typography level="h4">{comp.name}</Typography>
              <Typography>
                <b>Place:</b>&nbsp;{comp.venueAddress}
              </Typography>
              <Typography>
                <b>Date:</b>&nbsp;
                {new Date(comp.startdate).toLocaleDateString() +
                  " - " +
                  new Date(comp.enddate).toLocaleDateString()}
              </Typography>
              <Typography>
                <b>Competitors:</b>&nbsp;
                {comp.registered + "/" + comp.competitorLimit}
              </Typography>
              <Typography sx={{ display: "flex", alignItems: "center" }}>
                <b>Events:</b>&nbsp;
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
              </Typography>
              <Typography>
                <b>Delegates:</b>&nbsp;{comp.delegates.join(", ")}
              </Typography>
              <Typography>
                <b>Organizers:</b>&nbsp;{comp.organizers.join(", ")}
              </Typography>
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
