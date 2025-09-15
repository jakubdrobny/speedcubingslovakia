import {
  Autocomplete,
  Box,
  Button,
  Card,
  CircularProgress,
  Divider,
  Stack,
  Typography,
} from "@mui/joy";
import {
  findDuplicateUser,
  getError,
  getSearchUsers,
  isObjectEmpty,
  mergeUsers,
  renderResponseError,
} from "../../utils/utils";
import { LoadingState, ManageUser, User } from "../../Types";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const MergeUsers = () => {
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<ManageUser[]>([]);
  const [selectedUser, setSelectedUser] = useState<ManageUser | null>(null);
  const [duplicateUser, setDuplicateUser] = useState<User | null>(null);
  const [noDuplicateFound, setNoDuplicateFound] = useState(false);

  useEffect(() => {
    if (searchQuery.length < 2) {
      setSearchResults([]);
      return;
    }
    const delayDebounceFn = setTimeout(() => {
      getSearchUsers(searchQuery).then(setSearchResults);
    }, 300);

    return () => clearTimeout(delayDebounceFn);
  }, [searchQuery]);

  const handleFindDuplicate = () => {
    if (!selectedUser) return;
    setLoadingState({ isLoading: true, error: {} });
    setDuplicateUser(null);
    setNoDuplicateFound(false);

    findDuplicateUser(selectedUser.id)
      .then((res) => {
        setDuplicateUser(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        const error = getError(err);
        if (error.status === 404) {
          setNoDuplicateFound(true);
          setLoadingState({ isLoading: false, error: {} });
        } else {
          setLoadingState({ isLoading: false, error });
        }
      });
  };

  const handleMerge = () => {
    if (!selectedUser || !duplicateUser) return;
    setLoadingState({ isLoading: true, error: {} });

    const oldUserId = Math.min(selectedUser.id, duplicateUser.id);
    const newUserId = Math.max(selectedUser.id, duplicateUser.id);

    mergeUsers(oldUserId, newUserId)
      .then(() => {
        alert("Merge successful!");
        setSelectedUser(null);
        setDuplicateUser(null);
        setNoDuplicateFound(false);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) =>
        setLoadingState({ isLoading: false, error: getError(err) }),
      );
  };

  const renderUserCard = (user: ManageUser | User, title: string) => {
    const wcaId = "wca_id" in user ? user.wca_id : user.wcaid;
    return (
      <Card variant="outlined">
        <Typography level="h4">{title}</Typography>
        <Typography>
          ID: <b>{user.id}</b>
        </Typography>
        <Typography>
          Name: <Link to={`/profile/${wcaId || user.name}`}>{user.name}</Link>
        </Typography>
        <Typography>WCA ID: {wcaId || "N/A"}</Typography>
      </Card>
    );
  };

  return (
    <Stack spacing={2} sx={{ margin: "1em" }}>
      <Typography level="h2" className="bottom-divider">
        Merge duplicate users
      </Typography>

      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}

      <Card>
        <Stack spacing={2}>
          <Typography level="h3">1. Select a User</Typography>
          <Autocomplete
            placeholder="Search for a user..."
            options={searchResults}
            getOptionLabel={(option) => `${option.name} (${option.id})`}
            onInputChange={(_, newValue) => setSearchQuery(newValue)}
            onChange={(_, newValue) => {
              setSelectedUser(newValue);
              setDuplicateUser(null);
              setNoDuplicateFound(false);
            }}
            isOptionEqualToValue={(option, value) => option.id === value.id}
          />
          {selectedUser && (
            <Button
              onClick={handleFindDuplicate}
              loading={loadingState.isLoading}
            >
              Find Potential Duplicate
            </Button>
          )}
        </Stack>
      </Card>

      {loadingState.isLoading && <CircularProgress />}

      {noDuplicateFound && (
        <Card color="warning" variant="soft">
          <Typography>
            No potential duplicate account found for this user.
          </Typography>
        </Card>
      )}

      {selectedUser && duplicateUser && (
        <Card>
          <Stack spacing={2}>
            <Typography level="h3">2. Confirm Merge</Typography>
            <Typography>
              A potential duplicate has been found. All data from the user with
              the lower ID will be moved to the user with the higher ID, and the
              old user will be deleted. This action cannot be undone.
            </Typography>
            <Box display="flex" gap={2}>
              {renderUserCard(selectedUser, "Selected User")}
              {renderUserCard(duplicateUser, "Found Duplicate")}
            </Box>
            <Divider />
            <Button
              color="danger"
              onClick={handleMerge}
              loading={loadingState.isLoading}
            >
              Confirm and Merge Users
            </Button>
          </Stack>
        </Card>
      )}
    </Stack>
  );
};

export default MergeUsers;
