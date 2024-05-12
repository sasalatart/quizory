import { Auth } from '@supabase/auth-ui-react';
import { ThemeSupa } from '@supabase/auth-ui-shared';

import { supabase } from '../supabase';

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
