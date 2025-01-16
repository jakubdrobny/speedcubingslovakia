import { Link, useNavigate } from "react-router-dom";
import { NavItemSublistItem } from "../../Types";

const NavItem: React.FC<{
  Title: React.ReactNode;
  to?: string;
  TitleIcon: React.ReactNode;
  sublistItems: NavItemSublistItem[];
  onClick: () => void;
}> = ({ Title, to, TitleIcon, sublistItems, onClick }) => {
  const navigate = useNavigate();

  return (
    <div className="relative group h-full flex items-center">
      <Link
        to={to ? to : {}}
        onClick={to ? onClick : () => { }}
        className="flex items-center gap-2 no-underline text-black hover:bg-gray-100 mx-4"
      >
        {TitleIcon}
        {Title}
      </Link>
      <div className="absolute left-0 top-full w-full h-2"></div>
      {sublistItems && sublistItems.length > 0 && (
        <div
          className="absolute hidden left-1/2 transform -translate-x-1/2 mt-2 group-hover:flex group-hover:flex-col bg-gray-100 gap-2 rounded-md p-4 z-10"
          style={{ top: "100%" }}
        >
          <div className="absolute -top-1 left-1/2 transform -translate-x-1/2 w-4 h-4 bg-gray-100 rotate-45 border-l border-t border-gray-200"></div>
          {sublistItems.map((entry: NavItemSublistItem) => (
            <a
              onClick={() => {
                if (entry.to) {
                  onClick();
                  navigate(entry.to);
                } else if (entry.onClick !== undefined) {
                  entry.onClick();
                }
              }}
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
