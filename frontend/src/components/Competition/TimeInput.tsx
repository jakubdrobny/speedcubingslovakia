import { CompetitionContextType, ResultEntry } from "../../Types";

import { CompetitionContext } from "./CompetitionContext";
import { Input } from "@mui/joy";
import { useContext } from "react";

const TimeInput = () => {
    const { competitionState, updateSolve, saveResults } = useContext(CompetitionContext) as CompetitionContextType;
    const solveProp: keyof ResultEntry = `solve${competitionState.currentSolveIdx+1}` as keyof ResultEntry;
    const formattedTime = competitionState.results[solveProp].toString();
    
    const reformatTime = (oldFormattedTime: string, added: boolean = false): string => {
        if (added) {
            let idx = 0;
            while (idx < oldFormattedTime.length && /^\D/.test(oldFormattedTime[idx]) || oldFormattedTime[idx] === '0')
                idx++;
            oldFormattedTime = oldFormattedTime.slice(idx);
        }

        const matchedDigits = oldFormattedTime.match(/\d+/g);
        let digits = !matchedDigits ? '' : matchedDigits.join('');
        if (digits.length < 3)
            digits = digits.padStart(3, '0');

        let newFormattedTime = `${digits[digits.length - 1]}${digits[digits.length - 2]}.`;
        let idx = digits.length - 3;
        while (idx >= 0) {
            newFormattedTime += digits[idx--];
            if (idx >= 0)
                newFormattedTime += digits[idx--];
            if (idx >= 0)
                newFormattedTime += ':';
        }

        newFormattedTime = newFormattedTime.split('').reverse().join('');

        return newFormattedTime;
    }

    const handleTimeInputChange = (e: React.FormEvent<HTMLInputElement>) => {
        const newValue = e.currentTarget.value;
        
        // character deleted
        if (newValue.length + 1 === formattedTime.length) {
            if (newValue.endsWith('N')) {
                updateSolve("0.00");
                return;
            } else {
                updateSolve(reformatTime(newValue));
            }
        } else {
            console.log(newValue);
            if (newValue.endsWith("d")) {
                updateSolve("DNF");
            } else if (newValue.endsWith("s")) {
                updateSolve("DNS");
            } else if (/\d$/.test(newValue.slice(-1))) {
                updateSolve(reformatTime(newValue, true));
            } else {
                updateSolve("DNF");
            }
        }
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            saveResults();
        }
    }

    return (
        <Input
            size="lg"
            placeholder="Enter your time or solution..."
            sx={{ marginBottom: 2, marginTop: 2}}
            value={formattedTime}
            onChange={handleTimeInputChange}
            onKeyDown={handleKeyDown}
        />
    )
}

export default TimeInput;