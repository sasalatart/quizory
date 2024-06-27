import { useState } from 'react';
import { CenteredSpinner } from '@/layout';
import { QuestionFormCard } from './form-card';
import { Feedback, QuestionFeedbackCard } from './feedback-card';
import { NoQuestionsLeftCard } from './no-questions-left-card';
import { useCurrentQuestion, useSubmitAnswer } from './hooks';

export function Question(): JSX.Element {
  const [feedback, setFeedback] = useState<Feedback>();

  const {
    question,
    remainingTopics,
    isLoading,
    handleChangeTopic,
    refetch,
  } = useCurrentQuestion();

  const { handleSubmitAnswer } = useSubmitAnswer({
    onSubmit: setFeedback,
  });

  if (isLoading) {
    return <CenteredSpinner />;
  }

  if (!question) {
    return <NoQuestionsLeftCard />;
  }

  if (feedback) {
    return (
      <QuestionFeedbackCard
        question={question}
        feedback={feedback}
        onNext={async () => {
          await refetch();
          setFeedback(undefined);
        }}
      />
    );
  }

  return (
    <QuestionFormCard
      question={question}
      onChangeTopic={handleChangeTopic}
      onSubmit={handleSubmitAnswer}
      remainingTopics={remainingTopics ?? []}
    />
  );
}
