import {
  DialogContent,
  DialogTitle,
  Divider,
  Modal,
  ModalClose,
  ModalDialog,
} from "@mui/joy";
import { reformatFMCSolve, reformatMultiTime } from "../../utils/utils";

import { AverageInfo } from "../../Types";
import { EmojiEvents } from "@mui/icons-material";
import SizedConfetti from "./SizedConfetti";
import { useContainerDimensions } from "../../utils/useContainerDimensions";
import { useRef } from "react";

const ResultsModal: React.FC<{
  isModalOpen: boolean;
  setIsModalOpen: (newIsModalOpen: boolean) => void;
  averageInfo: AverageInfo;
  isfmc: boolean;
  ismbld: boolean;
  isbo1: boolean;
  party: boolean;
  setParty: (newParty: boolean) => void;
}> = ({
  isModalOpen,
  setIsModalOpen,
  averageInfo,
  isfmc,
  ismbld,
  isbo1,
  party,
  setParty,
}) => {
  const single = ismbld
    ? reformatMultiTime(averageInfo.single)
    : isfmc
    ? reformatFMCSolve(averageInfo.single)
    : averageInfo.single;
  const average = ismbld
    ? reformatMultiTime(averageInfo.average)
    : isfmc
    ? reformatFMCSolve(averageInfo.average)
    : averageInfo.average;
  const modalRef = useRef<HTMLDivElement>(null);
  const modalDimensionsRef = useContainerDimensions(modalRef);

  return (
    <div>
      <Modal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        ref={modalRef}
      >
        <>
          <SizedConfetti
            style={{ pointerEvents: "none" }}
            numberOfPieces={party ? 500 : 0}
            recycle={false}
            onConfettiComplete={(confetti: any) => {
              setParty(false);
              confetti.reset();
            }}
          />
          <ModalDialog
            color="success"
            layout="center"
            size="md"
            variant="soft"
            role="alertdialog"
            // ref={modalRef}
          >
            <DialogTitle>
              <EmojiEvents />
              Results
            </DialogTitle>
            <ModalClose />
            <Divider />
            <DialogContent>
              <div>
                You are currently in the <b>{averageInfo.place}</b> place with a{" "}
                <b>
                  {single}{" "}
                  <span style={{ color: averageInfo.singleRecordColor }}>
                    {averageInfo.singleRecord}
                    {averageInfo.singleRecord ? " " : ""}
                  </span>
                </b>
                Single
                {!isbo1 && !ismbld && (
                  <>
                    {" "}
                    and a{" "}
                    <b>
                      {average}{" "}
                      <span style={{ color: averageInfo.averageRecordColor }}>
                        {averageInfo.averageRecord}
                        {averageInfo.averageRecord ? " " : ""}
                      </span>
                    </b>{" "}
                    Average
                  </>
                )}
                .
              </div>
            </DialogContent>
          </ModalDialog>
        </>
      </Modal>
    </div>
  );
};

export default ResultsModal;
