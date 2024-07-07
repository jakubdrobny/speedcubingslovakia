import {
  Alert,
  Box,
  Card,
  CircularProgress,
  Stack,
  Table,
  Tooltip,
  Typography,
} from "@mui/joy";
import { WIN_SMALL, WIN_VERYSMALL } from "../../constants";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { CompetitionContextType } from "../../Types";
import { Help } from "@mui/icons-material";
import { Link } from "react-router-dom";
import { reformatMultiTime } from "../../utils";

const CompetitionResults = () => {
  const {
    competitionState,
    competitionStateRef,
    results,
    loadingState,
    fetchCompetitionResults,
  } = useContext(CompetitionContext) as CompetitionContextType;
  const averageFirst = (() => {
    const format =
      competitionState?.events[competitionState?.currentEventIdx]?.format;
    if (!format || format?.length === 0) return true;
    return format[0] !== "b";
  })();
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);
  const isOverall =
    competitionState?.events[competitionState?.currentEventIdx]?.id === -1;
  const isfmc =
    competitionState?.events[competitionState?.currentEventIdx]?.iconcode ===
    "333fm";
  const ismbld =
    competitionState?.events[competitionState?.currentEventIdx]?.iconcode ===
    "333mbf";

  useEffect(() => {
    const handleResize = () => {
      setWindowWidth(window.innerWidth);
    };

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  const columnNames = () => {
    let columnNames: string[] = [
      "#",
      "Name",
      "Country",
      "Average",
      "Single",
      isfmc ? "Moves" : "Times",
    ];
    if (!averageFirst)
      [columnNames[3], columnNames[4]] = [columnNames[4], columnNames[3]];
    if (isOverall) {
      columnNames.splice(
        3 - Number(windowWidth < WIN_VERYSMALL),
        3 + Number(windowWidth < WIN_VERYSMALL),
        "Score"
      );
    } else {
      if (windowWidth < WIN_SMALL)
        columnNames = [...columnNames.slice(0, 3), ...columnNames.slice(3)];
    }

    if (columnNames.includes("Average") && ismbld)
      columnNames = columnNames.filter((c) => c !== "Average");

    return columnNames;
  };

  useEffect(() => {
    if (competitionState.id !== undefined && competitionState.id !== "")
      fetchCompetitionResults();
  }, []);

  return (
    <div style={{ margin: "1.5em 0.5em", width: "100%" }}>
      {loadingState.results ? (
        <>
          <Typography level="h3" sx={{ display: "flex", alignItems: "center" }}>
            <CircularProgress />
            &nbsp; Loading results ...
          </Typography>
        </>
      ) : loadingState.error ? (
        <Alert color="danger">{loadingState.error}</Alert>
      ) : (
        <Card
          sx={{
            margin: 0,
            padding: 0,
            width: "100%",
          }}
        >
          <Table
            size="md"
            sx={{
              tableLayout: "auto",
              width: "100%",
            }}
          >
            <thead>
              <tr>
                {columnNames().map((val, idx) => {
                  return (
                    <th
                      style={{
                        height: "1em",
                        maxWidth:
                          idx === columnNames().length - 1 ? "400px" : "auto",
                      }}
                      key={idx}
                    >
                      <Stack direction="row" justifyContent="center">
                        <b>{val}</b>
                        {val === "Score" && (
                          <Tooltip
                            placement="right"
                            title={
                              <Box style={{ textAlign: "center" }}>
                                <Typography>
                                  <b>Kinch score</b>
                                </Typography>
                                For more information click{" "}
                                <a href="https://www.speedsolving.com/threads/all-round-rankings-kinchranks.53353/">
                                  here
                                </a>
                                .
                              </Box>
                            }
                            variant="outlined"
                            color="primary"
                            enterTouchDelay={0}
                          >
                            <span style={{ height: "21px" }}>
                              &nbsp;
                              <Help fontSize="small" />
                            </span>
                          </Tooltip>
                        )}
                      </Stack>
                    </th>
                  );
                })}
              </tr>
            </thead>
            <tbody>
              {results.map((result, idx) => {
                result.single = isfmc
                  ? result.single.split(".")[0]
                  : ismbld
                  ? reformatMultiTime(result.single)
                  : result.single;
                result.times = isfmc
                  ? result.times?.map((res) => res.split(".")[0])
                  : ismbld
                  ? result.times?.map((r) => reformatMultiTime(r))
                  : result.times;
                return (
                  <tr key={idx}>
                    <td style={{ height: "1em", width: "center" }}>
                      {result.place}
                    </td>
                    <td style={{ height: "1em, center" }}>
                      <Link
                        to={`/profile/${result.wca_id}`}
                        style={{
                          color: "#0B6BCB",
                          textDecoration: "none",
                          fontWeight: 555,
                        }}
                      >
                        {result.username}
                      </Link>
                    </td>
                    {(windowWidth >= WIN_SMALL ||
                      (isOverall && windowWidth >= WIN_VERYSMALL)) && (
                      <td style={{ height: "1em", textAlign: "center" }}>
                        <span
                          className={`fi fi-${result.country_iso2.toLowerCase()}`}
                        />
                        &nbsp;&nbsp;{result.country_name}
                      </td>
                    )}
                    {isOverall ? (
                      <td style={{ height: "1em", textAlign: "center" }}>
                        <b>{result.score}</b>
                      </td>
                    ) : (
                      <>
                        <td style={{ height: "1em", textAlign: "center" }}>
                          <b>
                            {!averageFirst ? result.single : result.average}
                          </b>
                        </td>
                        {!ismbld && (
                          <td style={{ height: "1em", textAlign: "center" }}>
                            {averageFirst ? result.single : result.average}
                          </td>
                        )}
                        <td style={{ height: "1em", textAlign: "center" }}>
                          {result.times?.join(", ")}
                        </td>
                      </>
                    )}
                  </tr>
                );
              })}
            </tbody>
          </Table>
        </Card>
      )}
    </div>
  );
};

export default CompetitionResults;
