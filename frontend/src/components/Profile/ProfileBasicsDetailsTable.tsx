import { Card, Table } from "@mui/joy";

import { Link } from "react-router-dom";
import { ProfileTypeBasics } from "../../Types";

const ProfileBasicsDetailsTable: React.FC<{ basics: ProfileTypeBasics }> = ({
  basics,
}) => {
  const center: React.CSSProperties = { textAlign: "center", height: "1em" };
  return (
    <Card sx={{ padding: "0.25em 0.5em", whiteSpace: "nowrap" }}>
      <Table>
        <thead>
          <tr>
            {[
              "Region",
              "WCA ID",
              "Sex",
              "Competitions",
              "Completed solves",
            ].map((columnTitle, idx) => (
              <th key={idx} style={center}>
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
            ].map((columnContent, idx) => (
              <td key={idx} style={center}>
                {idx === 0 ? (
                  <Link
                    to={`https://worldcubeassociation.org/persons/${columnContent}`}
                    style={{ color: "#0B6BCB", textDecoration: "none" }}
                  >
                    {columnContent}
                  </Link>
                ) : (
                  columnContent
                )}
              </td>
            ))}
          </tr>
        </tbody>
      </Table>
    </Card>
  );
};

export default ProfileBasicsDetailsTable;
