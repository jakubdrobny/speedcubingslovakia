import { Link, useNavigate } from "react-router-dom";
import { NavContextType, NavItemSublistItem } from "../../Types";
import { useContext } from "react";
import { NavContext } from "../../context/NavContext";

const NavItem: React.FC<{
  Title: React.ReactNode;
  to?: string;
  listItemType?: "competitions" | "results" | "profile";
  TitleIcon: React.ReactNode;
  sublistItems: NavItemSublistItem[];
  onClick: () => void;
}> = ({ Title, to, listItemType, TitleIcon, sublistItems, onClick }) => {
  const navigate = useNavigate();
  const { toggleSublistOpen, sublistOpen, closeSublists } = useContext(
    NavContext,
  ) as NavContextType;

  return (
    <div className="relative group w-full h-full flex items-center justify-end">
      <Link
        to={to ? to : {}}
        onClick={(e) => {
          if (sublistItems && sublistItems.length > 0) {
            e.preventDefault();
            toggleSublistOpen(listItemType);
          } else if (to) {
            onClick();
            closeSublists();
          }
        }}
        className={`flex items-center gap-2 w-full no-underline text-black hover:bg-gray-100 mx-4 justify-end ${sublistOpen(listItemType) && "bg-gray-100"}`}
      >
        {TitleIcon}
        {Title}
      </Link>
      <div className="absolute left-0 top-full w-full h-2"></div>
      {sublistItems && sublistItems.length > 0 && (
        <div
          className={`absolute left-1/2 transform -translate-x-1/2 mt-2 group-hover:flex group-hover:flex-col ${sublistOpen(listItemType) ? "flex flex-col" : "hidden"} bg-gray-100 gap-2 rounded-md p-4 z-10`}
          style={{ top: "100%" }}
        >
          <div className="absolute -top-1 left-1/2 transform -translate-x-1/2 w-4 h-4 bg-gray-100 rotate-45 border-l border-t border-gray-200"></div>
          {sublistItems.map((entry: NavItemSublistItem, idx: number) => (
            <a
              onClick={() => {
                if (entry.to) {
                  onClick();
                  navigate(entry.to);
                } else if (entry.onClick !== undefined) {
                  entry.onClick();
                }
                closeSublists();
              }}
              key={idx.toString() + entry.title}
              className="flex text-sm cursor-pointer items-center no-underline text-nowrap gap-2 text-gray-800 hover:text-black"
            >
              {entry.icon}
              <span>{entry.title}</span>
            </a>
          ))}
        </div>
      )}
    </div>
  );
};

export default NavItem;
