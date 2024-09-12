import {
  ComposableMap,
  Geographies,
  Geography,
  Graticule,
  Sphere,
  ZoomableGroup,
} from "react-simple-maps";
import { LoadingState, SearchUser } from "../../Types";
import { Stack, Typography } from "@mui/joy";
import { getError, getUsers, initialLoadingState } from "../../utils/utils";
import { useEffect, useState } from "react";

import { Tooltip } from "react-tooltip";
import { csv } from "d3-fetch";
import { scaleLinear } from "d3-scale";

const geoUrl = "https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json";
// const geoUrl = "/features.json";

const colorScale = scaleLinear<string, string>()
  .domain([0.29, 0.68])
  .range(["#ffedea", "#ff5233"]);

const Users = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [users, setUsers] = useState<SearchUser[]>([]);
  const [tooltipContent, setTooltipContent] = useState("");

  const searchForUsers = () => {
    setLoadingState({ isLoading: true, error: {} });

    getUsers(searchQuery)
      .then((res: SearchUser[]) => {
        setUsers(res);
        setLoadingState({ isLoading: false, error: {} });
      })
      .catch((err) => {
        setLoadingState({ isLoading: false, error: getError(err) });
      });
  };

  const [data, setData] = useState([]);

  useEffect(() => {
    csv(`${process.env.PUBLIC_URL}/vulnerability.csv`).then((data: any) => {
      setData(data);
    });
  }, []);

  return (
    <Stack sx={{ margin: "1em" }}>
      <Typography level="h2">Users</Typography>
      <div data-tooltip-id="my-tooltip" data-tooltip-content="">
        <ComposableMap
          projectionConfig={{
            scale: 147,
          }}
        >
          <ZoomableGroup>
            <Sphere
              id="sphere-id"
              fill="transparent"
              stroke="#E4E5E6"
              strokeWidth={1}
            />
            <Graticule stroke="#E4E5E6" strokeWidth={1} />
            {data.length > 0 && (
              <Geographies geography={geoUrl}>
                {({ geographies }) =>
                  geographies.map((geo) => {
                    const d = data.find((s: any) => s.ISO3 === geo.id);
                    return (
                      <Geography
                        key={geo.rsmKey}
                        geography={geo}
                        fill={"#0B6BCB"}
                        stroke="white"
                        strokeWidth={1}
                      />
                    );
                  })
                }
              </Geographies>
            )}
          </ZoomableGroup>
        </ComposableMap>
      </div>
      <Tooltip id="my-tooltip">{tooltipContent}</Tooltip>
    </Stack>
  );
};

export default Users;
