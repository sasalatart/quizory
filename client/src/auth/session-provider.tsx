import {
  createContext,
  useEffect,
  useState,
  type PropsWithChildren,
} from 'react';
import { supabase } from '../supabase';
import { type Session } from '@supabase/supabase-js';

type SessionContextValue = {
  session: Session | undefined;
  logout: () => unknown;
};

const SessionContext = createContext<SessionContextValue>({
  session: undefined,
  logout: () => undefined,
});

function SessionProvider({ children }: PropsWithChildren): JSX.Element {
  const [session, setSession] = useState<Session | undefined>();

  useEffect(() => {
    void supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session ?? undefined);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session ?? undefined);
    });

    return () => {
      subscription.unsubscribe();
    };
  }, []);

  return (
    <SessionContext.Provider
      value={{
        session,
        logout: () => supabase.auth.signOut(),
      }}
    >
      {children}
    </SessionContext.Provider>
  );
}

export { SessionContext, SessionProvider };
