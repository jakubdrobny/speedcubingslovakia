import {
  MapContainer,
  TileLayer,
  Tooltip,
  Marker,
  useMapEvents,
  Circle,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { useState } from "react";
import { Stack, Button, ButtonGroup, Chip, Input, Typography } from "@mui/joy";
import { MAX_RADIUS, MIN_RADIUS } from "../../constants";

type Marker = {
  lat: number;
  long: number;
  radius: number;
  new: boolean;
};

const SubscriptionMap = () => {
  const [markers, setMarkers] = useState<Marker[]>([]);

  const MapClickHandler = () => {
    useMapEvents({
      click: (e) => {
        const newMarker = {
          lat: e.latlng.lat,
          long: e.latlng.lng,
          radius: 50,
          new: true,
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

  const stopPropagation = (e) => {
    e.stopPropagation();
  };

  const handleRadiusChange = (e, idx: number) => {
    let newRadius = parseInt(e.target.value || "1");
    if (newRadius < MIN_RADIUS) newRadius = MIN_RADIUS;
    if (newRadius > MAX_RADIUS) newRadius = MAX_RADIUS;

    setMarkers((pm) =>
      pm.map((m, i) => (i !== idx ? m : { ...m, radius: newRadius })),
    );
  };

  const handleMarkerSave = (idx: number) => {
    console.log("saved " + idx);
  };

  return (
    <Stack sx={{ height: "512px" }} spacing={1}>
      <Chip sx={{ fontSize: 12, fontStyle: "italic" }}>
        Note: the circles displayed might not look accurate for large radiuses,
        but the calculations when sending announcements will be done correctly.
      </Chip>
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
          <div key={markerIdx}>
            <Marker position={[marker.lat, marker.long]}>
              <Tooltip
                className="m-0 p-0"
                direction="top"
                offset={[0, -10]}
                permanent
              >
                <div
                  onClick={stopPropagation}
                  onMouseDown={stopPropagation}
                  onTouchStart={stopPropagation}
                  style={{ pointerEvents: "auto", padding: 5 }}
                >
                  <Typography level="h4">Radius (km):</Typography>
                  <Stack direction="row" spacing={1}>
                    <Input
                      size="sm"
                      type="number"
                      sx={{ width: 100 }}
                      value={marker.radius}
                      onChange={(e) => handleRadiusChange(e, markerIdx)}
                    />
                    <Button
                      color="primary"
                      onClick={() => handleMarkerSave(markerIdx)}
                    >
                      Save!
                    </Button>
                  </Stack>
                </div>
              </Tooltip>
            </Marker>
            <Circle
              center={[marker.lat, marker.long]}
              radius={marker.radius * 1000}
            />
          </div>
        ))}
      </MapContainer>
    </Stack>
  );
};

export default SubscriptionMap;
