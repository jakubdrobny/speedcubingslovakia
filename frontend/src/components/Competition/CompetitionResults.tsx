import {
  Alert,
  Card,
  CircularProgress,
  Sheet,
  Table,
  Typography,
} from "@mui/joy";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { CompetitionContextType } from "../../Types";

const CompetitionResults = () => {
  const { competitionState, results, loadingState, fetchCompetitionResults } =
    useContext(CompetitionContext) as CompetitionContextType;
  const averageFirst = (() => {
    const format =
      competitionState?.events[competitionState?.currentEventIdx]?.format;
    if (!format || format?.length === 0) return true;
    return format[0] != "b";
  })();
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);

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
      "Times",
      "",
    ];
    if (!averageFirst)
      [columnNames[4], columnNames[5]] = [columnNames[5], columnNames[4]];
    if (windowWidth < 1000)
      columnNames = [...columnNames.slice(0, 3), ...columnNames.slice(4)];
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
                {columnNames().map((val, idx) => (
                  <th
                    style={
                      val === ""
                        ? { height: "1em", width: "0%" }
                        : val === "#"
                        ? { height: "1em", width: "3%" }
                        : val === "Times"
                        ? { height: "1em", width: "25%" }
                        : val === "Name"
                        ? { height: "1em", width: "20%" }
                        : val == "Average" || val == "Single"
                        ? { height: "1em", width: "10%" }
                        : { height: "1em" }
                    }
                    key={idx}
                  >
                    <b>{val}</b>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {results.map((result, idx) => (
                <tr key={idx}>
                  <td style={{ height: "1em", width: "1%" }}></td>
                  <td style={{ height: "1em", width: "3%" }}>{idx + 1}.</td>
                  <td style={{ height: "1em" }}>{result.username}</td>
                  {windowWidth >= 1000 && (
                    <td style={{ height: "1em" }}>
                      <span
                        className={`fi fi-${result.country_iso2.toLowerCase()}`}
                      />
                      &nbsp;&nbsp;{result.country_name}
                    </td>
                  )}
                  <td style={{ height: "1em" }}>
                    <b>{!averageFirst ? result.single : result.average}</b>
                  </td>
                  <td style={{ height: "1em" }}>
                    {averageFirst ? result.single : result.average}
                  </td>
                  <td style={{ height: "1em" }}>{result.times?.join(", ")}</td>
                  <td style={{ height: "1em", width: "0%" }}></td>
                </tr>
              ))}
            </tbody>
          </Table>
        </Card>
      )}
    </div>
  );
};

export default CompetitionResults;
