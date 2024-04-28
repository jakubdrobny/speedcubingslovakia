import { Card, IconButton, Stack, Table, Typography } from "@mui/joy";
import { Link, Navigate } from "react-router-dom";

import { ProfileTypeResultHistory } from "../../Types";
import { useState } from "react";

const ProfileResultsHistory: React.FC<{
  resultsHistory: ProfileTypeResultHistory[];
}> = ({ resultsHistory }) => {
  const goodHeight: React.CSSProperties = { height: "1em" };
  const left: React.CSSProperties = { textAlign: "left", ...goodHeight };
  const center: React.CSSProperties = { textAlign: "center", ...goodHeight };
  const right: React.CSSProperties = { textAlign: "right", ...goodHeight };
  const [currentHistoryIdx, setCurrentHistoryIdx] = useState(0);

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
              {["", "Competition", "Place", "Single", "Average", "Solves"].map(
                (columnTitle, idx) => (
                  <th
                    key={columnTitle}
                    style={{
                      ...(columnTitle === ""
                        ? goodHeight
                        : idx < 3
                        ? left
                        : idx < 5
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
                                : idx < 5
                                ? "10%"
                                : "40%",
                          }),
                    }}
                  >
                    <b>{columnTitle}</b>
                  </th>
                )
              )}
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
              {[0, 1, 2, 3].map((val) => (
                <td key={val} style={goodHeight}></td>
              ))}
            </tr>
            {resultsHistory[currentHistoryIdx].history.map((entry, idx) => (
              <tr key={idx}>
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
                {[entry.place, entry.single, entry.average].map((val, idx1) => (
                  <td key={idx1 + 100000} style={idx1 === 0 ? left : right}>
                    {idx1 === 0 ? val : <b>{val}</b>}
                  </td>
                ))}
                <td style={center}>{entry.solves.join(", ")}</td>
                <td style={goodHeight}></td>
              </tr>
            ))}
          </tbody>
        </Table>
      </Card>
    </Stack>
  );
};

export default ProfileResultsHistory;
