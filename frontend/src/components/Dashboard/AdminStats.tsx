import { Alert, Card, Typography } from "@mui/joy";
import { Stack } from "@mui/system";
import { useEffect } from "react";
import { Chart } from "react-google-charts";
import { LoadingState } from "../../Types";
import {
  getAdminStats,
  getError,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import LoadingComponent from "../Loading/LoadingComponent";
import useState from "react-usestateref";
import { StackedBarChart } from "@mui/icons-material";

const AdminStats = () => {
  const [chartData, setChartData, chartDataRef] = useState<any[]>([]);
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [overallStats, setOverallStats] = useState<{
    total?: number;
    max?: number;
    median?: number;
    average?: number;
  }>({});

  useEffect(() => {
    setLoadingState((p) => ({ ...p, isLoading: true }));
    getAdminStats()
      .then((res) => {
        setChartData([
          res.chartData.columnNames,
          ...res.chartData.data.map((e) =>
            e.map((e2, i) => (i === 0 ? e2 : parseInt(e2))),
          ),
        ]);
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
    <Stack sx={{ margin: "1em" }} spacing={1} direction="column">
      <Typography level="h2">Admin Stats</Typography>
      {loadingState.isLoading ? (
        <LoadingComponent title="Fetching data..." />
      ) : !isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <Stack spacing={1}>
          <Stack
            spacing={1}
            direction={{ xs: "column", lg: "row" }}
            flexWrap="wrap"
          >
            <Card
              sx={{
                display: "flex",
                alignItems: "center",
                flexDirection: "row",
              }}
            >
              <Typography fontWeight="bold">Total competitors:</Typography>
              <Typography>{overallStats.total}</Typography>
            </Card>
            <Card
              sx={{
                display: "flex",
                alignItems: "center",
                flexDirection: "row",
              }}
            >
              <Typography fontWeight="bold">Max. competitors:</Typography>
              <Typography>{overallStats.max}</Typography>
            </Card>
            <Card
              sx={{
                display: "flex",
                alignItems: "center",
                flexDirection: "row",
              }}
            >
              <Typography fontWeight="bold">Median competitors:</Typography>
              <Typography>{overallStats.median}</Typography>
            </Card>
            <Card
              sx={{
                display: "flex",
                alignItems: "center",
                flexDirection: "row",
              }}
            >
              <Typography fontWeight="bold">Average competitors:</Typography>
              <Typography>{overallStats.average}</Typography>
            </Card>
          </Stack>
          <Card sx={{ backgroundColor: "white" }}>
            <Chart
              chartType="LineChart"
              data={chartDataRef.current}
              options={{
                title: "Competitors over time",
                curveType: "function",
                legend: { position: "bottom" },
              }}
              legendToggle
            />
          </Card>
        </Stack>
      )}{" "}
    </Stack>
  );
};

export default AdminStats;
