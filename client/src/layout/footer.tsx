import { GITHUB_USER_LINK, GITHUB_REPO_LINK } from '@/config';

export function Footer(): JSX.Element {
  return (
    <footer className="footer footer-center p-4 bg-base-300 text-base-content opacity-85">
      <aside>
        <p>
          <a href={GITHUB_USER_LINK} target="_blank" className="link">
            Sebastián Salata R-T
          </a>
          {' | '}
          See code on{' '}
          <a href={GITHUB_REPO_LINK} target="_blank" className="link">
            GitHub
          </a>
        </p>
      </aside>
    </footer>
  );
}
