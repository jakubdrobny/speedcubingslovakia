import {
  MapContainer,
  TileLayer,
  Tooltip,
  useMapEvents,
  Marker,
  Circle,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { useEffect } from "react";
import useState from "react-usestateref";
import { Stack, Button, Chip, Input, Typography, IconButton } from "@mui/joy";
import { MAX_RADIUS, MIN_RADIUS } from "../../constants";
import { LoadingState, MarkerType, ResponseError } from "../../Types";
import { AxiosError } from "axios";
import {
  DeleteMarker,
  getError,
  GetMarkers,
  initialLoadingState,
  isObjectEmpty,
  renderResponseError,
  saveCurrentLocation,
  SaveMarker,
} from "../../utils/utils";
import { Close } from "@mui/icons-material";
import L from "leaflet";

import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";
import { Link } from "react-router-dom";
import PleaseLoginButton from "./PleaseLoginButton";

const DefaultIcon = L.icon({
  iconUrl: icon,
  shadowUrl: iconShadow,
  iconAnchor: [12.5, 41],
  tooltipAnchor: [0, -41],
});

L.Marker.prototype.options.icon = DefaultIcon;

const SubscriptionMap = () => {
  const [loadingState, setLoadingState] =
    useState<LoadingState>(initialLoadingState);
  const [markers, setMarkers] = useState<MarkerType[]>([]);

  useEffect(() => {
    setLoadingState({
      isLoading: true,
      error: {},
    });

    GetMarkers()
      .then((res: MarkerType[]) => {
        setMarkers(res);
        setLoadingState({
          isLoading: false,
          error: {},
        });
      })
      .catch((err: AxiosError) => {
        setLoadingState({
          isLoading: false,
          error: customGetError(err),
        });
      });
  }, []);

  const customGetError = (err: AxiosError): ResponseError => {
    if (err.response?.status === 401) {
      return {
        element: <PleaseLoginButton />,
      };
    }
    return getError(err);
  };

  const MapClickHandler = () => {
    useMapEvents({
      click: (e) => {
        const newMarker = {
          id: 0,
          lat: e.latlng.lat,
          long: e.latlng.lng,
          radius: 50,
          new: true,
          open: true,
        };
        if (markers.length > 0 && markers[markers.length - 1].new) {
          newMarker.radius = markers[markers.length - 1].radius;
          setMarkers((pm) => pm.map((m) => (!m.new ? m : newMarker)));
        } else {
          setMarkers((pm) => [...pm, newMarker]);
        }
      },
    });
    return null;
  };

  const stopPropagation = (e: React.SyntheticEvent) => {
    e.stopPropagation();
  };

  const handleRadiusChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    idx: number,
  ) => {
    let newRadius = parseInt(e.target.value || "0");
    if (newRadius < MIN_RADIUS) newRadius = MIN_RADIUS;
    if (newRadius > MAX_RADIUS) newRadius = MAX_RADIUS;

    setMarkers((pm) =>
      pm.map((m, i) => (i !== idx ? m : { ...m, radius: newRadius })),
    );
  };

  const handleMarkerSave = (idx: number) => {
    setLoadingState({
      isLoading: true,
      error: {},
    });

    SaveMarker(markers[idx])
      .then(() => {
        setMarkers((p: MarkerType[]) =>
          p.map((m, i) => (i !== idx ? m : { ...m, new: false, open: false })),
        );
        setLoadingState({
          isLoading: false,
          error: {},
        });
      })
      .catch((err: AxiosError) => {
        setLoadingState({ isLoading: false, error: customGetError(err) });
      });
  };

  const handleMarkerDelete = (idx: number) => {
    if (markers[idx].new) {
      handleMarkerClose(idx);
      return;
    }

    setLoadingState({
      isLoading: true,
      error: {},
    });

    DeleteMarker(markers[idx])
      .then(() => {
        setMarkers((p: MarkerType[]) => p.filter((_, i) => i !== idx));
        setLoadingState({
          isLoading: false,
          error: {},
        });
      })
      .catch((err: AxiosError) => {
        setLoadingState({ isLoading: false, error: customGetError(err) });
      });
  };

  const handleMarkerOpenToggle = (idx: number) => {
    if (markers[idx].new) return;
    setMarkers((p: MarkerType[]) =>
      p.map((m, i) => ({
        ...m,
        open: i !== idx ? m.new : !m.open,
      })),
    );
  };

  const handleMarkerClose = (idx: number) => {
    if (markers[idx].new) {
      setMarkers((p: MarkerType[]) => p.filter((_, i) => i !== idx));
    } else {
      setMarkers((p: MarkerType[]) =>
        p.map((m, i) => (i !== idx ? m : { ...m, open: false })),
      );
    }
  };

  const formatRadius = (radius: number): string => {
    let newRadius: string = radius.toString();
    while (newRadius.length > 1 && newRadius[0] === "0")
      newRadius = newRadius.substring(1);
    return newRadius;
  };

  return (
    <Stack spacing={2}>
      {!isObjectEmpty(loadingState.error) &&
        renderResponseError(loadingState.error)}
      <Stack spacing={0}>
        <Chip
          sx={{
            fontSize: 12,
            maxWidth: "100%",
            fontStyle: "italic",
            borderRadius: "16px 16px 0 0",
            px: 2,
            "& .MuiChip-label": { overflow: "auto" },
          }}
        >
          <b>Note:</b> the circles displayed might not look accurate for *VERY*
          large radiuses, but the calculations when sending announcements will
          be done correctly.
        </Chip>
        <div style={{ height: "512px" }}>
          <MapContainer
            center={[0, 0]}
            zoom={2}
            style={{ height: "100%", minHeight: "100%" }}
          >
            <TileLayer
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              attribution="&copy; <a href='https://www.openstreetmap.org/copyright'>OpenStreetMap</a> contributors"
            />
            <MapClickHandler />
            {markers.map((marker, markerIdx) => (
              <div key={markerIdx + "" + marker.id}>
                <Marker
                  position={[marker.lat, marker.long]}
                  eventHandlers={{
                    click: () => handleMarkerOpenToggle(markerIdx),
                  }}
                >
                  {marker.open && (
                    <Tooltip className="m-0 p-0" direction="top" permanent>
                      <div
                        onClick={stopPropagation}
                        onMouseDown={stopPropagation}
                        onTouchStart={stopPropagation}
                        style={{ pointerEvents: "auto", padding: 5 }}
                      >
                        <Stack
                          direction="row"
                          sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "center",
                          }}
                        >
                          <Typography level="h4">Radius (km):</Typography>
                          <IconButton
                            onClick={() => handleMarkerClose(markerIdx)}
                          >
                            <Close fontSize="small" />
                          </IconButton>
                        </Stack>
                        <Stack direction="row" spacing={1}>
                          <Input
                            size="sm"
                            type="number"
                            sx={{ width: 100 }}
                            value={formatRadius(marker.radius)}
                            onChange={(e) => handleRadiusChange(e, markerIdx)}
                          />
                          <Button
                            color="primary"
                            onClick={() => handleMarkerSave(markerIdx)}
                          >
                            Save!
                          </Button>
                          {!marker.new && (
                            <Button
                              color="danger"
                              onClick={() => handleMarkerDelete(markerIdx)}
                            >
                              Delete!
                            </Button>
                          )}
                        </Stack>
                      </div>
                    </Tooltip>
                  )}
                </Marker>
                <Circle
                  center={[marker.lat, marker.long]}
                  radius={marker.radius * 1000}
                />
              </div>
            ))}
          </MapContainer>
        </div>
      </Stack>
    </Stack>
  );
};

export default SubscriptionMap;
