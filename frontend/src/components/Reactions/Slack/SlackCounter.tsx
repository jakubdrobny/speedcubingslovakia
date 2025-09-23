import React from "react";
import { CounterObject, groupBy, Hover } from "./helpers";
import SlackCounterGroup from "./SlackCounterGroup";
import SlackCSS from "./SlackCSS";

export interface SlackCounterProps {
  counters?: CounterObject[];
  user?: string;
  onSelect?: (emoji: string) => void;
  onAdd?: () => void;
}

export const SlackCounter = React.forwardRef<HTMLDivElement, SlackCounterProps>(
  (
    {
      counters = defaultProps.counters,
      user = defaultProps.user,
      onSelect = defaultProps.onSelect,
      onAdd = defaultProps.onAdd,
    },
    ref,
  ) => {
    const groups = groupBy(counters, "emoji");

    return (
      <>
        <SlackCSS />
        <Hover ref={ref} style={counterStyle}>
          {groups &&
            Object.keys(groups).map((emoji: string) => {
              const names = groups[emoji].map(({ by }: CounterObject) => {
                return by;
              });
              return (
                <div style={groupStyle} key={emoji}>
                  <SlackCounterGroup
                    emoji={emoji}
                    count={names.length}
                    names={names}
                    active={names.includes(user)}
                    onSelect={onSelect}
                  />
                </div>
              );
            })}
          <div style={addStyle} onClick={onAdd}>
            <SlackCounterGroup emoji={""} />
          </div>
        </Hover>
      </>
    );
  },
);

const defaultProps: Required<SlackCounterProps> = {
  counters: [
    {
      emoji: "👍",
      by: "Case Sandberg",
    },
    {
      emoji: "👎",
      by: "Charlie!!!!!",
    },
  ],
  user: "Charlie",
  onSelect: (emoji: string) => {
    console.log(emoji);
  },
  onAdd: () => {
    console.log("add");
  },
};

const counterStyle = {
  display: "flex",
};
const addStyle = {
  cursor: "pointer",
  fontFamily: "Slack",
  opacity: "1",
  transition: "opacity 0.1s ease-in-out",
};
const groupStyle = {
  marginRight: "4px",
};

export default SlackCounter;
