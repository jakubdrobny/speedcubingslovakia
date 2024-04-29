import {
  Button,
  ButtonGroup,
  Card,
  List,
  ListItem,
  Select,
  Stack,
  Table,
  Typography,
} from "@mui/joy";
import {
  CompetitionEvent,
  RankingsEntry,
  RegionSelectGroup,
} from "../../Types";

import { Link } from "react-router-dom";
import React from "react";

const defaultRegionGroup = "World;World";

const RankingsTable: React.FC<{
  rankings: RankingsEntry[];
  single: boolean;
}> = ({ rankings, single }) => {
  const columnNames = (() => {
    let columnNames = ["", "#", "Name", "Result", "Represeting", "Competition"];
    if (!single) columnNames.push("Times");
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
                      ? "2%"
                      : idx === 1
                      ? "4%"
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
          {rankings.map((ranking, idx) => (
            <tr key={idx}>
              <td style={{ height: "1em", width: "2%" }}></td>
              <td style={{ height: "1em", width: "4%" }}>{idx + 1}.</td>
              <td style={{ height: "1em", width: "10%" }}>
                {ranking.username}
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
                <td style={{ height: "1em" }}>{ranking.times.join(", ")}</td>
              )}
              <td style={{ height: "1em", width: "2%" }}></td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Card>
  );
};

export default RankingsTable;
