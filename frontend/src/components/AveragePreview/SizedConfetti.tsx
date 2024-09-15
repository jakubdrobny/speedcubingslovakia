import React from "react";
import ReactConfetti from "react-confetti";
import useWindowSize from "../../utils/useWindowSize";

const SizedConfetti = React.forwardRef<HTMLCanvasElement, any>(
  (passedProps, ref) => {
    const { width, height } = useWindowSize();
    return (
      <ReactConfetti width={width} height={height} {...passedProps} ref={ref} />
    );
  }
);

export default SizedConfetti;
