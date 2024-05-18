interface Props {
  height?: number;
  width?: number;
}

export function NapoleonicHatIcon({ height = 32, width = 80 }: Props): JSX.Element {
  return (
    <svg height={height} width={width} viewBox="0 30 240 140" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M30 130 Q120 0, 210 130 Q180 140, 60 140 Q30 130, 30 130 Z"
        fill="#2e2e2e"
        stroke="#ffffff"
        strokeWidth="2"
      />
      <path d="M50 130 Q120 40, 190 130" fill="none" stroke="#FFD700" strokeWidth="2" />
      <circle cx="140" cy="90" r="12" fill="#FF0000" />
      <circle cx="140" cy="90" r="7" fill="#ffffff" />
      <circle cx="140" cy="90" r="4" fill="#000000" />
    </svg>
  );
}

export function HamburgerIcon(): JSX.Element {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className="h-5 w-5"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M4 6h16M4 12h8m-8 6h16"
      />
    </svg>
  );
}
