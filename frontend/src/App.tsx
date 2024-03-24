import { Grid, List, ListItemButton, ListItemDecorator } from '@mui/joy';
import { Link, Route, Routes } from 'react-router-dom';

import Home from './routes/Home'
import LanguageIcon from '@mui/icons-material/Language';
import OnlineCompetitions from './routes/OnlineCompetitions';

const App = () => {
    return (
        <Grid container>
            <Grid xs={12} borderBottom={"2px solid lightgrey"}>
                <List style={{ display: 'flex', flexDirection: 'row', padding: 20, justifyContent: 'space-around'}}>
                    <ListItemButton component={Link} to="/">
                        Speedcubing Slovakia
                    </ListItemButton>
                    <ListItemButton component={Link} to="/online-competitions">
                        <ListItemDecorator>
                            <LanguageIcon/>
                        </ListItemDecorator>
                        Online Competitions
                    </ListItemButton>
                </List>
            </Grid>
            <Grid xs={0} md={1} lg={2}/>
            <Grid xs={12} md={10} lg={8}>
                <Routes>
                    <Route path="/" Component={Home} />
                    <Route path="/online-competitions" Component={OnlineCompetitions} />
                </Routes>
            </Grid>
            <Grid xs={0} md={1} lg={2}/>
        </Grid>
    );
}

export default App;
