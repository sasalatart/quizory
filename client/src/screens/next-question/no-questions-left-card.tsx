export function NoQuestionsLeftCard(): JSX.Element {
  return (
    <div className="card bg-neutral shadow-xl">
      <div className="card-body">
        <h2 className="card-title">No questions are left</h2>
        <p>You have answered all the available questions so far. More will be generated soon.</p>
      </div>
    </div>
  );
}
