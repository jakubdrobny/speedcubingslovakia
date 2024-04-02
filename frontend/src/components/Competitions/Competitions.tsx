import { Button, ButtonGroup, Card, Stack, Typography } from "@mui/joy";
import { CompetitionData, FilterValue } from "../../Types";
import { useEffect, useState } from "react";

import FormControl from '@mui/joy/FormControl';
import { Link } from "react-router-dom";
import Table from '@mui/joy/Table';
import { loadCompetitionData } from "../../utils";

const Competitions = () => {
    const [competitionData, setCompetitionData] = useState<CompetitionData[]>([])
    const [filterValue, setFilterValue] = useState<FilterValue>(FilterValue.Current)

    useEffect(() => {
        loadCompetitionData(filterValue).then(setCompetitionData).catch(console.error);
    }, [filterValue]);

    const formatDate = (date: Date): String => date.toLocaleDateString() + " " + date.toLocaleTimeString()

    const handleFilterChange = (newFilterValue: FilterValue) => setFilterValue(newFilterValue);

    return (
        <Card sx={{margin: "1em 0"}}>
            <Typography level="h1" sx={{margin: "0.25em 0", borderBottom: '1px solid '}}>Online competitions</Typography>
            <Stack direction="row" sx={{display: 'flex', alignItems: 'center'}} spacing={1}>
                <Typography level="title-md">Filters:</Typography>
                <ButtonGroup>
                {Object.keys(FilterValue).map((key) => {
                    const filterKey = key as keyof typeof FilterValue;
                    return (
                        <Button 
                            key={FilterValue[filterKey]} 
                            onClick={() => handleFilterChange(FilterValue[filterKey])}
                            variant={filterValue === FilterValue[filterKey] ? "solid" : "outlined"}
                            color="primary"
                        >
                            {FilterValue[filterKey]}
                        </Button>
                    );
                })}
                </ButtonGroup>
            </Stack>
            <Table aria-label="basic table">
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
                                    <Link to={`/competition/${competition.id}`}>{competition.name}</Link>
                                </td>
                                <td>{formatDate(competition.startdate)}</td>
                                <td>{formatDate(competition.enddate)}</td>
                            </tr>
                        );
                    })}
                </tbody>
            </Table>
        </Card>
    );
}

export default Competitions;