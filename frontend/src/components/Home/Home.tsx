import "../../styles/Home.css";

import { Typography } from "@mui/joy";

const Home = () => {
  return (
    <div
      style={{
        width: "100%",
        textAlign: "center",
      }}
    >
      <Typography
        level="h1"
        className="welcome-title"
        sx={{
          width: "100%",
          zIndex: -1,
        }}
      >
        Welcome to Speedcubing Slovakia!
      </Typography>
      <img
        src="/speedcubingslovakialogo256.png"
        style={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          zIndex: -2,
          opacity: 0.5,
        }}
        alt="SpeedcubingSlovakia Logo"
      />
    </div>
  );
};

export default Home;
