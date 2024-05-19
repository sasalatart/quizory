import { createContext, useCallback, useEffect, useState, type PropsWithChildren } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { type Session } from '@supabase/supabase-js';
import { supabase } from '@/supabase';
import { ROUTES } from '@/router';

type SessionContextValue = {
  session: Session | undefined;
  handleLogOut: () => unknown;
  isLoggingOut: boolean;
};

export const SessionContext = createContext<SessionContextValue>({
  session: undefined,
  handleLogOut: () => undefined,
  isLoggingOut: false,
});

export function SessionProvider({ children }: PropsWithChildren): JSX.Element | null {
  const [session, setSession] = useState<Session | undefined>();
  const [isGettingSession, setIsGettingSession] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogOut = useCallback(async () => {
    setIsLoggingOut(true);
    await supabase.auth.signOut();
    setIsLoggingOut(false);
  }, []);

  useEffect(() => {
    void supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session ?? undefined);
      setIsGettingSession(false);
    });
  }, [navigate]);

  useEffect(() => {
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session ?? undefined);

      if (!session) {
        navigate(ROUTES.login);
      } else if (location.pathname === ROUTES.login) {
        navigate(ROUTES.questions.current);
      }
    });

    return () => {
      subscription.unsubscribe();
    };
  }, [location.pathname, navigate]);

  if (isGettingSession) {
    return null;
  }

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
