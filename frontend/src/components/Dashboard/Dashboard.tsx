import { Button, Stack } from "@mui/joy";

import { Link } from "react-router-dom";

const Dashboard = () => {
  return (
    <div style={{ margin: "1em" }}>
      <Stack direction="column" spacing={1}>
        <Button
          component={Link}
          to="/admin/stats"
          color="primary"
          variant="outlined"
        >
          Show stats
        </Button>
        <Button
          component={Link}
          to={import.meta.env.VITE_MONITORING_PATH}
          reloadDocument
          color="primary"
          variant="outlined"
        >
          Monitoring
        </Button>
        <Button
          component={Link}
          to="/admin/manage-roles"
          color="primary"
          variant="outlined"
        >
          Manage roles
        </Button>
        <Button
          component={Link}
          to="/competition/create"
          color="primary"
          variant="outlined"
        >
          Create competition
        </Button>
        <Button
          component={Link}
          to="/results/edit"
          color="primary"
          variant="outlined"
        >
          Edit results
        </Button>
        <Button
          component={Link}
          to="/announcement/create"
          color="primary"
          variant="outlined"
        >
          Create announcement
        </Button>
      </Stack>
    </div >
  );
};

export default Dashboard;
