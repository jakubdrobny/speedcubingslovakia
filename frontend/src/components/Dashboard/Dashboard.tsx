import { Button, ButtonGroup } from "@mui/joy"
import { Link, Navigate } from "react-router-dom"

import { AuthContext } from "../../context/AuthContext"
import { AuthContextType } from "../../Types"
import { useContext } from "react"

const Dashboard = () => {
    const { authState } = useContext(AuthContext) as AuthContextType

    if (!authState.authenticated || !authState.admin) {
        return <Navigate to="/" />
    }

    return (
        <div style={{margin: '1em'}}>
            <ButtonGroup>
                <Button component={Link} to="/admin/manage-roles" color="primary">
                    Manage roles
                </Button>
                <Button component={Link} to="/competition/create" color="primary">
                    Create competition
                </Button>
                <Button component={Link} to="/results/edit" color="primary">
                    Edit results
                </Button>
            </ButtonGroup>
        </div>
    )
}

export default Dashboard