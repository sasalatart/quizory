import { createContext, useCallback, useEffect, useState, type PropsWithChildren } from 'react';
import { supabase } from '../supabase';
import { type Session } from '@supabase/supabase-js';

type SessionContextValue = {
  session: Session | undefined;
  handleLogOut: () => unknown;
  isLoggingOut: boolean;
};

const SessionContext = createContext<SessionContextValue>({
  session: undefined,
  handleLogOut: () => undefined,
  isLoggingOut: false,
});

function SessionProvider({ children }: PropsWithChildren): JSX.Element {
  const [session, setSession] = useState<Session | undefined>();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  const handleLogOut = useCallback(async () => {
    setIsLoggingOut(true);
    await supabase.auth.signOut();
    setIsLoggingOut(false);
  }, []);

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
        handleLogOut,
        isLoggingOut,
      }}
    >
      {children}
    </SessionContext.Provider>
  );
}

export { SessionContext, SessionProvider };
