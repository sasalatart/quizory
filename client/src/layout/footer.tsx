import { CONFIG } from './config';

export function Footer(): JSX.Element {
  return (
    <footer className="footer footer-center p-4 bg-base-300 text-base-content">
      <aside>
        <p>
          <a href={CONFIG.githubUserLink} target="_blank" className="link">
            Sebastián Salata R-T
          </a>
          {' | '}
          See code on{' '}
          <a href={CONFIG.githubRepoLink} target="_blank" className="link">
            GitHub
          </a>
        </p>
      </aside>
    </footer>
  );
}
