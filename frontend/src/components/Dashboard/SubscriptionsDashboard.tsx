import { Card, CardContent, Grid, Stack, Table, Typography } from "@mui/joy";
import {
  getError,
  getSubscriptionStats,
  getUserSubscriptionDetails,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import {
  LoadingState,
  SubscriptionStats as StatsType,
  UserSubscriptionDetail,
} from "../../Types";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import LoadingComponent from "../Loading/LoadingComponent";

const StatCard = ({ title, value }: { title: string; value: number }) => (
  <Card variant="outlined">
    <CardContent>
      <Typography level="title-md">{title}</Typography>
      <Typography level="h1" sx={{ mt: 1 }}>
        {value.toLocaleString()}
      </Typography>
    </CardContent>
  </Card>
);

const DetailsTable = ({ details }: { details: UserSubscriptionDetail[] }) => {
  const columnNames = () => ["Name", "Country", "Countries", "Positions"];
  return (
    <Card sx={{ padding: 0, margin: 0, overflowX: "auto" }}>
      <Table
        size="md"
        sx={{
          tableLayout: "auto",
          width: "100%",
          "& td": { whiteSpace: "nowrap" },
        }}
      >
        <thead>
          <tr>
            {columnNames().map((val, idx) => (
              <th key={idx} style={{ height: "1em" }}>
                <b>{val}</b>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {details.map((row) => (
            <tr key={row.id}>
              <td style={{ height: "1em" }}>
                <Link
                  to={`/profile/${row.wca_id ? row.wca_id : row.name}`}
                  style={{
                    color: "#0B6BCB",
                    textDecoration: "none",
                    fontWeight: 555,
                  }}
                >
                  {row.name}
                  {row.wca_id && ` (${row.wca_id})`}
                </Link>
              </td>
              <td style={{ height: "1em" }}>
                <span className={`fi fi-${row.country_iso2.toLowerCase()}`} />
                &nbsp;&nbsp;{row.country_name}
              </td>
              <td style={{ height: "1em" }}>
                <b>{row.country_sub_count}</b> countries
              </td>
              <td style={{ height: "1em" }}>
                <b>{row.position_sub_count}</b> positions
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Card>
  );
};

const SubscriptionsDashboard = () => {
  const [stats, setStats] = useState<StatsType | null>(null);
  const [details, setDetails] = useState<UserSubscriptionDetail[]>([]);
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: true,
    error: {},
  });

  useEffect(() => {
    Promise.all([getSubscriptionStats(), getUserSubscriptionDetails()])
      .then(([statsData, detailsData]) => {
        console.log(statsData, detailsData);
        setStats(statsData);
        setDetails(detailsData);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) =>
        setLoadingState({ isLoading: false, error: getError(err) }),
      );
  }, []);

  return (
    <Stack spacing={2} sx={{ margin: "1em" }}>
      <Typography level="h2" className="bottom-divider">
        Subscriptions Dashboard
      </Typography>

      {loadingState.isLoading && (
        <LoadingComponent title="Fetching subscription data..." />
      )}
      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}

      {!loadingState.isLoading && stats && (
        <Grid container spacing={2} sx={{ flexGrow: 1 }}>
          <Grid xs={12} md={4}>
            <StatCard
              title="Users Subscribed by Country"
              value={stats.country_subscriptions}
            />
          </Grid>
          <Grid xs={12} md={4}>
            <StatCard
              title="Users Subscribed by Position"
              value={stats.position_subscriptions}
            />
          </Grid>
          <Grid xs={12} md={4}>
            <StatCard
              title="Total Unique Subscribers"
              value={stats.total_unique_users}
            />
          </Grid>
        </Grid>
      )}

      {!loadingState.isLoading && details.length > 0 && (
        <DetailsTable details={details} />
      )}
    </Stack>
  );
};

export default SubscriptionsDashboard;
