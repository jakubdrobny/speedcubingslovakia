import { AddAlarm, EmojiEvents } from "@mui/icons-material";
import { Button, ButtonGroup } from "@mui/joy";

import { ResultsCompeteChoiceEnum } from "../../Types";

const ResultsCompeteChoice: React.FC<{
  resultsCompeteChoice: ResultsCompeteChoiceEnum;
  setResultsCompeteChoice: (
    newResultsCompeteChoice: ResultsCompeteChoiceEnum
  ) => void;
  loading: boolean;
}> = ({ resultsCompeteChoice, setResultsCompeteChoice, loading }) => {
  return (
    <ButtonGroup style={{ padding: "1em", flexWrap: "wrap", margin: 0 }}>
      <Button
        onClick={() =>
          setResultsCompeteChoice(ResultsCompeteChoiceEnum.Results)
        }
        variant={
          resultsCompeteChoice === ResultsCompeteChoiceEnum.Results
            ? "solid"
            : "soft"
        }
        color="primary"
        loading={loading}
      >
        <EmojiEvents />
        &nbsp; Results
      </Button>
      <Button
        onClick={() =>
          setResultsCompeteChoice(ResultsCompeteChoiceEnum.Compete)
        }
        variant={
          resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete
            ? "solid"
            : "soft"
        }
        color="primary"
        loading={loading}
      >
        <AddAlarm />
        &nbsp; Compete
      </Button>
    </ButtonGroup>
  );
};

export default ResultsCompeteChoice;
