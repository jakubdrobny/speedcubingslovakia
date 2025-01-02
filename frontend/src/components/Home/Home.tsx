import "../../styles/Home.css";

import { Typography } from "@mui/joy";
import mainLogo from './speedcubingslovakialogo256.png'
import { useEffect, useState } from "react";


const Home = () => {
  const [logoReady, setLogoReady] = useState(false);

  useEffect(() => {
    const img = new Image();
    img.onload = () => setLogoReady(true);
    img.src = mainLogo;
  }, []);

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
      {logoReady && <img
        src={mainLogo}
        style={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          zIndex: -2,
          opacity: 0.5,
        }}
        alt="SpeedcubingSlovakia Logo"
      />}
    </div>
  );
};

export default Home;
