import { AverageInfo, CompetitionContextType, LoadingState } from "../../Types";
import { Card, Chip, Stack, Typography } from "@mui/joy";
import {
  GetAverageInfo,
  GetAverageInfoRecords,
  getError,
  initialAverageInfo,
  isObjectEmpty,
  reformatFMCSolve,
  reformatMultiTime,
  renderResponseError,
} from "../../utils/utils";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "../../context/CompetitionContext";
import LoadingComponent from "../Loading/LoadingComponent";
import ResultsModal from "./ResultsModal";

const AveragePreview: React.FC<{
  showResultsModal: boolean;
  loadingResults: boolean;
}> = ({ showResultsModal, loadingResults }) => {
  const [averageInfo, setAverageInfo] =
    useState<AverageInfo>(initialAverageInfo);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [party, setParty] = useState<boolean>(false);
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    error: {},
  });
  const { currentResults } = useContext(
    CompetitionContext
  ) as CompetitionContextType;
  const ismbld = currentResults.iconcode === "333mbf";
  const isfmc = currentResults.iconcode === "333fm";
  const isBo1 = currentResults.format === "bo1";

  useEffect(() => {
    if (currentResults.format && !loadingResults) {
      setLoadingState({ isLoading: true, error: {} });
      GetAverageInfo(currentResults)
        .then((res) => {
          setAverageInfo(res);
          if (res.finishedCompeting && showResultsModal)
            return GetAverageInfoRecords(currentResults, res);
        })
        .then((res) => {
          if (!res) {
            setLoadingState({ isLoading: false, error: {} });
            return;
          }
          setAverageInfo(res);
          if (res.finishedCompeting && showResultsModal) {
            setIsModalOpen(true);
            if (hasAnyRecord(res)) setParty(true);
          }
          setLoadingState({ isLoading: false, error: {} });
        })
        .catch((err) => {
          setLoadingState({ isLoading: false, error: getError(err) });
        });
    }
  }, [currentResults.format, loadingResults]);

  const hasAnyRecord = (averageInfo: AverageInfo): boolean => {
    return !(!averageInfo.singleRecord && !averageInfo.averageRecord);
  };

  return (
    <Stack
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        overflow: "auto",
      }}
      component={Card}
    >
      <ResultsModal
        isModalOpen={isModalOpen}
        setIsModalOpen={setIsModalOpen}
        averageInfo={averageInfo}
        isfmc={isfmc}
        ismbld={ismbld}
        isbo1={isBo1}
        party={party}
        setParty={setParty}
      />
      <h3 style={{ padding: 0, margin: 0, borderBottom: "1px solid #CDD7E1" }}>
        Average preview:
      </h3>
      {!isObjectEmpty(loadingState.error) ? (
        renderResponseError(loadingState.error)
      ) : (
        <>
          <Stack
            spacing={2}
            direction="row"
            display="flex"
            justifyContent="center"
          >
            <div
              style={
                averageInfo.showPossibleAverage
                  ? {
                      display: "flex",
                      justifyContent: "flex-end",
                    }
                  : {}
              }
            >
              <Typography
                sx={{ display: "flex", alignItems: "center" }}
                component="div"
              >
                <b>Single:</b>&nbsp;
                <Chip variant="soft" color="primary">
                  {ismbld
                    ? reformatMultiTime(averageInfo.single)
                    : isfmc
                    ? reformatFMCSolve(averageInfo.single)
                    : averageInfo.single}
                </Chip>
              </Typography>
            </div>
            {averageInfo.finishedCompeting && !ismbld && !isBo1 && (
              <div>
                <Typography
                  sx={{ display: "flex", alignItems: "center" }}
                  component="div"
                >
                  <b>Average:</b>&nbsp;
                  <Chip variant="soft" color="primary">
                    {averageInfo.average}
                  </Chip>
                </Typography>
              </div>
            )}
          </Stack>
          <div
            style={{
              display: "flex",
              justifyContent: "center",
            }}
          >
            <Typography
              sx={{ display: "flex", alignItems: "center" }}
              component="div"
            >
              <b>{ismbld ? "Attempts" : isfmc ? "Solves" : "Times"}</b>:&nbsp;
              <Chip variant="soft" color="warning">
                <Stack spacing={1} direction="row">
                  {averageInfo.times.map((solveTime, idx) => (
                    <div key={idx.toString() + "#" + solveTime}>
                      {ismbld
                        ? reformatMultiTime(solveTime)
                        : isfmc
                        ? reformatFMCSolve(solveTime)
                        : solveTime}
                    </div>
                  ))}
                </Stack>
              </Chip>
            </Typography>
          </div>
          {averageInfo.showPossibleAverage && (
            <Stack
              spacing={2}
              direction="row"
              display="flex"
              justifyContent="center"
            >
              <div style={{ display: "flex", justifyContent: "flex-end" }}>
                <Typography
                  sx={{ display: "flex", alignItems: "center" }}
                  component="div"
                >
                  <b>BPA:</b>&nbsp;
                  <Chip variant="soft" color="success">
                    {averageInfo.bpa}
                  </Chip>
                </Typography>
              </div>
              <div>
                <Typography
                  sx={{ display: "flex", alignItems: "center" }}
                  component="div"
                >
                  <b>WPA:</b>&nbsp;
                  <Chip variant="soft" color="danger">
                    {averageInfo.wpa}
                  </Chip>
                </Typography>
              </div>
            </Stack>
          )}
        </>
      )}
    </Stack>
  );
};

export default AveragePreview;
