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

const CompetitionResults = () => {
  const { competitionState, results, loadingState, fetchCompetitionResults } =
    useContext(CompetitionContext) as CompetitionContextType;
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
      "",
      "#",
      "Name",
      "Country",
      "Average",
      "Single",
      isfmc ? "Moves" : "Times",
      "",
    ];
    if (!averageFirst)
      [columnNames[4], columnNames[5]] = [columnNames[5], columnNames[4]];
    if (isOverall) {
      columnNames.splice(
        4 - Number(windowWidth < WIN_VERYSMALL),
        3 + Number(windowWidth < WIN_VERYSMALL),
        "Score"
      );
    } else {
      if (windowWidth < WIN_SMALL)
        columnNames = [...columnNames.slice(0, 3), ...columnNames.slice(4)];
    }

    return columnNames;
  };

  useEffect(() => {
    if (competitionState.id !== undefined && competitionState.id !== "")
      fetchCompetitionResults();
  }, []);

  return (
    <div style={{ margin: "1.5em 0.5em" }}>
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
        <Card sx={{ margin: 0, padding: 0 }}>
          <Table size="md">
            <thead>
              <tr>
                {columnNames().map((val, idx) => {
                  let style1: Object = isOverall
                    ? val === ""
                      ? { height: "1em", width: "0%" }
                      : val === "#"
                      ? { height: "1em", width: "3%" }
                      : { height: "1em" }
                    : val === ""
                    ? { height: "1em", width: "0%" }
                    : val === "#"
                    ? { height: "1em", width: "3%" }
                    : val === "Times"
                    ? { height: "1em", width: "25%" }
                    : val === "Name"
                    ? { height: "1em", width: "20%" }
                    : val == "Average" || val == "Single"
                    ? { height: "1em", width: "10%" }
                    : { height: "1em" };
                  let thStyle = style1;
                  return (
                    <th style={thStyle} key={idx}>
                      <Stack direction="row" alignItems="flex-end">
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
                  : result.single;
                result.times = isfmc
                  ? result.times.map((res) => res.split(".")[0])
                  : result.times;
                return (
                  <tr key={idx}>
                    <td style={{ height: "1em", width: "1%" }}></td>
                    <td style={{ height: "1em", width: "3%" }}>{idx + 1}.</td>
                    <td style={{ height: "1em" }}>
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
                      <td style={{ height: "1em" }}>
                        <span
                          className={`fi fi-${result.country_iso2.toLowerCase()}`}
                        />
                        &nbsp;&nbsp;{result.country_name}
                      </td>
                    )}
                    {isOverall ? (
                      <td style={{ height: "1em" }}>
                        <b>{result.score}</b>
                      </td>
                    ) : (
                      <>
                        <td style={{ height: "1em" }}>
                          <b>
                            {!averageFirst ? result.single : result.average}
                          </b>
                        </td>
                        <td style={{ height: "1em" }}>
                          {averageFirst ? result.single : result.average}
                        </td>
                        <td style={{ height: "1em" }}>
                          {result.times?.join(", ")}
                        </td>
                      </>
                    )}
                    <td style={{ height: "1em", width: "0%" }}></td>
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
