import {
  ComposableMap,
  Geographies,
  Geography,
  Graticule,
  Sphere,
  ZoomableGroup,
} from "react-simple-maps";
import { FeatureCollection, GeoJsonObject } from "geojson";
import {
  GetMapData,
  getError,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import { Stack, Typography } from "@mui/joy";
import { useEffect, useState } from "react";

import LoadingComponent from "../Loading/LoadingComponent";
import { LoadingState } from "../../Types";
import { Tooltip } from "react-tooltip";
import { scaleLinear } from "d3-scale";

const Users = () => {
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [mapData, setMapData] = useState<FeatureCollection>();
  const [tooltipContent, setTooltipContent] = useState<string>("");
  const colorScale = scaleLinear<string, string>()
    .domain([0, 1])
    .range(["#F5F4F6", "rgb(38, 62, 89)"]);
  const [maxNoOfCompetitors, setMaxNoOfCompetitors] = useState<number>(0);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    console.log("loading");
    GetMapData()
      .then((res: FeatureCollection) => {
        let _maxNoOfCompetitors = 0;
        res.features.map((f) =>
          Math.max(maxNoOfCompetitors, f.properties?.users.length)
        );
        setMaxNoOfCompetitors(_maxNoOfCompetitors);

        setMapData(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  }, []);

  return (
    <Stack sx={{ margin: "1em" }}>
      <Typography level="h2">Users</Typography>
      {loadingState.isLoading ? (
        <LoadingComponent title="Loading map data..." />
      ) : !isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <>
          <div
            data-tooltip-id="my-tooltip"
            data-tooltip-content=""
            data-tooltip-float={true}
            data-tooltip-place="bottom-start"
            data-tooltip-offset={10}
          >
            <ComposableMap
              projectionConfig={{
                scale: 147,
              }}
            >
              <ZoomableGroup>
                <Sphere
                  id="sphere-id"
                  fill="transparent"
                  stroke="#ebe8eb"
                  strokeWidth={0.5}
                />
                <Graticule
                  fill="transparent"
                  stroke="#ebe8eb"
                  strokeWidth={0.5}
                />
                <Geographies geography={mapData}>
                  {({ geographies }) =>
                    geographies.map((geo) => {
                      return (
                        <Geography
                          key={geo.rsmKey}
                          geography={geo}
                          stroke="#FFFFFF"
                          strokeWidth={0.5}
                          onMouseEnter={() => {
                            setTooltipContent(
                              `<span className={fi fi-${geo.properties.countryIso2.toLowerCase()}}/>&nbsp;&nbsp;${
                                geo.properties.countryName
                              }`
                            );
                          }}
                          onMouseLeave={() => {
                            setTooltipContent("");
                          }}
                          className="geo-no-outline"
                          style={{
                            default: {
                              fill: colorScale(
                                geo.properties.users.length / maxNoOfCompetitors
                              ),
                            },
                            hover: {
                              fill: "#F53",
                            },
                          }}
                        />
                      );
                    })
                  }
                </Geographies>
              </ZoomableGroup>
            </ComposableMap>
          </div>
          <Tooltip id="my-tooltip" noArrow>
            {tooltipContent}
          </Tooltip>
        </>
      )}
    </Stack>
  );
};

export default Users;
