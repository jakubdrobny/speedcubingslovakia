import { Button, Card, Switch, Table, Typography } from "@mui/joy";
import { Navigate, useNavigate } from "react-router-dom";
import { getUsers, updateUsers } from "../../utils";
import { useContext, useEffect, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import { AuthContextType } from "../../Types";
import { User } from "../../Types";

const ManageRoles = () => {
    const [users, setUsers] = useState<User[]>([]);
    const { authState } = useContext(AuthContext) as AuthContextType
    const navigate = useNavigate();

    useEffect(() => {
        if (!authState.authenticated || !authState.admin) {
            navigate("/");
        }
        
        getUsers()
            .then(res => setUsers(res))
            .catch(console.error);
    }, []);

    const handleUserRoleChange = (userid: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
        const checked: boolean = e.target.checked;
        setUsers(users.map((u: User): User => u.id === userid ? {...u, isadmin: checked} : {...u}));
    }

    const handleUserRolesSubmit = () => {
        updateUsers(users);
    }

    return (
        <div style={{marginTop: "2em"}}>
            <Card>
                <div style={{borderBottom: "1px solid #CDD7E1"}}>
                    <Typography fontWeight={"bold"}>Manage Roles</Typography>
                </div>
                <div>
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
                                        <Switch checked={user.isadmin} onChange={handleUserRoleChange(user.id)} />
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </Table>
                    <Button sx={{marginTop: "1em"}} onClick={() => handleUserRolesSubmit()}>Save</Button>
                </div>
            </Card>
        </div>
    )
}

export default ManageRoles