import {
  Box,
  Button,
  Card,
  FormControl,
  FormHelperText,
  FormLabel,
  IconButton,
  Input,
  Stack,
  Typography,
} from "@mui/joy";
import { LoadingState, SearchUser } from "../../Types";
import { getError, getUsers, initialLoadingState } from "../../utils";

import { Link } from "react-router-dom";
import { Search } from "@mui/icons-material";
import { useState } from "react";

const Users = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [users, setUsers] = useState<SearchUser[]>([]);

  const searchForUsers = () => {
    setLoadingState({ isLoading: true, error: "" });

    getUsers(searchQuery)
      .then((res: SearchUser[]) => {
        setUsers(res);
        setLoadingState({ isLoading: false, error: "" });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  return (
    <Stack sx={{ margin: "1em 0" }} spacing={2}>
      <Typography level="h2">Find users</Typography>
      <FormControl disabled={loadingState.isLoading}>
        <FormLabel>Enter WCA ID or username:</FormLabel>
        <Input
          placeholder="Enter WCA ID or username... (eg. 2016DROB01)"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          startDecorator={
            <Button
              startDecorator={<Search />}
              variant="soft"
              color="neutral"
              onClick={searchForUsers}
            >
              Search
            </Button>
          }
        />
        <FormHelperText>Leave empty for all users.</FormHelperText>
      </FormControl>
      {users.map((u: SearchUser) => (
        <Card sx={{ display: "flex", flexDirection: "row" }}>
          <span style={{ fontWeight: "bold" }}>{u.username}: </span>
          <Link to={`/profile/${u.wcaid}`} style={{ textDecoration: "none" }}>
            {u.wcaid}
          </Link>
        </Card>
      ))}
    </Stack>
  );
};

export default Users;
