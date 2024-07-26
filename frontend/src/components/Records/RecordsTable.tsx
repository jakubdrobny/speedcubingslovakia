import { Card, Stack, Table, Typography } from "@mui/joy";
import { RankingsEntry, RecordsItem } from "../../Types";

import { Link } from "react-router-dom";
import React from "react";
import { reformatMultiTime } from "../../utils";

const RecordsTable: React.FC<{
  recordItems: RecordsItem[];
  loading: boolean;
  isfmc: boolean;
  ismbld: boolean;
}> = ({ recordItems, loading, isfmc, ismbld }) => {
  const columnNames = (() => {
    let columnNames = [
      "Type",
      "Name",
      "Result",
      "Represeting",
      "Competition",
      "Solves",
    ];
    return columnNames;
  })();

  return (
    <Stack spacing={3} direction="column">
      {recordItems.map((item: RecordsItem) => (
        <Stack direction="column" spacing={2}>
        <Typography level="h3"><span className={`cubing-icon event-${item.iconcode} profile-cubing-icon-mock`} />&nbsp;{item.eventname}</Typography>
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
        </>
      ))}
    </Stack>
  );
};

export default RecordsTable;
