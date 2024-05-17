import { useContext } from 'react';
import { useQuery } from 'react-query';
import capitalize from 'lodash/capitalize';
import { Difficulty } from '@/generated/api';
import { QueryClientContext } from '@/providers';
import { HintButton } from './hint-button';

const READABLE_DIFFICULTIES: Record<Difficulty, string> = {
  DifficultyNoviceHistorian: 'Novice Historian',
  DifficultyAvidHistorian: 'Avid Historian',
  DifficultyHistoryScholar: 'History Scholar',
};

export function Question(): JSX.Element {
  const apiClient = useContext(QueryClientContext);

  const { data: question, isLoading } = useQuery('current-question', {
    queryFn: () => apiClient.getNextQuestion(),
  });

  if (isLoading || !question) {
    return <span className="loading loading-spinner loading-xs"></span>;
  }

  return (
    <div className="flex flex-grow justify-center align-middle">
      <div className="card bg-neutral shadow-xl m-4">
        <div className="card-body">
          <h2 className="card-title">{question.question}</h2>
          <div className="flex flex-col sm:flex-row">
            <div className="badge badge-primary badge-outline">
              Topic: {question.topic.split(' ').map(capitalize).join(' ')}
            </div>
            <div className="badge badge-secondary badge-outline mt-2 sm:mt-0 sm:ml-2">
              Difficulty: {READABLE_DIFFICULTIES[question.difficulty]}
            </div>
          </div>

          {question.choices.map(({ id, choice }) => (
            <div key={id} className="form-control">
              <label className="cursor-pointer label">
                <span className="label-text">{choice}</span>
                <input type="radio" name="choice" value={id} className="radio" />
              </label>
            </div>
          ))}

          <div className="card-actions justify-end">
            <button className="btn btn-primary">Submit</button>
            <HintButton hint={question.hint} />
          </div>
        </div>
      </div>
    </div>
  );
}
