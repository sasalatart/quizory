import { useContext } from 'react';
import { Link } from 'react-router-dom';
import { SessionContext } from '@/providers';
import { HamburgerIcon, NapoleonicHatIcon } from '@/icons';
import { InlineSpinner } from './spinner';

export function Navbar(): JSX.Element {
  const { handleLogOut, isLoggingOut } = useContext(SessionContext);

  return (
    <div className="navbar bg-base-100">
      <div className="navbar-start">
        <div className="dropdown">
          <div tabIndex={0} role="button" className="btn btn-ghost lg:hidden">
            <HamburgerIcon />
          </div>
          <ul
            tabIndex={0}
            className="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52"
          >
            <li>
              <Link to="/questions/next">Current Question</Link>
            </li>
            <li>
              <Link to="/questions/log">Past Questions</Link>
            </li>
          </ul>
        </div>
        <NapoleonicHatIcon />
        <Link to="/questions/next" className="btn btn-ghost text-xl">
          Quizory
        </Link>
      </div>
      <div className="navbar-center hidden lg:flex">
        <ul className="menu menu-horizontal px-1">
          <li>
            <Link to="/questions/next">Current Question</Link>
          </li>
          <li>
            <Link to="/questions/log">Past Questions</Link>
          </li>
        </ul>
      </div>
      <div className="navbar-end">
        <button onClick={handleLogOut} className="btn" disabled={isLoggingOut}>
          {isLoggingOut ? <InlineSpinner /> : null}
          Log Out
        </button>
      </div>
    </div>
  );
}
