import { Alert, Card, Typography } from "@mui/joy";
import { Stack } from "@mui/system";
import { useEffect, useState } from "react";
import { Chart } from "react-google-charts";
import { LoadingState } from "../../Types";
import {
  getAdminStats,
  getError,
  initialLoadingState,
  renderResponseError,
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";

const AdminStats = () => {
  const [chartData, setChartData] = useState<any[]>([]);
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [overallState, setOverallStats] = useState<{
    total?: number;
    max?: number;
    median?: number;
    average?: number;
  }>({});

  useEffect(() => {
    setLoadingState((p) => ({ ...p, isLoading: true }));
    getAdminStats()
      .then((res) => {
        setChartData([res.chartData.columnNames, ...res.chartData.data]);
        setOverallStats({
          total: res.total,
          max: res.max,
          median: res.median,
          average: res.average,
        });
        setLoadingState((p) => ({ ...p, isLoading: false }));
      })
      .catch((err) => {
        setLoadingState((p) => ({
          ...p,
          isLoading: false,
          error: getError(err),
        }));
      });
  }, []);

  return (
    <Stack sx={{ margin: "1em 0" }} spacing={1} direction="column">
      <Typography level="h2">Admin Stats</Typography>
      {loadingState.isLoading ? (
        <LoadingComponent title="Fetching data..." />
      ) : loadingState.error ? (
        renderResponseError(loadingState.error)
      ) : (
        <Card sx={{ backgroundColor: "white" }}>
          <Chart
            chartType="Line"
            data={chartData}
            options={{
              title: "Competitors over time",
              curveType: "",
              legend: { position: "bottom" },
            }}
            legendToggle
          />
        </Card>
      )}{" "}
    </Stack>
  );
};

export default AdminStats;
