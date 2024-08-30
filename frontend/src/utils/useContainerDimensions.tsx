import { useEffect } from "react";
import useState from "react-usestateref";

export const useContainerDimensions = (myRef: any) => {
  const [dimensions, setDimensions, dimensionsRef] = useState({
    width: 0,
    height: 0,
  });

  useEffect(() => {
    const getDimensions = () => ({
      width: myRef.current.offsetWidth,
      height: myRef.current.offsetHeight,
    });

    const handleResize = () => {
      setDimensions(getDimensions());
    };

    if (myRef.current) {
      setDimensions(getDimensions());
    }

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, [myRef]);

  return dimensionsRef;
};
