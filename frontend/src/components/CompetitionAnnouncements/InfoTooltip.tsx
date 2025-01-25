import { HelpOutline } from "@mui/icons-material";
import { IconButton, Typography, Tooltip, Stack } from "@mui/joy";
import { Dispatch, SetStateAction } from "react";

const InfoTooltipTitle: React.FC<{ notInsideTooltip: boolean }> = ({
  notInsideTooltip,
}) => {
  return (
    <Stack
      spacing={1}
      className="shadow-md rounded-md p-4 bg-[#e3effb]"
      sx={notInsideTooltip ? { fontSize: "0.875rem" } : {}}
    >
      <Typography
        sx={{
          fontWeight: "625",
          textDecoration: "underline",
        }}
        level="h4"
      >
        Tired of checking the WCA website for new competitions?
      </Typography>
      <Typography fontSize="1em">
        Subscribe to our <b>newsletter</b> to receive emails when{" "}
        <b>new WCA competitions</b> are <b>announced</b> in any country of your
        choice.{" "}
      </Typography>
      <Typography fontSize="1em">You can subscribe to:</Typography>
      <Typography fontSize="1em" sx={{ pl: 2 }}>
        - any <b>country</b> (or <b>multiple countries</b>) of your choice
        <br />- any <b>position with radius around it</b> (or multiple such
        positions)
        <br />- any (or multiple) <b>US states</b>
      </Typography>
      <Typography fontSize="1em">
        You can unsubscribe from any of the mentioned above at any time.
      </Typography>
      <Typography fontWeight="bold" fontSize="1em">
        Enjoy :)
      </Typography>
    </Stack>
  );
};

const InfoTooltip: React.FC<{
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}> = ({ open, setOpen }) => {
  return (
    <Tooltip
      variant="soft"
      color="primary"
      title={<InfoTooltipTitle notInsideTooltip={false} />}
      open={open}
      disableInteractive={false}
      enterTouchDelay={0}
      enterDelay={0}
      leaveDelay={0}
      sx={{ padding: 0, margin: 0 }}
    >
      <IconButton
        onMouseEnter={() => setOpen(true)}
        onMouseLeave={() => setOpen(false)}
        onTouchStart={() => setOpen((p) => !p)}
      >
        <HelpOutline fontSize="small" />
      </IconButton>
    </Tooltip>
  );
};

export { InfoTooltip, InfoTooltipTitle };
