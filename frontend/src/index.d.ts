import { HTMLAttributes } from "react";

declare global {
    namespace JSX {
      interface IntrinsicElements {
        "scramble-display": React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
      }
    }
  }