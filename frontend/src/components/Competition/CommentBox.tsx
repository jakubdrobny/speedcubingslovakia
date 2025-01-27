import { Card, Textarea } from "@mui/joy";

import { CompetitionContext } from "../../context/CompetitionContext";
import { CompetitionContextType } from "../../Types";
import { useContext } from "react";
import { emojify } from "node-emoji";

const CommentBox: React.FC<{ disabled: boolean }> = ({ disabled }) => {
  const { setCurrentResults, currentResultsRef } = useContext(
    CompetitionContext,
  ) as CompetitionContextType;

  const handleCommentChange = (newComment: string) => {
    setCurrentResults({
      ...currentResultsRef.current,
      comment: newComment,
    });
  };

  return (
    <Card>
      <h3 style={{ textAlign: "center", margin: "0.25em 0" }}>Comment:</h3>
      <Textarea
        value={emojify(currentResultsRef.current.comment)}
        onChange={(e) => handleCommentChange(e.target.value)}
        placeholder="Enter a comment to your solutions..."
        minRows={4}
        style={{ marginBottom: "1.25em" }}
        disabled={disabled}
      />
    </Card>
  );
};

export default CommentBox;
