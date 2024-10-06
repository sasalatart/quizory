import { GITHUB_REPO_LINK, GITHUB_USER_LINK, GITHUB_USER_NAME } from '@/config';

export function Credits(): JSX.Element {
  return (
    <aside>
      <p>
        <a href={GITHUB_USER_LINK} target="_blank" className="link">
          {GITHUB_USER_NAME}
        </a>
        {' | '}
        See on{' '}
        <a href={GITHUB_REPO_LINK} target="_blank" className="link">
          GitHub
        </a>
      </p>
    </aside>
  );
}
