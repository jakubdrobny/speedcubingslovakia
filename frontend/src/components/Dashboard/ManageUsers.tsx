import { Card, Stack, Switch, Table, Typography } from "@mui/joy";
import { LoadingState, ManageUser } from "../../Types";
import {
  getError,
  getManageUsers,
  initialLoadingState,
  renderResponseError,
  updateUserRoles,
  isObjectEmpty,
} from "../../utils/utils";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import LoadingComponent from "../Loading/LoadingComponent";

const ManageRoles = () => {
  const [users, setUsers] = useState<ManageUser[]>([]);
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    getManageUsers()
      .then((res: ManageUser[]) => {
        setUsers(res.sort((u1: ManageUser, u2: ManageUser) => u2.id - u1.id));
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, []);

  const handleUserRoleChange =
    (users_idx: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const checked: boolean = e.target.checked;
      const user = users[users_idx];
      user.is_admin = checked;

      setLoadingState({ isLoading: true, error: {} });
      updateUserRoles(user)
        .then((_) => {
          setUsers(
            users.map(
              (u: ManageUser, idx: number): ManageUser =>
                idx === users_idx ? { ...u, is_admin: checked } : { ...u },
            ),
          );
          setLoadingState({ isLoading: false, error: {} });
        })
        .catch((err) => {
          setLoadingState({ isLoading: false, error: getError(err) });
        });
    };

  const columnNames = () => ["Order", "Name", "Country", "Is Admin?"];

  return (
    <Stack style={{ padding: 16 }} spacing={2}>
      <div style={{ borderBottom: "1px solid #CDD7E1" }}>
        <Typography fontWeight={"bold"} level="h2">
          Manage Roles
        </Typography>
      </div>

      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}
      {loadingState.isLoading && (!users || users.length === 0) ? (
        <LoadingComponent title="Fetching user data..." />
      ) : (
        <Card sx={{ padding: 0, margin: 0, overflowX: "auto" }}>
          <Table
            size="md"
            sx={{
              tableLayout: "auto",
              width: "100%",
              whiteSpace: "nowrap",
            }}
          >
            <thead>
              <tr>
                {columnNames().map((val, idx) => {
                  return (
                    <th
                      style={{
                        height: "1em",
                        maxWidth: "auto",
                        textAlign: idx === 0 ? "right" : "left",
                      }}
                      key={idx}
                    >
                      <b>{val}</b>
                    </th>
                  );
                })}
              </tr>
            </thead>
            <tbody>
              {users.map((user, idx) => {
                return (
                  <tr key={idx}>
                    <td style={{ height: "1em", textAlign: "right" }}>
                      {user.id}
                    </td>
                    <td style={{ height: "1em", textAlign: "left" }}>
                      <Link
                        to={`/profile/${user.wca_id ? user.wca_id : user.name}`}
                        style={{
                          color: "#0B6BCB",
                          textDecoration: "none",
                          fontWeight: 555,
                        }}
                      >
                        {user.name}
                      </Link>
                    </td>
                    <td style={{ height: "1em", textAlign: "left" }}>
                      <span
                        className={`fi fi-${user.country_iso2.toLowerCase()}`}
                      />
                      &nbsp;&nbsp;{user.country_name}
                    </td>
                    <td style={{ height: "1em", textAlign: "left" }}>
                      <Switch
                        disabled={loadingState.isLoading}
                        checked={user.is_admin}
                        onChange={handleUserRoleChange(idx)}
                      />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </Table>
        </Card>
      )}
    </Stack>
  );
};

export default ManageRoles;
