import { Button, ButtonGroup } from "@mui/joy";

import { Link } from "react-router-dom";

const Dashboard = () => {
  return (
    <div style={{ margin: "1em" }}>
      <ButtonGroup>
        <Button component={Link} to="/admin/manage-roles" color="primary">
          Manage roles
        </Button>
        <Button component={Link} to="/competition/create" color="primary">
          Create competition
        </Button>
        <Button component={Link} to="/results/edit" color="primary">
          Edit results
        </Button>
      </ButtonGroup>
    </div>
  );
};

export default Dashboard;
