import { Credits } from './credits';

export function Footer(): JSX.Element {
  return (
    <footer className="footer footer-center p-4 bg-base-300 text-base-content opacity-85">
      <Credits />
    </footer>
  );
}
