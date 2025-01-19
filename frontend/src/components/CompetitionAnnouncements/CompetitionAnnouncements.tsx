import {
  Button,
  extendTheme,
  CssVarsProvider,
  Stack,
  Typography,
} from "@mui/joy";
import { useContext, useEffect } from "react";
import useState from "react-usestateref";
import {
  AuthContextType,
  CompetitionAnnouncementSubcriptionUpdateResponse,
  CompetitionAnnouncementSubscription,
  LoadingState,
  RegionSelectGroup,
} from "../../Types";
import {
  GetAnnouncementSubscriptions,
  getError,
  GetStateFromRegionPrecise,
  GetWCARegionGroups,
  isObjectEmpty,
  renderResponseError,
  saveCurrentLocation,
  UpdateCompetitionAnnouncementSubscription,
} from "../../utils/utils";
import { Link } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";
import { AxiosError } from "axios";
import { InfoTooltipTitle } from "./InfoTooltip";
import RegionGroupSelect from "../RegionGroupSelect";
import SubscriptionMap from "./SubscriptionMap";

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

const CompetitionAnnouncements = () => {
  const [loadingState, setLoadingState] = useState<
    LoadingState & { isLoadingSubs: boolean }
  >({
    isLoading: false,
    error: {},
    isLoadingSubs: false,
  });
  const [regionGroups, setRegionGroups] = useState<RegionSelectGroup[]>([]);
  const [regionValue, setRegionValue] = useState<string>(defaultRegionGroup);
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
    const fetchAnnouncementSubscriptions = () => {
      setLoadingState((p) => ({ ...p, isLoadingSubs: true }));
      GetAnnouncementSubscriptions()
        .then((res: CompetitionAnnouncementSubscription[]) => {
          const newSubscriptions = new Map<
            string,
            CompetitionAnnouncementSubscription
          >();
          for (const entry of res) {
            const location =
              entry.countryName + (entry.state ? ", " + entry.state : "");
            newSubscriptions.set(location, entry);
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

    GetWCARegionGroups()
      .then((res: RegionSelectGroup[]) => {
        setRegionGroups(res);
        fetchAnnouncementSubscriptions();
      })
      .catch((err) => {
        setLoadingState((p) => ({
          ...p,
          isLoading: false,
          error: getError(err),
        }));
      });
  }, []);

  const handleSubscribeChange = () => {
    setLoadingState((p) => ({ ...p, isLoadingSubs: true }));
    UpdateCompetitionAnnouncementSubscription(
      regionPrecise,
      !subscriptions.get(regionPrecise)?.subscribed,
    )
      .then((res: CompetitionAnnouncementSubcriptionUpdateResponse) => {
        const newSub = subscriptions.get(regionPrecise) || {
          countryId: regionPrecise,
          countryName: regionPrecise,
          subscribed: false,
          state: GetStateFromRegionPrecise(regionPrecise),
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
        Competition Announcements Newsletter
      </Typography>
      <InfoTooltipTitle notInsideTooltip={true} />
      <CssVarsProvider theme={subscriptionTheme}>
        <Stack
          direction={{ xs: "column", md: "row" }}
          spacing={2}
          sx={{ px: 2 }}
        >
          <Stack spacing={2} direction="row">
            <Typography level="h3">Region:</Typography>
            <RegionGroupSelect
              regionValue={regionValue}
              handleRegionChange={(newRegionValue: string) =>
                setRegionValue(newRegionValue)
              }
              regionGroups={regionGroups}
              disabled={loadingState.isLoading}
            />
          </Stack>
          {!loggedIn ? (
            <Button
              variant="soft"
              component={Link}
              color="warning"
              sx={{ px: { xs: 0, md: 2 }, width: "auto" }}
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
      </CssVarsProvider>
      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}
      <SubscriptionMap />
    </Stack>
  );
};

export default CompetitionAnnouncements;
