import { Grid, List, ListItemButton, ListItemDecorator } from '@mui/joy';
import { Link, Navigate, Route, Routes } from 'react-router-dom';

import Competition from './components/Competition/Competition';
import Competitions from './components/Competitions/Competitions';
import Home from './components/Home/Home'
import LanguageIcon from '@mui/icons-material/Language';
import NotFound from './components/NotFound/NotFound';

const App = () => {
    return (
        <Grid container>
            <Grid xs={12} borderBottom={"2px solid lightgrey"}>
                <List style={{ display: 'flex', flexDirection: 'row', padding: 20, justifyContent: 'space-around'}}>
                    <ListItemButton component={Link} to="/">
                        Speedcubing Slovakia
                    </ListItemButton>
                    <ListItemButton component={Link} to="/competitions">
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
                    <Route path="/competitions" Component={Competitions} />
                    <Route path="/competition/:id" Component={Competition} />
                    <Route path='/not-found' Component={NotFound} />
                    <Route path="*" element={ <Navigate to="/not-found" replace />} />
                </Routes>
            </Grid>
            <Grid xs={0} md={1} lg={2}/>
        </Grid>
    );
}

export default App;
