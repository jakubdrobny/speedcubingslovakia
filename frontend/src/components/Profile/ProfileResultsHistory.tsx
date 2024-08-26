import { Card, Stack, Table, Typography } from "@mui/joy";
import { getCubingIconClassName, isBoX, reformatMultiTime } from "../../utils";

import { Link } from "react-router-dom";
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
  const isfmc = resultsHistory[currentHistoryIdx]?.eventIconcode === "333fm";
  const ismbld = resultsHistory[currentHistoryIdx]?.eventIconcode === "333mbf";
  const isboX = isBoX(resultsHistory[currentHistoryIdx]?.eventFormat);

  const getColumnNames = () => {
    let columnNames = [
      "",
      "Competition",
      "Place",
      "Single",
      "",
      "Average",
      "",
      resultsHistory[currentHistoryIdx]?.eventIconcode === "333fm"
        ? "Moves"
        : "Solves",
    ];

    if (columnNames.includes("Average") && isboX)
      columnNames.splice(
        columnNames.findIndex((x) => x === "Average"),
        2
      );

    return columnNames;
  };

  return (
    <Stack spacing={2} sx={{ whiteSpace: "nowrap" }}>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <Typography level="h3">Results</Typography>
      </div>
      <Card sx={{ padding: "0.4em 0.5em", gap: 0 }}>
        <div
          style={{
            textAlign: "center",
            padding: 0,
            margin: 0,
            marginTop: "5px",
            overflowX: "auto",
          }}
        >
          {resultsHistory.map((entry, idx) => (
            <span
              key={idx}
              className={`${getCubingIconClassName(
                entry.eventIconcode
              )} profile-cubing-icon-mock`}
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
        <div style={{ overflowX: "auto" }}>
          <Table sx={{ tableLayout: "auto" }}>
            <thead>
              <tr>
                {getColumnNames().map((columnTitle, idx) => {
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
                <td style={{ ...goodHeight }}></td>
                <td
                  style={{
                    ...left,
                    display: "flex",
                    alignItems: "center",
                    height: "1.5em",
                  }}
                >
                  <span
                    className={getCubingIconClassName(
                      resultsHistory[currentHistoryIdx].eventIconcode
                    )}
                  />
                  &nbsp;{resultsHistory[currentHistoryIdx].eventName}
                </td>
                {(isboX ? [0, 1, 2, 3] : [0, 1, 2, 3, 4, 5]).map((val) => (
                  <td key={val + 10} style={goodHeight}></td>
                ))}
              </tr>
              {resultsHistory[currentHistoryIdx].history.map((entry, idx) => (
                <tr key={idx + 1000}>
                  <td style={{ ...goodHeight }}></td>
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
                  {(isboX
                    ? [entry.place, entry.single, entry.singleRecord]
                    : [
                        entry.place,
                        entry.single,
                        entry.singleRecord,
                        entry.average,
                        entry.averageRecord,
                      ]
                  ).map((val, idx1) => {
                    return (
                      <td key={idx1 + 100000} style={idx1 === 0 ? left : right}>
                        {idx1 === 0 ? (
                          val
                        ) : (
                          <b
                            style={{
                              color:
                                idx1 === 1
                                  ? entry.singleRecordColor
                                  : idx1 === 3
                                  ? entry.averageRecordColor
                                  : "black",
                            }}
                          >
                            {idx1 === 1 && isfmc
                              ? val.split(".")[0]
                              : ismbld
                              ? reformatMultiTime(val)
                              : val}
                          </b>
                        )}
                      </td>
                    );
                  })}
                  <td style={center}>
                    {entry.solves
                      ? (isfmc
                          ? entry.solves.map((x) => x.split(".")[0])
                          : ismbld
                          ? entry.solves.map((x) => reformatMultiTime(x))
                          : entry.solves
                        ).join(", ")
                      : entry.solves}
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </Card>
    </Stack>
  );
};

export default ProfileResultsHistory;
