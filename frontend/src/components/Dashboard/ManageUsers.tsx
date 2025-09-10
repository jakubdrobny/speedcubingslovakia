import { Card, CircularProgress, Switch, Table, Typography } from "@mui/joy";
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

const ManageRoles = () => {
  const [users, setUsers] = useState<ManageUser[]>([]);
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    getManageUsers()
      .then((res: ManageUser[]) => {
        setUsers(res);
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
    <div style={{ marginTop: "2em" }}>
      <Card>
        <div style={{ borderBottom: "1px solid #CDD7E1" }}>
          <Typography fontWeight={"bold"}>Manage Roles</Typography>
        </div>
        <div>
          {!isObjectEmpty(loadingState.error) &&
            renderResponseError(loadingState.error)}{" "}
          {loadingState.isLoading && (!users || users.length === 0) ? (
            <>
              <CircularProgress />
              &nbsp;Loading...
            </>
          ) : (
            <>
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
                            to={`/profile/${user.wca_id}`}
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
                          &nbsp;&nbsp;{user.country}
                        </td>
                        <td style={{ height: "1em", textAlign: "left" }}>
                          <Switch
                            checked={user.is_admin}
                            onChange={handleUserRoleChange(idx)}
                          />
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </Table>
            </>
          )}
        </div>
      </Card>
    </div>
  );
};

export default ManageRoles;
