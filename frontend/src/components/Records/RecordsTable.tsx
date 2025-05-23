import { Card, Stack, Table, Typography } from "@mui/joy";
import { RecordsItem, RecordsItemEntry } from "../../Types";
import { getCubingIconClassName, reformatMultiTime } from "../../utils/utils";

import { Link } from "react-router-dom";
import React from "react";

const RecordsTable: React.FC<{
  recordItems: RecordsItem[];
}> = ({ recordItems }) => {
  const columnNames = (() => {
    let columnNames = ["Type", "Name", "Result", "Represeting", "Competition"];
    return columnNames;
  })();

  return (
    <Stack spacing={3} direction="column">
      {recordItems.map((item: RecordsItem, idx0) => {
        const isfmc = item.iconcode === "333fm";
        const ismbld = item.iconcode === "333mbf";
        return (
          <Stack direction="column" spacing={2} key={idx0}>
            <Typography
              level="h3"
              alignItems="center"
              sx={{ display: "flex", alignItems: "center", height: "1.5em" }}
            >
              <span
                className={`${getCubingIconClassName(
                  item.iconcode
                )} profile-cubing-icon-mock`}
              />
              &nbsp;{item.eventname}
            </Typography>
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
                    {columnNames
                      .concat(
                        item.entries[item.entries.length - 1].type !== "Single"
                          ? ["Solves"]
                          : []
                      )
                      .map((val, idx) => (
                        <th
                          key={idx}
                          style={{
                            height: "1em",
                            textAlign:
                              idx === 2
                                ? "right"
                                : val == "Solves"
                                ? "center"
                                : "left",
                          }}
                        >
                          <b>{val}</b>
                        </th>
                      ))}
                  </tr>
                </thead>
                <tbody>
                  {item.entries.map((itemEntry: RecordsItemEntry, idx) => {
                    itemEntry.result =
                      isfmc && itemEntry.type === "Single"
                        ? itemEntry.result.split(".")[0]
                        : ismbld
                        ? reformatMultiTime(itemEntry.result)
                        : itemEntry.result;
                    itemEntry.solves = isfmc
                      ? itemEntry.solves.map((res) => res.split(".")[0])
                      : ismbld
                      ? itemEntry.solves.map((res) => reformatMultiTime(res))
                      : itemEntry.solves;
                    return (
                      <tr key={idx}>
                        <td
                          style={{
                            height: "1em",
                            opacity:
                              idx > 0 &&
                              item.entries[idx - 1].type === itemEntry.type
                                ? 0.5
                                : 1,
                          }}
                        >
                          {itemEntry.type}
                        </td>
                        <td style={{ height: "1em" }}>
                          <Link
                            to={`/profile/${itemEntry.wcaId}`}
                            style={{
                              color: "#0B6BCB",
                              textDecoration: "none",
                              fontWeight: 555,
                            }}
                          >
                            {itemEntry.username}
                          </Link>
                        </td>
                        <td style={{ height: "1em", textAlign: "right" }}>
                          <b>{itemEntry.result}</b>
                        </td>
                        <td style={{ height: "1em" }}>
                          <span
                            className={`fi fi-${itemEntry.countryIso2.toLowerCase()}`}
                          />
                          &nbsp;&nbsp;{itemEntry.countryName}
                        </td>
                        <td style={{ height: "1em" }}>
                          <Link to={`/competition/${itemEntry.competitionId}`}>
                            {itemEntry.competitionName}
                          </Link>
                        </td>
                        {item.entries[item.entries.length - 1].type !==
                          "Single" && (
                          <td style={{ height: "1em", textAlign: "center" }}>
                            {itemEntry.solves.join(", ")}
                          </td>
                        )}
                      </tr>
                    );
                  })}
                </tbody>
              </Table>
            </Card>
          </Stack>
        );
      })}
    </Stack>
  );
};

export default RecordsTable;
