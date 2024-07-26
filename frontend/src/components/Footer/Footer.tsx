import { Facebook, GitHub, Mail } from "@mui/icons-material";
import { Grid, Link, Stack, Typography } from "@mui/joy";

import DiscordIcon from "./DiscordIcon";

const Footer = () => {
  return (
    <Grid
      xs={12}
      sx={{
        position: "fixed",
        bottom: 0,
        width: "100%",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        padding: "1em",
        borderTop: "1px solid #D3D3D3",
        zIndex: 1000,
        backgroundColor: "white",
      }}
      component={Stack}
      direction="column"
    >
      <Typography>
        <b>Contact Us / Bug reporting:</b>
      </Typography>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          paddingTop: "0.25em",
          backgroundColor: "white",
          zIndex: 1000,
        }}
      >
        <Stack direction="row" spacing={2}>
          <Stack spacing={1} direction="row" alignItems="center">
            <Mail />
            <Link
              href="mailto:speedcubingsvk@gmail.com"
              style={{ color: "#0B6BCB", textDecoration: "none" }}
            >
              Email
            </Link>
          </Stack>
          <Stack spacing={1} direction="row" alignItems="center">
            <Facebook />
            <Link
              href="https://www.facebook.com/speedcubingslovakia"
              style={{ color: "#0B6BCB", textDecoration: "none" }}
            >
              Our page!
            </Link>
          </Stack>
          <Stack spacing={1} direction="row" alignItems="center">
            <DiscordIcon />
            <Link
              href="https://discord.com/invite/vKQs7htk"
              style={{ color: "#0B6BCB", textDecoration: "none" }}
            >
              Join us!
            </Link>
          </Stack>
          <Stack spacing={1} direction="row" alignItems="center">
            <GitHub />
            <Link
              href="https://github.com/jakubdrobny/speedcubingslovakia"
              style={{ color: "#0B6BCB", textDecoration: "none" }}
            >
              Github
            </Link>
          </Stack>
        </Stack>
      </div>
    </Grid>
  );
};

export default Footer;
