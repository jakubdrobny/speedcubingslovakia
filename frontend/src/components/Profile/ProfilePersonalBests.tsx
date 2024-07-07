import { Card, Grid, Stack, Table, Typography } from "@mui/joy";

import { ProfileTypePersonalBests } from "../../Types";
import { reformatMultiTime } from "../../utils";

const ProfilePersonalBests: React.FC<{ pbs: ProfileTypePersonalBests[] }> = ({
  pbs,
}) => {
  const goodHeight: React.CSSProperties = { height: "1em" };
  const right: React.CSSProperties = { textAlign: "right", ...goodHeight };

  return (
    <Stack spacing={2}>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <Typography level="h3">Personal Best Records</Typography>
      </div>
      <Card sx={{ padding: "0.4em 0.5em", overflowX: "auto" }}>
        <Table sx={{ tableLayout: "auto" }}>
          <thead>
            <tr>
              {[
                "",
                "Event",
                "NR",
                "CR",
                "WR",
                "Single",
                "Average",
                "WR",
                "CR",
                "NR",
                "",
              ].map((columnTitle, idx) => (
                <th
                  key={idx}
                  style={{
                    ...(columnTitle !== "Event" ? right : goodHeight),
                  }}
                >
                  <b>{columnTitle}</b>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {pbs &&
              pbs.map((entry) => (
                <tr key={entry.eventIconcode}>
                  <td style={{ ...goodHeight }}></td>
                  <td
                    style={{
                      ...goodHeight,
                      display: "flex",
                      alignItems: "center",
                      height: "1.5em",
                      whiteSpace: "nowrap",
                    }}
                  >
                    <span
                      className={`cubing-icon event-${entry.eventIconcode}`}
                    />
                    &nbsp;{entry.eventName}
                  </td>
                  {[
                    entry.single.nr,
                    entry.single.cr,
                    entry.single.wr,
                    entry.single.value,
                    entry.average.value,
                    entry.average.wr,
                    entry.average.cr,
                    entry.average.nr,
                  ].map((val, idx) => (
                    <td
                      key={idx}
                      style={
                        val !== ""
                          ? {
                              ...right,
                              whiteSpace: "nowrap",
                            }
                          : { ...goodHeight, whiteSpace: "nowrap" }
                      }
                    >
                      {idx === 3 || idx === 4 ? (
                        <b>
                          {idx === 3 && entry.eventIconcode === "333fm"
                            ? val.split(".")[0]
                            : entry.eventIconcode === "333mbf"
                            ? reformatMultiTime(val)
                            : val}
                        </b>
                      ) : (
                        val
                      )}
                    </td>
                  ))}
                  <td style={{ ...goodHeight }}></td>
                </tr>
              ))}
          </tbody>
        </Table>
      </Card>
    </Stack>
  );
};

export default ProfilePersonalBests;
