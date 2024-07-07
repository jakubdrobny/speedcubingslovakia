import { Card, Table } from "@mui/joy";

import { Link } from "react-router-dom";
import { RankingsEntry } from "../../Types";
import React from "react";
import { reformatMultiTime } from "../../utils";

const RankingsTable: React.FC<{
  rankings: RankingsEntry[];
  single: boolean;
  loading: boolean;
  isfmc: boolean;
  ismbld: boolean;
}> = ({ rankings, single, loading, isfmc, ismbld }) => {
  const columnNames = (() => {
    let columnNames = ["#", "Name", "Result", "Represeting", "Competition"];
    if (!single) columnNames.push(isfmc ? "Moves" : "Times");
    return columnNames;
  })();

  return (
    <Card
      sx={{
        margin: 0,
        padding: 0,
        overflow: "auto",
        boxSizing: "border-box",
        MozBoxSizing: "border-box",
        WebkitBoxSizing: "border-box",
        width: `calc(100%-1em)`,
      }}
    >
      <Table
        sx={{
          //   width: "100%",
          tableLayout: "auto",
          whiteSpace: "nowrap",
          width: `calc(100%-1em)`,
        }}
      >
        <thead>
          <tr>
            {columnNames.map((val, idx) => (
              <th
                key={idx}
                style={{
                  height: "1em",
                  textAlign: idx === 0 || idx === 2 ? "right" : "left",
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
                isfmc && single
                  ? ranking.result.split(".")[0]
                  : ismbld
                  ? reformatMultiTime(ranking.result)
                  : ranking.result;
              ranking.times = isfmc
                ? ranking.times.map((res) => res.split(".")[0])
                : ismbld
                ? ranking.times.map((res) => reformatMultiTime(res))
                : ranking.times;
              return (
                <tr key={idx}>
                  <td style={{ height: "1em", textAlign: "right" }}>
                    {ranking.place}
                  </td>
                  <td style={{ height: "1em" }}>
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
