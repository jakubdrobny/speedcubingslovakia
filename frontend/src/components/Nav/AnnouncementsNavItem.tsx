import { IconSpeakerphone } from "@tabler/icons-react";
import NavItem from "./NavItem";
import { useContext } from "react";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";

const AnnouncementsNavItem: React.FC<{
  newAnnouncements: number;
  isAuthenticated: boolean;
}> = ({ newAnnouncements, isAuthenticated }) => {
  const { closeNav } = useContext(NavContext) as NavContextType;

  return (
    <NavItem
      Title={
        <div className="relative inline-block">
          <div>Announcements</div>

          {isAuthenticated && newAnnouncements > 0 && (
            <div className="absolute -top-3 -right-3 w-5 h-5 bg-red-100 text-red-700 text-xs rounded-full flex items-center justify-center border border-red-200">
              {newAnnouncements}
            </div>
          )}
        </div>
      }
      onClick={closeNav}
      TitleIcon={<IconSpeakerphone />}
      to="/announcements"
      sublistItems={[]}
    />
  );
};

export default AnnouncementsNavItem;
