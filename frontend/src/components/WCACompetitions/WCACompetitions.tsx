import {
  Button,
  Divider,
  IconButton,
  Option,
  Select,
  Stack,
  ThemeProvider,
  Tooltip,
  Typography,
} from "@mui/joy";
import { Box, createTheme } from "@mui/system";
import { useContext, useEffect } from "react";
import useState from "react-usestateref";
import {
  AuthContextType,
  CompetitionAnnouncementSubcriptionUpdateResponse,
  CompetitionAnnouncementSubscription,
  LoadingState,
  RegionSelectGroup,
  WCACompetitionType,
} from "../../Types";
import {
  getAnnouncementSubscriptions,
  getError,
  getRegionGroups,
  GetWCACompetitions,
  isObjectEmpty,
  renderResponseError,
  saveCurrentLocation,
  updateCompetitionAnnouncementSubscription,
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";
import { Link } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";
import { AxiosError } from "axios";
import { HelpOutline } from "@mui/icons-material";
import WCACompetition from "./WCACompetition";

const defaultRegionGroup = "Country+Slovakia";

const subscriptionTheme = createTheme({
  breakpoints: {
    values: {
      xs: 0,
      sm: 500,
      md: 750,
      lg: 1200,
      xl: 1536,
    },
  },
});

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
  const [subscriptionTooltipOpen, setSubscriptionTooltipOpen] = useState(false);

  useEffect(() => {
    getRegionGroups()
      .then((res: RegionSelectGroup[]) => {
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
      <ThemeProvider theme={subscriptionTheme}>
        <Stack
          direction={{ xs: "column", md: "row" }}
          spacing={2}
          sx={{ pl: 2 }}
        >
          <Stack spacing={2} direction="row">
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
              {regionGroups.map(
                (regionGroup: RegionSelectGroup, idx: number) => (
                  <div key={idx}>
                    <Option
                      value={regionGroup.groupName}
                      disabled
                      sx={{ pl: 2 }}
                    >
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
                ),
              )}
            </Select>
          </Stack>
          <Stack
            spacing={1}
            direction="row"
            sx={{
              display: "flex",
              alignItems: "center",
            }}
          >
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
            <Tooltip
              variant="soft"
              color="primary"
              title={
                <Box>
                  <Typography fontWeight="bold">
                    Tired of checking the WCA website for new competitions?
                  </Typography>
                  <Typography fontSize="1em">
                    Subscribe to our <b>newsletter</b> to receive emails when{" "}
                    <b>new WCA competitions</b> are <b>announced</b> in any
                    country of your choice.{" "}
                  </Typography>
                  <Typography fontSize="1em">
                    You can choose one or <b>multiple countries</b> and
                    unsubscribe from any of them at any time.
                  </Typography>
                  <Typography fontWeight="bold" fontSize="0.9em">
                    Enjoy :)
                  </Typography>
                </Box>
              }
              open={subscriptionTooltipOpen}
              disableInteractive={false}
              enterTouchDelay={0}
              enterDelay={0}
              leaveDelay={0}
            >
              <IconButton
                onMouseEnter={() => setSubscriptionTooltipOpen(true)}
                onMouseLeave={() => setSubscriptionTooltipOpen(false)}
                onTouchStart={() => setSubscriptionTooltipOpen((p) => !p)}
              >
                <HelpOutline fontSize="small" />
              </IconButton>
            </Tooltip>
          </Stack>
        </Stack>
      </ThemeProvider>
      <Divider />
      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}
      {loadingState.isLoading ? (
        <LoadingComponent title="Loading upcoming WCA competitions..." />
      ) : (
        <Stack spacing={2}>
          {competitions.map((comp: WCACompetitionType, idx1: number) => (
            <WCACompetition comp={comp} key={idx1} />
          ))}
        </Stack>
      )}
    </Stack>
  );
};

export default WCACompetitions;
