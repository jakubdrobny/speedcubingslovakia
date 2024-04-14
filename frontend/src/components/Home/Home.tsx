import "../../styles/Home.css";

import { Typography } from "@mui/joy";

const Home = () => {
  return (
    <div style={{ width: "100%", textAlign: "center" }}>
      <Typography level="h1" className="welcome-title" sx={{ width: "100%" }}>
        Welcome to Speedcubing Slovakia!
      </Typography>
    </div>
  );
};

export default Home;
