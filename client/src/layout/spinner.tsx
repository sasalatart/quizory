export function CenteredSpinner(): JSX.Element {
  return (
    <div className="flex justify-center items-center h-full">
      <InlineSpinner />
    </div>
  );
}

export function InlineSpinner(): JSX.Element {
  return <span className="loading loading-spinner" />;
}
