import { Grid, List, ListItemButton, ListItemDecorator } from "@mui/joy";
import { Link, Navigate, Route, Routes } from "react-router-dom";
import { authorizeAdmin, setBearerIfPresent } from "./utils";
import { useContext, useEffect } from "react";

import { AuthContext } from "./context/AuthContext";
import { AuthContextType } from "./Types";
import Competition from "./components/Competition/Competition";
import CompetitionEdit from "./components/Dashboard/CompetitionEdit";
import Competitions from "./components/Competitions/Competitions";
import Dashboard from "./components/Dashboard/Dashboard";
import Home from "./components/Home/Home";
import LanguageIcon from "@mui/icons-material/Language";
import { ListAlt } from "@mui/icons-material";
import LogIn from "./components/Login/LogIn";
import ManageRoles from "./components/Dashboard/ManageRoles";
import NotFound from "./components/NotFound/NotFound";
import ProfileListItem from "./components/Profile/ProfileListItem";
import ProtectedRoute from "./components/Login/ProtectedRoute";
import ResultsEdit from "./components/Dashboard/ResultsEdit";
import WCALogoNoText from "./images/WCALogoNoText";

const App = () => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;

  useEffect(() => {
    setBearerIfPresent(authState.token);

    if (authState.token) {
      authorizeAdmin()
        .then((_) => setAuthState({ ...authState, isadmin: true }))
        .catch((_) => null);
    }
  }, []);

  return (
    <Grid container>
      <Grid
        xs={0}
        md={1}
        lg={2}
        borderBottom={"2px solid lightgrey"}
        width={"100%"}
      />
      <Grid
        xs={12}
        md={10}
        lg={8}
        borderBottom={"2px solid lightgrey"}
        width={"100%"}
      >
        <List
          style={{
            display: "flex",
            flexDirection: "row",
            padding: 20,
            width: "100%",
          }}
        >
          <Grid sx={{ display: "flex", justifyContent: "center" }}>
            <ListItemButton component={Link} to="/">
              Speedcubing Slovakia
            </ListItemButton>
          </Grid>
          <Grid
            sx={{ display: "flex", justifyContent: "flex-end", width: "100%" }}
          >
            <ListItemButton component={Link} to="/competitions">
              <ListItemDecorator>
                <LanguageIcon />
              </ListItemDecorator>
              Online Competitions
            </ListItemButton>
            {authState.isadmin && (
              <ListItemButton component={Link} to="/admin/dashboard">
                <ListItemDecorator>
                  <ListAlt />
                </ListItemDecorator>
                Dashboard
              </ListItemButton>
            )}
            {authState.token ? (
              <ProfileListItem />
            ) : (
              <ListItemButton
                component={Link}
                to={process.env.REACT_APP_WCA_GET_CODE_URL || ""}
              >
                <ListItemDecorator>
                  <WCALogoNoText />
                </ListItemDecorator>
                Log In
              </ListItemButton>
            )}
          </Grid>
        </List>
      </Grid>
      <Grid
        xs={0}
        md={1}
        lg={2}
        borderBottom={"2px solid lightgrey"}
        width={"100%"}
      />
      <Grid xs={0} md={1} lg={2} />
      <Grid xs={12} md={10} lg={8}>
        <Routes>
          <Route path="/" Component={Home} />
          <Route path="/competitions" Component={Competitions} />
          <Route path="/competition/:id" Component={Competition} />
          <Route path="/not-found" Component={NotFound} />
          <Route path="/login" Component={LogIn} />
          <Route Component={ProtectedRoute}>
            <Route
              path="/competition/:id/edit"
              Component={() => <CompetitionEdit edit={true} />}
            />
            <Route
              path="/competition/create"
              Component={() => <CompetitionEdit edit={false} />}
            />
            <Route path="/admin/dashboard" Component={Dashboard} />
            <Route path="/admin/manage-roles" Component={ManageRoles} />
            <Route path="/results/edit" Component={ResultsEdit} />
          </Route>
          <Route path="*" element={<Navigate to="/not-found" replace />} />
        </Routes>
      </Grid>
      <Grid xs={0} md={1} lg={2} />
    </Grid>
  );
};

export default App;
