import { useEffect, useState } from "react";

import Table from '@mui/joy/Table';

type Competition = {
    name: String,
    startdate: Date,
    enddate: Date
}

const OnlineCompetitions = () => {
    const [competitionData, setCompetitionData] = useState<Competition[]>([])

    useEffect(() => {
        const loadCompetitionData = async (): Promise<Competition[]> => {
            const result: Competition[] = [];
            let startdate = new Date();
            let enddate = startdate;
            enddate.setDate(startdate.getDate() + 7);

            for (let i = 0; i < 3; i++) {
                if (i) {
                    startdate = result[result.length - 1].enddate
                    enddate = startdate
                    enddate.setDate(startdate.getDate() + 7)
                }

                result.push(
                    {
                        'name': 'Weekly Competition ' + (i + 1).toString(),
                        startdate,
                        enddate
                    }
                )
            }

            return result
        }

        loadCompetitionData().then(setCompetitionData).catch(console.error);
    }, )

    return (
        <div>
            <h1>Online competitions</h1>
            <Table aria-label="basic table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Start date</th>
                        <th>End date</th>
                    </tr>
                </thead>
                <tbody>
                    {competitionData.map((competition: Competition, index) => {
                        return (
                            <tr key={index}>
                                <td>{competition.name}</td>
                                <td>{competition.startdate.toISOString()}</td>
                                <td>{competition.enddate.toISOString()}</td>
                            </tr>
                        );
                    })}
                </tbody>
            </Table>
        </div>
    );
}

export default OnlineCompetitions;