import { Auth } from '@supabase/auth-ui-react';
import { ThemeSupa } from '@supabase/auth-ui-shared';

import { supabase } from '../supabase';
import { PropsWithChildren, useContext } from 'react';
import { SessionContext } from '../providers/session-provider';

export function Login(): JSX.Element {
  return (
    <Auth
      supabaseClient={supabase}
      providers={['google']}
      appearance={{ theme: ThemeSupa }}
      onlyThirdPartyProviders
    />
  );
}

export function EnsureLoggedIn({ children }: PropsWithChildren): JSX.Element {
  const { session } = useContext(SessionContext);
  return session ? <>{children}</> : <Login />;
}
