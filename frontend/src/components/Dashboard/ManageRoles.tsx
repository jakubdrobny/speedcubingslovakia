import {
  Button,
  Card,
  CircularProgress,
  Switch,
  Table,
  Typography,
} from "@mui/joy";
import { ManageRolesUser, ResponseError } from "../../Types";
import {
  getError,
  getManageUsers,
  renderResponseError,
  updateUserRoles,
} from "../../utils/utils";
import { useEffect, useState } from "react";

const ManageRoles = () => {
  const [users, setUsers] = useState<ManageRolesUser[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>();
  const [error, setError] = useState<ResponseError>();

  useEffect(() => {
    setIsLoading(true);

    getManageUsers()
      .then((res: ManageRolesUser[]) => {
        setUsers(res);
        setIsLoading(false);
      })
      .catch((err) => {
        setIsLoading(false);
        setError(getError(err));
      });
  }, []);

  const handleUserRoleChange =
    (userid: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const checked: boolean = e.target.checked;
      setUsers(
        users.map(
          (u: ManageRolesUser): ManageRolesUser =>
            u.id === userid ? { ...u, isadmin: checked } : { ...u },
        ),
      );
    };

  const handleUserRolesSubmit = () => {
    setIsLoading(true);
    updateUserRoles(users)
      .then((res: ManageRolesUser[]) => {
        setUsers(res);
        setIsLoading(false);
      })
      .catch((err) => {
        setIsLoading(false);
        setError(getError(err));
      });
  };

  return (
    <div style={{ marginTop: "2em" }}>
      <Card>
        <div style={{ borderBottom: "1px solid #CDD7E1" }}>
          <Typography fontWeight={"bold"}>Manage Roles</Typography>
        </div>
        <div>
          {error ? (
            renderResponseError(error)
          ) : (
            <>
              <Table borderAxis="both" size="md">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>IsAdmin</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((user) => (
                    <tr key={user.id}>
                      <td>{user.name}</td>
                      <td>
                        {isLoading ? (
                          <CircularProgress />
                        ) : (
                          <Switch
                            checked={user.isadmin}
                            onChange={handleUserRoleChange(user.id)}
                          />
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </Table>
              {isLoading ? (
                <CircularProgress />
              ) : (
                <Button
                  sx={{ marginTop: "1em" }}
                  onClick={() => handleUserRolesSubmit()}
                >
                  Save
                </Button>
              )}
            </>
          )}
        </div>
      </Card>
    </div>
  );
};

export default ManageRoles;
