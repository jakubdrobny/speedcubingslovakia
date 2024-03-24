import { Button, ButtonGroup, Typography } from "@mui/joy";
import { CompetitionData, FilterValue } from "../Types";
import { useEffect, useState } from "react";

import FormControl from '@mui/joy/FormControl';
import FormLabel from '@mui/joy/FormLabel';
import { Link } from "react-router-dom";
import Table from '@mui/joy/Table';
import { loadCompetitionData } from "../utils";

const Competitions = () => {
    const [competitionData, setCompetitionData] = useState<CompetitionData[]>([])
    const [filterValue, setFilterValue] = useState<FilterValue>(FilterValue.Current)

    useEffect(() => {
        loadCompetitionData(filterValue).then(setCompetitionData).catch(console.error);
    }, [filterValue]);

    const formatDate = (date: Date): String => date.toLocaleDateString() + " " + date.toLocaleTimeString()

    const handleFilterChange = (newFilterValue: FilterValue) => setFilterValue(newFilterValue);

    return (
        <div>
            <h1>Online competitions</h1>
            <FormControl>
                <FormLabel>Filters:</FormLabel>
                <ButtonGroup variant="outlined">
                {Object.keys(FilterValue).map((key) => {
                    const filterKey = key as keyof typeof FilterValue;
                    return (
                        <Button 
                            key={FilterValue[filterKey]} 
                            onClick={() => handleFilterChange(FilterValue[filterKey])}
                            variant={filterValue === FilterValue[filterKey] ? "solid" : "plain"}
                            color="primary"
                        >
                            {FilterValue[filterKey]}
                        </Button>
                    );
                })}
                </ButtonGroup>
            </FormControl>

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
        </div>
    );
}

export default Competitions;