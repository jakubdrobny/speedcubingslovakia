import { Card, Table } from "@mui/joy";

import { Link } from "react-router-dom";
import { RankingsEntry } from "../../Types";
import React from "react";

const RankingsTable: React.FC<{
  rankings: RankingsEntry[];
  single: boolean;
  loading: boolean;
  isfmc: boolean;
}> = ({ rankings, single, loading, isfmc }) => {
  const columnNames = (() => {
    let columnNames = ["", "#", "Name", "Result", "Represeting", "Competition"];
    if (!single) columnNames.push(isfmc ? "Moves" : "Times");
    columnNames.push("");
    return columnNames;
  })();

  return (
    <Card sx={{ margin: 0, padding: 0 }}>
      <Table>
        <thead>
          <tr>
            {columnNames.map((val, idx) => (
              <th
                key={idx}
                style={{
                  height: "1em",
                  width:
                    idx === 0 || idx === columnNames.length - 1
                      ? "0.2%"
                      : idx === 1
                      ? "1%"
                      : idx == 2
                      ? "10%"
                      : "auto",
                  textAlign: idx === 3 ? "right" : "left",
                }}
              >
                <b>{val}</b>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {loading ? (
            <></>
          ) : (
            rankings.map((ranking, idx) => {
              ranking.result =
                isfmc && single ? ranking.result.split(".")[0] : ranking.result;
              ranking.times = isfmc
                ? ranking.times.map((res) => res.split(".")[0])
                : ranking.times;
              return (
                <tr key={idx}>
                  <td style={{ height: "1em", width: "2%" }}></td>
                  <td style={{ height: "1em", width: "4%" }}>{idx + 1}.</td>
                  <td style={{ height: "1em", width: "10%" }}>
                    <Link
                      to={`/profile/${ranking.wca_id}`}
                      style={{
                        color: "#0B6BCB",
                        textDecoration: "none",
                        fontWeight: 555,
                      }}
                    >
                      {ranking.username}
                    </Link>
                  </td>
                  <td style={{ height: "1em", textAlign: "right" }}>
                    <b>{ranking.result}</b>
                  </td>
                  <td style={{ height: "1em" }}>
                    <span
                      className={`fi fi-${ranking.country_iso2.toLowerCase()}`}
                    />
                    &nbsp;&nbsp;{ranking.country_name}
                  </td>
                  <td style={{ height: "1em" }}>
                    <Link to={`/competition/${ranking.competitionId}`}>
                      {ranking.competitionName}
                    </Link>
                  </td>
                  {!single && (
                    <td style={{ height: "1em" }}>
                      {ranking.times.join(", ")}
                    </td>
                  )}
                  <td style={{ height: "1em", width: "2%" }}></td>
                </tr>
              );
            })
          )}
        </tbody>
      </Table>
    </Card>
  );
};

export default RankingsTable;
