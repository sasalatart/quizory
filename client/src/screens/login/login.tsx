import { PropsWithChildren, useContext } from 'react';
import { Auth } from '@supabase/auth-ui-react';
import { ThemeSupa } from '@supabase/auth-ui-shared';
import { SessionContext } from '@/providers';
import { supabase } from '../../supabase';

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
