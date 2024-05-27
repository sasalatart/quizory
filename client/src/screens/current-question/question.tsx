import { useContext, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext } from '@/providers';
import { CenteredSpinner } from '@/layout';
import { QuestionFormCard } from './form-card';
import { Feedback, QuestionFeedbackCard } from './feedback-card';
import { NoQuestionsLeftCard } from './no-questions-left-card';

export function Question(): JSX.Element {
  const [feedback, setFeedback] = useState<Feedback | undefined>();

  const { questionsApi } = useContext(QueryClientContext);
  const {
    data: question,
    isLoading,
    isRefetching,
    refetch: refetchNextQuestion,
  } = useQuery('current-question', {
    queryFn: () => questionsApi.getNextQuestion(),
    refetchOnWindowFocus: false,
  });

  if (isLoading && !question) {
    return <CenteredSpinner />;
  }

  const shouldShowFeedback = question?.choices.some(({ id }) => id === feedback?.selectedChoiceId);
  if (question && feedback && shouldShowFeedback) {
    return (
      <QuestionFeedbackCard
        question={question}
        feedback={feedback}
        isLoadingNext={isRefetching}
        onNext={refetchNextQuestion}
      />
    );
  }

  if (!question) {
    return <NoQuestionsLeftCard />;
  }

  return <QuestionFormCard question={question} onSubmit={(feedback) => setFeedback(feedback)} />;
}
