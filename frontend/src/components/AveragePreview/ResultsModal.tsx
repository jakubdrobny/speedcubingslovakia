import {
  DialogContent,
  DialogTitle,
  Divider,
  Modal,
  ModalClose,
  ModalDialog,
} from "@mui/joy";
import { reformatFMCSolve, reformatMultiTime } from "../../utils";

import { AverageInfo } from "../../Types";
import { EmojiEvents } from "@mui/icons-material";

const ResultsModal: React.FC<{
  isModalOpen: boolean;
  setIsModalOpen: (newIsModalOpen: boolean) => void;
  averageInfo: AverageInfo;
  isfmc: boolean;
  ismbld: boolean;
  isbo1: boolean;
}> = ({ isModalOpen, setIsModalOpen, averageInfo, isfmc, ismbld, isbo1 }) => {
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

  return (
    <Modal open={isModalOpen} onClose={() => setIsModalOpen(false)}>
      <ModalDialog
        color="success"
        layout="center"
        size="md"
        variant="soft"
        role="alertdialog"
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
              </span>
            </b>
            single
            {!isbo1 && !ismbld && (
              <>
                {" "}
                and <b>{average}</b>{" "}
                <span style={{ color: averageInfo.averageRecordColor }}>
                  {averageInfo.averageRecord}
                </span>{" "}
                average
              </>
            )}
            .
          </div>
        </DialogContent>
      </ModalDialog>
    </Modal>
  );
};

export default ResultsModal;
