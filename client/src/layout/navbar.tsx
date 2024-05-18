import { useContext } from 'react';
import { SessionContext } from '@/providers';

export function Navbar(): JSX.Element {
  const { handleLogOut, isLoggingOut } = useContext(SessionContext);

  return (
    <div className="navbar bg-base-100">
      <div className="navbar-start">
        <div className="dropdown">
          <div tabIndex={0} role="button" className="btn btn-ghost lg:hidden">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 6h16M4 12h8m-8 6h16"
              />
            </svg>
          </div>
          <ul
            tabIndex={0}
            className="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52"
          >
            <li>
              <a>Current Question</a>
            </li>
            <li>
              <a>Past Questions</a>
            </li>
          </ul>
        </div>
        <a className="btn btn-ghost text-xl">Quizory</a>
      </div>
      <div className="navbar-center hidden lg:flex">
        <ul className="menu menu-horizontal px-1">
          <li>
            <a>Current Question</a>
          </li>
          <li>
            <a>Past Questions</a>
          </li>
        </ul>
      </div>
      <div className="navbar-end">
        <button onClick={handleLogOut} className="btn" disabled={isLoggingOut}>
          {isLoggingOut ? <span className="loading loading-spinner"></span> : null}
          Log Out
        </button>
      </div>
    </div>
  );
}
