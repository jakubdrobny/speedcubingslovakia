import { Button, ButtonGroup } from "@mui/joy"
import { CompetitionContextType, CompetitionEvent } from "../../Types"

import { CompetitionContext } from "./CompetitionContext"
import { useContext } from "react"

export const EventSelector = () => {
    const { competitionState, updateCurrentEvent } = useContext(CompetitionContext) as CompetitionContextType

    return (
        <ButtonGroup style={{padding: '1em'}}>
            {competitionState.events.map((e: CompetitionEvent, idx: number) => {
                return (
                    <Button 
                        key={idx} 
                        onClick={() => updateCurrentEvent(idx)}
                        variant={idx === competitionState.currentEventIdx ? "solid" : "soft"}
                        color="primary"
                    >
                        <span className={`cubing-icon event-${e.iconcode}`}>&ensp;</span>
                        {e.displayname}
                    </Button>
                )
            })}
        </ButtonGroup>
    )
}