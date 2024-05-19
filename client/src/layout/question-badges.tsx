import capitalize from 'lodash/capitalize';
import { Difficulty } from '@/generated/api';

const READABLE_DIFFICULTIES: Record<Difficulty, string> = {
  DifficultyNoviceHistorian: 'Novice Historian',
  DifficultyAvidHistorian: 'Avid Historian',
  DifficultyHistoryScholar: 'History Scholar',
};

interface Props {
  topic: string;
  difficulty: Difficulty;
}

export function QuestionBadges({ topic, difficulty }: Props): JSX.Element {
  return (
    <div className="flex flex-col sm:flex-row">
      <div className="badge badge-primary badge-outline">
        Topic: {topic.split(' ').map(capitalize).join(' ')}
      </div>
      <div className="badge badge-secondary badge-outline mt-2 sm:mt-0 sm:ml-2">
        Difficulty: {READABLE_DIFFICULTIES[difficulty]}
      </div>
    </div>
  );
}
