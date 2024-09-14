import {
  ComposableMap,
  Geographies,
  Geography,
  Graticule,
  Sphere,
  ZoomableGroup,
} from "react-simple-maps";
import {
  GetMapData,
  getError,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
} from "../../utils/utils";
import { Stack, Typography } from "@mui/joy";

import { EmojiEvents } from "@mui/icons-material";
import { FeatureCollection } from "geojson";
import LoadingComponent from "../Loading/LoadingComponent";
import { LoadingState } from "../../Types";
import { Tooltip } from "react-tooltip";
import { scaleLinear } from "d3-scale";
import { useEffect } from "react";
import useState from "react-usestateref";

const Users = () => {
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [mapData, setMapData] = useState<FeatureCollection>();
  const [tooltipContent, setTooltipContent] = useState<any>("");
  const defaultColorScale = scaleLinear<string, string>()
    .domain([0, 1])
    .range(["rgba(11, 107, 203, 0.1)", "rgba(11, 107, 203, 1)"]);
  const hoverColorScale = scaleLinear<string, string>()
    .domain([0, 1])
    .range(["rgba(196, 28, 28, 0.1)", "rgba(196, 28, 28, 1)"]);
  const [maxNoOfCompetitors, setMaxNoOfCompetitors] = useState<number>(0);
  const [
    totalNoOfCompetitors,
    setTotalNoOfCompetitors,
    totalNoOfCompetitorsRef,
  ] = useState<number>(0);

  useEffect(() => {
    setLoadingState({ isLoading: true, error: {} });

    setTotalNoOfCompetitors(0);
    GetMapData()
      .then((res: FeatureCollection) => {
        let _maxNoOfCompetitors = 0;
        res?.features?.map((f) => {
          if (f && f.properties && f.properties.users) {
            const countryNoOfCompetitors = f.properties?.users?.length;
            _maxNoOfCompetitors = Math.max(
              _maxNoOfCompetitors,
              countryNoOfCompetitors
            );
            setTotalNoOfCompetitors(
              totalNoOfCompetitorsRef.current + countryNoOfCompetitors
            );
          }
        });
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
      <Typography level="h2">
        Users {totalNoOfCompetitors !== 0 && <>({totalNoOfCompetitors})</>}
      </Typography>
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
            data-tooltip-offset={20}
            style={{ marginTop: "-100px" }}
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
                              <Stack spacing={1} direction="column">
                                <div
                                  style={{
                                    display: "flex",
                                    alignItems: "center",
                                    borderBottom: "1px solid white",
                                  }}
                                >
                                  <span
                                    className={`fi fi-${geo.properties.countryIso2.toLowerCase()}`}
                                  />
                                  &nbsp;&nbsp;
                                  <Typography sx={{ color: "white" }}>
                                    {geo.properties.name}
                                  </Typography>
                                </div>
                                {geo &&
                                  geo.properties &&
                                  geo.properties.users &&
                                  geo.properties.users.map((user: any) => (
                                    <Typography
                                      sx={{
                                        color: "white",
                                        display: "flex",
                                        alignItems: "center",
                                      }}
                                    >
                                      <b>{user.username}</b>&nbsp;(
                                      <Stack spacing={0.5} direction="row">
                                        <div>{user.score}</div>
                                        <EmojiEvents />
                                      </Stack>
                                      )
                                    </Typography>
                                  ))}
                              </Stack>
                            );
                          }}
                          onMouseLeave={() => {
                            setTooltipContent("");
                          }}
                          className="geo-no-outline"
                          style={{
                            default: {
                              fill:
                                geo.properties.users.length === 0
                                  ? "#F5F4F6"
                                  : defaultColorScale(
                                      geo.properties.users.length /
                                        maxNoOfCompetitors
                                    ),
                            },
                            hover: {
                              fill:
                                geo.properties.users.length === 0
                                  ? "#F5F4F6"
                                  : hoverColorScale(
                                      geo.properties.users.length /
                                        maxNoOfCompetitors
                                    ),
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
          <Tooltip id="my-tooltip" noArrow children={tooltipContent}></Tooltip>
        </>
      )}
    </Stack>
  );
};

export default Users;
