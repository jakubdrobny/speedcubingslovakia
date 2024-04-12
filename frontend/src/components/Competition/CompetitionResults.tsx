import {
  Alert,
  Card,
  CircularProgress,
  Sheet,
  Table,
  Typography,
} from "@mui/joy";
import { CompetitionContextType, CompetitionResult } from "../../Types";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "./CompetitionContext";
import { getCompetitionResults } from "../../utils";

const CompetitionResults = () => {
  const { competitionState } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const [results, setResult] = useState<CompetitionResult[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const averageFirst = (() => {
    const format =
      competitionState?.events[competitionState?.currentEventIdx]?.format;
    if (!format || format?.length === 0) return true;
    return format[0] != "b";
  })();

  useEffect(() => {
    setIsLoading(true);
    setError("");

    getCompetitionResults(
      competitionState.id,
      competitionState.events[competitionState.currentEventIdx]
    )
      .then((res) => {
        setResult(res);
        setIsLoading(false);
        setError("");
      })
      .catch((err) => {
        setIsLoading(false);
        setError(err.message);
      });
  }, [competitionState.currentEventIdx]);

  const columnNames = () => {
    return averageFirst
      ? ["", "#", "Name", "Country", "Average", "Single", "Times", ""]
      : ["", "#", "Name", "Country", "Single", "Average", "Times", ""];
  };

  return (
    <>
      {isLoading || competitionState.loadingState.results ? (
        <>
          <Typography level="h3" sx={{ display: "flex", alignItems: "center" }}>
            <CircularProgress />
            &nbsp; Loading results ...
          </Typography>
        </>
      ) : error ? (
        <Alert color="danger">{error}</Alert>
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
                        : { height: "1em" }
                    }
                    key={idx}
                  >
                    {val}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {results.map((result, idx) => (
                <tr>
                  <td style={{ height: "1em", width: "1%" }}></td>
                  <td style={{ height: "1em", width: "3%" }}>{idx + 1}.</td>
                  <td style={{ height: "1em" }}>{result.username}</td>
                  <td style={{ height: "1em" }}>
                    <span
                      className={`fi fi-${result.country_iso2.toLowerCase()}`}
                    />
                    &nbsp;&nbsp;{result.country_name}
                  </td>
                  <td style={{ height: "1em" }}>
                    {!averageFirst ? result.single : result.average}
                  </td>
                  <td style={{ height: "1em" }}>
                    {averageFirst ? result.average : result.single}
                  </td>
                  <td style={{ height: "1em" }}>{result.times?.join(", ")}</td>
                  <td style={{ height: "1em", width: "0%" }}></td>
                </tr>
              ))}
            </tbody>
          </Table>
        </Card>
      )}
    </>
  );
};

export default CompetitionResults;
