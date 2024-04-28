import { Card, Table } from "@mui/joy";

import { ProfileTypeBasics } from "../../Types";

const ProfileBasicsDetailsTable: React.FC<{ basics: ProfileTypeBasics }> = ({
  basics,
}) => {
  const center: React.CSSProperties = { textAlign: "center", height: "1em" };
  return (
    <Card sx={{ padding: "0.25em 0.5em" }}>
      <Table>
        <thead>
          <tr>
            {[
              "Region",
              "WCA ID",
              "Sex",
              "Competitions",
              "Completed solves",
            ].map((columnTitle) => (
              <th style={center}>
                <b>{columnTitle}</b>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          <tr>
            <td style={center}>
              <span className={`fi fi-${basics.region.iso2.toLowerCase()}`} />
              &nbsp;&nbsp;{basics.region.name}
            </td>
            {[
              basics.wcaid,
              basics.sex,
              basics.noOfCompetitions,
              basics.completedSolves,
            ].map((columnContent) => (
              <td style={center}>{columnContent}</td>
            ))}
          </tr>
        </tbody>
      </Table>
    </Card>
  );
};

export default ProfileBasicsDetailsTable;
