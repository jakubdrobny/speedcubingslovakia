import {
  Alert,
  Button,
  ButtonGroup,
  Card,
  CircularProgress,
  Stack,
  Typography,
} from "@mui/joy";
import { CompetitionData, FilterValue } from "../../Types";
import { formatDate, loadFilteredCompetitions } from "../../utils";
import { useEffect, useState } from "react";

import { Link } from "react-router-dom";
import Table from "@mui/joy/Table";
import { getError } from "../../utils";

const Competitions = () => {
  const [competitionData, setCompetitionData] = useState<CompetitionData[]>([]);
  const [filterValue, setFilterValue] = useState<FilterValue>(
    FilterValue.Current
  );
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<any>("");

  useEffect(() => {
    setIsLoading(true);
    loadFilteredCompetitions(filterValue)
      .then((res) => {
        setIsLoading(false);
        setCompetitionData(res);
      })
      .catch((err) => {
        setIsLoading(false);
        setError(getError(err));
      });
  }, [filterValue]);

  const handleFilterChange = (newFilterValue: FilterValue) =>
    setFilterValue(newFilterValue);

  return (
    <Card sx={{ margin: "1em 0.5em" }}>
      <Typography
        level="h2"
        sx={{ margin: "0.25em 0", borderBottom: "1px solid #CDD7E1" }}
      >
        Online competitions
      </Typography>
      <Stack
        direction="row"
        sx={{ display: "flex", alignItems: "center" }}
        spacing={1}
      >
        <Typography level="title-md">Filters:</Typography>
        <ButtonGroup>
          {Object.keys(FilterValue).map((key) => {
            const filterKey = key as keyof typeof FilterValue;
            return (
              <Button
                key={FilterValue[filterKey]}
                onClick={() => handleFilterChange(FilterValue[filterKey])}
                variant={
                  filterValue === FilterValue[filterKey] ? "solid" : "outlined"
                }
                color="primary"
              >
                {FilterValue[filterKey]}
              </Button>
            );
          })}
        </ButtonGroup>
      </Stack>
      {error ? (
        <Alert color="danger">{error}</Alert>
      ) : isLoading ? (
        <CircularProgress />
      ) : (
        <div style={{ overflowX: "auto" }}>
          <Table
            aria-label="basic table"
            sx={{
              tableLayout: "auto",
              width: "100%",
              whiteSpace: "nowrap",
            }}
          >
            <thead>
              <tr>
                <th>Name</th>
                <th>Start date</th>
                <th>End date</th>
              </tr>
            </thead>
            <tbody>
              {competitionData.map((competition: CompetitionData, index) => {
                return (
                  <tr key={index}>
                    <td>
                      <Link to={`/competition/${competition.id}`}>
                        {competition.name}
                      </Link>
                    </td>
                    <td>{formatDate(competition.startdate)}</td>
                    <td>{formatDate(competition.enddate)}</td>
                  </tr>
                );
              })}
            </tbody>
          </Table>
        </div>
      )}
    </Card>
  );
};

export default Competitions;
