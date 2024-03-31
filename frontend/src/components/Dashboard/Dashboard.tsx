import { Accordion, AccordionDetails, AccordionGroup, AccordionSummary, Button, Switch, Table, Typography } from "@mui/joy"
import { AuthContextType, User } from "../../Types"
import { getUsers, updateUsers } from "../../utils"
import { useContext, useEffect, useState } from "react"

import { AuthContext } from "../../context/AuthContext"
import { Navigate } from "react-router-dom"
import { accordionClasses } from '@mui/joy/Accordion';
import { accordionSummaryClasses } from '@mui/joy/AccordionSummary';

enum Panel {
    ManageRoles, None
}

const Dashboard = () => {
    const { authState } = useContext(AuthContext) as AuthContextType
    const [panel, setPanel] = useState<Panel>(Panel.None)
    const [users, setUsers] = useState<User[]>([]);

    useEffect(() => {
        getUsers()
            .then(res => setUsers(res))
            .catch(console.error);
    }, []);

    if (!authState.authenticated || !authState.admin) {
        return <Navigate to="/" />
    }

    const handlePanelChange = (panel: Panel) => (e: React.SyntheticEvent, newExpanded: boolean) => {
        setPanel(newExpanded ? panel : Panel.None);
    }

    const handleUserRoleChange = (userid: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
        const checked: boolean = e.target.checked;
        setUsers(users.map((u: User): User => u.id === userid ? {...u, isadmin: checked} : {...u}));
    }

    const handleUserRolesSubmit = () => {
        updateUsers(users);
    }

    return (
        <div style={{margin: '1em'}}>
            <AccordionGroup
                size="lg"
                sx={{
                    [`& .${accordionClasses.root}`]: {
                      marginTop: '0.5rem',
                      transition: '0.1s ease',
                      '& button:not([aria-expanded="true"])': {
                        transition: '0.1s ease',
                        paddingBottom: '0.625rem',
                      },
                    },
                    [`& .${accordionClasses.root}.${accordionClasses.expanded}`]: {
                      bgcolor: '',
                      borderRadius: 'md',
                      border: '1px solid',
                      borderColor: 'background.level2',
                      boxShadow: (theme) => `${theme.vars.palette.divider} 2px 2px`,
                    },
                    '& [aria-expanded="true"]': {
                        borderBottom: '1px solid'
                    },
                    [`& .${accordionSummaryClasses.button}:hover`]: {
                        borderTopRightRadius: "7px",
                        borderTopLeftRadius: "7px"
                    },
                  }}
            >
                <Accordion expanded={panel === Panel.ManageRoles} onChange={handlePanelChange(Panel.ManageRoles)}>
                    <AccordionSummary>
                        <Typography>Manage Roles</Typography>
                    </AccordionSummary>
                    <AccordionDetails sx={{paddingTop: 1, margin: 0}}>
                    <Table borderAxis={"both"} size="sm" sx={{margin: "1em 0 1em "}}>
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>IsAdmin</th>
                            </tr>
                        </thead>
                        <tbody>
                            {users.map((user) => {
                                console.log(user.isadmin); return (
                                <tr key={user.id}>
                                    <td>{user.name}</td>
                                    <td>
                                        <Switch checked={user.isadmin} onChange={handleUserRoleChange(user.id)} />
                                    </td>
                                </tr>
                            )})}
                        </tbody>
                    </Table>
                    <Button sx={{marginTop: "1em"}} onClick={() => handleUserRolesSubmit()}>Save</Button>
                    </AccordionDetails>
                </Accordion>
            </AccordionGroup>
        </div>
    )
}

export default Dashboard