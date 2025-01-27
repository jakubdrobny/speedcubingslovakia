import { Button } from "@mui/joy";
import { Link } from "react-router-dom";
import { saveCurrentLocation } from "../../utils/utils";

const PleaseLoginButton = () => {
  return (
    <Button
      variant="soft"
      color="warning"
      component={Link}
      to={import.meta.env.VITE_WCA_GET_CODE_URL || ""}
      onClick={() => saveCurrentLocation(window.location.pathname)}
    >
      Login to subscribe to positions on map
    </Button>
  );
};

export default PleaseLoginButton;
