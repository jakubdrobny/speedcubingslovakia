import { AddAlarm, EmojiEvents } from "@mui/icons-material";
import { Button, ButtonGroup } from "@mui/joy";

import { ResultsCompeteChoiceEnum } from "../../Types";
import { useSearchParams } from "react-router-dom";
import { RESULTS_COMPETE_CHOICE_QUERY_PARAM_NAME } from "../../constants";

const ResultsCompeteChoice: React.FC<{
  resultsCompeteChoice: ResultsCompeteChoiceEnum;
  setResultsCompeteChoice: (
    newResultsCompeteChoice: ResultsCompeteChoiceEnum,
  ) => void;
  loading: boolean;
}> = ({ resultsCompeteChoice, setResultsCompeteChoice, loading }) => {
  const [searchParams, setSearchParams] = useSearchParams();
  return (
    <ButtonGroup sx={{ pb: 1, flexWrap: "wrap", margin: 0 }}>
      <Button
        onClick={() => {
          searchParams.set(RESULTS_COMPETE_CHOICE_QUERY_PARAM_NAME, "results");
          setSearchParams(searchParams);
          setResultsCompeteChoice(ResultsCompeteChoiceEnum.Results);
        }}
        variant={
          resultsCompeteChoice === ResultsCompeteChoiceEnum.Results
            ? "solid"
            : "soft"
        }
        color="primary"
        disabled={loading}
      >
        <EmojiEvents />
        &nbsp; Results
      </Button>
      <Button
        onClick={() => {
          searchParams.set(RESULTS_COMPETE_CHOICE_QUERY_PARAM_NAME, "compete");
          setSearchParams(searchParams);
          setResultsCompeteChoice(ResultsCompeteChoiceEnum.Compete);
        }}
        variant={
          resultsCompeteChoice === ResultsCompeteChoiceEnum.Compete
            ? "solid"
            : "soft"
        }
        color="primary"
        disabled={loading}
      >
        <AddAlarm />
        &nbsp; Compete
      </Button>
    </ButtonGroup>
  );
};

export default ResultsCompeteChoice;
