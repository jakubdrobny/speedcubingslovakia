import { Card, IconButton, Stack, Table, Typography } from "@mui/joy";
import { Link, Navigate } from "react-router-dom";

import { ProfileTypeResultHistory } from "../../Types";
import { reformatMultiTime } from "../../utils";
import { useState } from "react";

const ProfileResultsHistory: React.FC<{
  resultsHistory: ProfileTypeResultHistory[];
}> = ({ resultsHistory }) => {
  const goodHeight: React.CSSProperties = { height: "1em" };
  const left: React.CSSProperties = { textAlign: "left", ...goodHeight };
  const center: React.CSSProperties = { textAlign: "center", ...goodHeight };
  const right: React.CSSProperties = { textAlign: "right", ...goodHeight };
  const [currentHistoryIdx, setCurrentHistoryIdx] = useState(0);
  const isfmc = resultsHistory[currentHistoryIdx]?.eventIconcode === "333fm";
  const ismbld = resultsHistory[currentHistoryIdx]?.eventIconcode === "333mbf";

  const getColumnNames = () => {
    let columnNames = [
      "",
      "Competition",
      "Place",
      "Single",
      "Average",
      resultsHistory[currentHistoryIdx]?.eventIconcode === "333fm"
        ? "Moves"
        : "Solves",
    ];

    if (columnNames.includes("Average") && ismbld)
      columnNames = columnNames.filter((c) => c !== "Average");

    return columnNames;
  };

  return (
    <Stack spacing={2}>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <Typography level="h3">Results</Typography>
      </div>
      <Card sx={{ padding: "0.4em 0.5em", gap: 0 }}>
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            padding: 0,
            margin: 0,
            marginTop: "5px",
          }}
        >
          {resultsHistory.map((entry, idx) => (
            <span
              key={idx}
              className={`cubing-icon event-${entry.eventIconcode} profile-cubing-icon-mock`}
              onClick={() => setCurrentHistoryIdx(idx)}
              style={{
                padding: "0 0.25em",
                fontSize: "1.75em",
                color: currentHistoryIdx === idx ? "#0B6BCB" : "",
                cursor: "pointer",
              }}
            />
          ))}
        </div>
        <Table>
          <thead>
            <tr>
              {getColumnNames().map((columnTitle, idx) => {
                console.log({
                  ...(columnTitle === ""
                    ? goodHeight
                    : idx < 3
                    ? left
                    : idx < 5 - (ismbld ? 1 : 0)
                    ? right
                    : center),
                  ...(idx === 0
                    ? { width: "1%" }
                    : {
                        width:
                          idx === 1
                            ? "30%"
                            : idx === 2
                            ? "4%"
                            : idx < 5 - (ismbld ? 1 : 0)
                            ? "10%"
                            : "40%",
                      }),
                });
                return (
                  <th
                    key={idx}
                    style={{
                      ...(columnTitle === ""
                        ? goodHeight
                        : idx < 3
                        ? left
                        : idx < 5 - (ismbld ? 1 : 0)
                        ? right
                        : center),
                      ...(idx === 0
                        ? { width: "1%" }
                        : {
                            width:
                              idx === 1
                                ? "30%"
                                : idx === 2
                                ? "4%"
                                : idx < 5 - (ismbld ? 1 : 0)
                                ? "10%"
                                : ismbld
                                ? "30%"
                                : "40%",
                          }),
                    }}
                  >
                    <b>{columnTitle}</b>
                  </th>
                );
              })}
            </tr>
          </thead>
          <tbody>
            <tr>
              <td style={{ ...goodHeight, width: "2%" }}></td>
              <td
                style={{
                  ...left,
                  display: "flex",
                  alignItems: "center",
                  height: "1.5em",
                }}
              >
                <span
                  className={`cubing-icon event-${resultsHistory[currentHistoryIdx].eventIconcode}`}
                />
                &nbsp;{resultsHistory[currentHistoryIdx].eventName}
              </td>
              {(ismbld ? [0, 1, 2] : [0, 1, 2, 3]).map((val) => (
                <td key={val + 10} style={goodHeight}></td>
              ))}
            </tr>
            {resultsHistory[currentHistoryIdx].history.map((entry, idx) => (
              <tr key={idx + 1000}>
                <td style={{ ...goodHeight, width: "2%" }}></td>
                <td
                  style={{
                    ...goodHeight,
                  }}
                >
                  <Link
                    to={`/competition/${entry.competitionId}`}
                    style={{ textDecoration: "none", color: "#0B6BCB" }}
                  >
                    <b>{entry.competitionName}</b>
                  </Link>
                </td>
                {(ismbld
                  ? [entry.place, entry.single]
                  : [entry.place, entry.single, entry.average]
                ).map((val, idx1) => (
                  <td key={idx1 + 100000} style={idx1 === 0 ? left : right}>
                    {idx1 === 0 ? (
                      val
                    ) : (
                      <b>
                        {idx1 === 1 && isfmc
                          ? val.split(".")[0]
                          : ismbld
                          ? reformatMultiTime(val)
                          : val}
                      </b>
                    )}
                  </td>
                ))}
                <td style={center}>
                  {(isfmc
                    ? entry.solves.map((x) => x.split(".")[0])
                    : ismbld
                    ? entry.solves.map((x) => reformatMultiTime(x))
                    : entry.solves
                  ).join(", ")}
                </td>
                {!ismbld && <td style={goodHeight}></td>}
              </tr>
            ))}
          </tbody>
        </Table>
      </Card>
    </Stack>
  );
};

export default ProfileResultsHistory;
