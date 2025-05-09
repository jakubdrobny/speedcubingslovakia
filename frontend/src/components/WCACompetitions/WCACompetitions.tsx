import {
  Button,
  Divider,
  Stack,
  extendTheme,
  CssVarsProvider,
  Typography,
} from "@mui/joy";
import { useContext, useEffect } from "react";
import useState from "react-usestateref";
import {
  AuthContextType,
  CompetitionAnnouncementSubscription,
  LoadingState,
  RegionSelectGroup,
  WCACompetitionType,
} from "../../Types";
import {
  GetAnnouncementSubscriptions,
  getError,
  GetWCACompetitions,
  GetWCARegionGroups,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";
import { Link } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";
import { AxiosError } from "axios";
import WCACompetition from "./WCACompetition";
import { InfoTooltip } from "../CompetitionAnnouncements/InfoTooltip";
import RegionGroupSelect from "../RegionGroupSelect";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

const defaultRegionGroup = "Country+Slovakia";

const subscriptionTheme = extendTheme({
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
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const loggedIn =
    authStateRef.current.token !== undefined &&
    authStateRef.current.token !== "";
  const [subscriptionTooltipOpen, setSubscriptionTooltipOpen] = useState(false);

  useEffect(() => {
    GetWCARegionGroups()
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

  const handleRegionChange = (newRegionValue: string) => {
    setRegionValue(newRegionValue);
    fetchWCACompetitions();
  };

  const fetchAnnouncementSubscriptions = () => {
    setLoadingState((p) => ({ ...p, isLoadingSubs: true }));
    GetAnnouncementSubscriptions()
      .then((res: CompetitionAnnouncementSubscription[]) => {
        const newSubscriptions = new Map<
          string,
          CompetitionAnnouncementSubscription
        >();
        for (const entry of res) {
          newSubscriptions.set(entry.countryName, entry);
        }
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

  return (
    <Stack spacing={3} sx={{ mt: 3 }}>
      <Typography
        level="h2"
        sx={{ pl: 1, borderBottom: "1px solid #636d7433" }}
      >
        Upcoming WCA Competitions
      </Typography>
      <CssVarsProvider theme={subscriptionTheme}>
        <Stack
          direction={{ xs: "column", md: "row" }}
          spacing={2}
          sx={{ pl: 2 }}
        >
          <Stack spacing={2} direction="row">
            <Typography level="h3">Region:</Typography>
            <RegionGroupSelect
              regionGroups={regionGroups}
              handleRegionChange={handleRegionChange}
              regionValue={regionValue}
              disabled={loadingState.isLoading}
            />
          </Stack>
          <Stack
            spacing={1}
            direction="row"
            sx={{
              display: "flex",
              alignItems: "center",
            }}
          >
            <Button
              variant="soft"
              component={Link}
              color="warning"
              sx={{ px: 2 }}
              to={"/competitions/announcements"}
            >
              Subscribe
            </Button>
            <InfoTooltip
              open={subscriptionTooltipOpen}
              setOpen={setSubscriptionTooltipOpen}
            />
          </Stack>
        </Stack>
      </CssVarsProvider>
      <Divider />
      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}
      {loadingState.isLoading ? (
        <LoadingComponent title="Loading upcoming WCA competitions..." />
      ) : (
        <Stack spacing={2}>
          {competitions.map(
            (comp: WCACompetitionType, idx1: number) =>
              dayjs().isBefore(dayjs(comp.enddate).add(1, "day")) && (
                <WCACompetition comp={comp} key={idx1} />
              ),
          )}
        </Stack>
      )}
    </Stack>
  );
};

export default WCACompetitions;
