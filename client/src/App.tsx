import { useContext } from 'react';
import { SessionContext, SessionProvider } from './auth/session-provider';
import { Login } from './auth';

export default function App() {
  return (
    <SessionProvider>
      <Landing />
    </SessionProvider>
  );
}

function Landing(): JSX.Element {
  const { session, logout } = useContext(SessionContext);

  if (!session) {
    return <Login />;
  }

  return (
    <div>
      <button onClick={logout} className="btn btn-outline btn-primary">
        Sign Out
      </button>
    </div>
  );
}
