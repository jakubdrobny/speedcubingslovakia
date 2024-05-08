import { Card, Grid, Stack, Table, Typography } from "@mui/joy";

import { ProfileTypePersonalBests } from "../../Types";

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
      <Card sx={{ padding: "0.4em 0.5em" }}>
        <Table>
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
                    ...(idx === 0 || idx === 10
                      ? { width: idx === 0 ? "0.5%" : "2%" }
                      : (idx >= 2 && idx < 5) || (idx > 7 && idx < 10)
                      ? { width: "4%" }
                      : idx === 1
                      ? { width: "28%" }
                      : idx === 7
                      ? { width: "7%" }
                      : {}),
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
                  <td style={{ ...goodHeight, width: "2%" }}></td>
                  <td
                    style={{
                      ...goodHeight,
                      display: "flex",
                      alignItems: "center",
                      height: "1.5em",
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
                    "",
                  ].map((val, idx) => (
                    <td
                      key={idx}
                      style={
                        val !== ""
                          ? { ...right, width: idx === 0 ? "3%" : "auto" }
                          : { ...goodHeight, width: "2%" }
                      }
                    >
                      {idx === 3 || idx === 4 ? (
                        <b>
                          {idx === 3 && entry.eventIconcode === "333fm"
                            ? val.split(".")[0]
                            : val}
                        </b>
                      ) : (
                        val
                      )}
                    </td>
                  ))}
                </tr>
              ))}
          </tbody>
        </Table>
      </Card>
    </Stack>
  );
};

export default ProfilePersonalBests;
