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
import { useContext, useEffect } from "react";
import useState from "react-usestateref";
import {
  AuthContextType,
  CompetitionAnnouncementSubcriptionUpdateResponse,
  CompetitionAnnouncementSubscription,
  CompetitionEvent,
  LoadingState,
  RegionSelectGroup,
  WCACompetitionType,
} from "../../Types";
import {
  getAnnouncementSubscriptions,
  getCubingIconClassName,
  getError,
  getRegionGroups,
  GetWCACompetitions,
  isObjectEmpty,
  renderResponseError,
  renderUpcomingWCACompetitionDateRange,
  saveCurrentLocation,
  updateCompetitionAnnouncementSubscription,
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";
import { Link } from "react-router-dom";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { AuthContext } from "../../context/AuthContext";
import { AxiosError } from "axios";
dayjs.extend(relativeTime);

const defaultRegionGroup = "Country+Slovakia";

const WCACompetitions = () => {
  const [loadingState, setLoadingState] = useState<
    LoadingState & { isLoadingSubs: boolean }
  >({
    isLoading: false,
    error: {},
    isLoadingSubs: false,
  });
  const [competitions, setCompetitions] = useState<WCACompetitionType[]>([]);
  const [regionGroups, setRegionGroups] = useState<RegionSelectGroup[]>([]);
  const [regionValue, setRegionValue, regionValueRef] =
    useState<string>(defaultRegionGroup);
  const [subscriptions, setSubscriptions] = useState<
    Map<string, CompetitionAnnouncementSubscription>
  >(new Map());
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const loggedIn =
    authStateRef.current.token !== undefined &&
    authStateRef.current.token !== "";
  const regionPrecise = regionValue.split("+")[1];
  const currentlySubscribed = subscriptions.get(regionPrecise)?.subscribed;

  useEffect(() => {
    getRegionGroups()
      .then((res: RegionSelectGroup[]) => {
        res = res.filter((g: RegionSelectGroup) => g.groupName === "Country");
        setRegionGroups(res);
        fetchWCACompetitions();
      })
      .catch((err) => {
        setLoadingState((p) => ({
          ...p,
          isLoading: false,
          error: getError(err),
        }));
      });
  }, []);

  const fetchWCACompetitions = () => {
    setLoadingState((p) => ({
      ...p,
      isLoading: true,
      error: {},
      isLoadingSubs: false,
    }));

    const _regionValueSplit = regionValueRef.current.split("+");
    const regionPrecise = _regionValueSplit[_regionValueSplit.length - 1];
    GetWCACompetitions(regionPrecise)
      .then((res: WCACompetitionType[]) => {
        setCompetitions(res);
        setLoadingState((p) => ({
          ...p,
          isLoading: false,
          error: {},
          isLoadingSubs: true,
        }));
        fetchAnnouncementSubscriptions();
      })
      .catch((err) => {
        setLoadingState((p) => ({
          ...p,
          isLoadingSubs: false,
          isLoading: false,
          error: getError(err),
        }));
      });
  };

  const fetchAnnouncementSubscriptions = () => {
    setLoadingState((p) => ({ ...p, isLoadingSubs: true }));
    getAnnouncementSubscriptions()
      .then((res: CompetitionAnnouncementSubscription[]) => {
        const newSubscriptions = new Map<
          string,
          CompetitionAnnouncementSubscription
        >();
        for (const entry of res) {
          newSubscriptions.set(entry.countryName, entry);
        }
        setSubscriptions(new Map(newSubscriptions));
        setLoadingState((p) => ({ ...p, isLoadingSubs: false }));
      })
      .catch((err: AxiosError) => {
        setLoadingState((p) => ({
          ...p,
          isLoadingSubs: false,
          error: err.status === 401 && !loggedIn ? {} : getError(err),
        }));
      });
  };

  const handleSubscribeChange = () => {
    setLoadingState((p) => ({ ...p, isLoadingSubs: true }));
    updateCompetitionAnnouncementSubscription(
      regionPrecise,
      !subscriptions.get(regionPrecise)?.subscribed,
    )
      .then((res: CompetitionAnnouncementSubcriptionUpdateResponse) => {
        const newSub = subscriptions.get(regionPrecise) || {
          countryId: regionPrecise,
          countryName: regionPrecise,
          subscribed: false,
        };
        newSub.subscribed = res.subscribed;
        setSubscriptions(new Map(subscriptions).set(regionPrecise, newSub));
        setLoadingState((p) => ({ ...p, isLoadingSubs: false }));
      })
      .catch((err) => {
        setLoadingState((p) => ({
          ...p,
          isLoadingSubs: false,
          error: getError(err),
        }));
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
        {!loggedIn ? (
          <Button
            variant="soft"
            component={Link}
            color="warning"
            sx={{ px: 2 }}
            to={import.meta.env.VITE_WCA_GET_CODE_URL || ""}
            onClick={() => saveCurrentLocation(window.location.pathname)}
          >
            Login to subscribe
          </Button>
        ) : (
          subscriptions &&
          subscriptions.size > 0 && (
            <Button
              onClick={handleSubscribeChange}
              variant="soft"
              color={currentlySubscribed ? "success" : "danger"}
              disabled={loadingState.isLoadingSubs}
            >
              {currentlySubscribed ? "Subscribed!" : "Not subscribed"}
            </Button>
          )
        )}
      </Stack>
      <Divider />
      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}
      {loadingState.isLoading ? (
        <LoadingComponent title="Loading upcoming WCA competitions..." />
      ) : (
        <Stack spacing={2}>
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
