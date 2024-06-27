import { useContext, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext, SessionContext } from '@/providers';
import { RemainingTopic } from '@/generated/api';

/**
 * Returns the next topic to query, based on the following priority order:
 * 1. Current topic if there are remaining questions for it.
 * 2. Last answered topic if there are remaining questions for it.
 * 3. The first remaining topic as a final fallback.
 */
function getTopicToQuery(
  selectedTopic: string | undefined,
  lastAnsweredTopic: string | undefined,
  remainingTopics: RemainingTopic[],
) {
  const hasRemainingQuestionsFor = (forTopic: string) => {
    return remainingTopics.some(({ topic }) => topic === forTopic);
  };

  if (selectedTopic && hasRemainingQuestionsFor(selectedTopic)) {
    return selectedTopic;
  }

  if (lastAnsweredTopic && hasRemainingQuestionsFor(lastAnsweredTopic)) {
    return lastAnsweredTopic;
  }

  return remainingTopics?.[0]?.topic;
}

export function useCurrentQuestion() {
  const { session } = useContext(SessionContext);
  const { answersApi, questionsApi } = useContext(QueryClientContext);
  const [selectedTopic, setSelectedTopic] = useState<string>();

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['current-question', selectedTopic],
    queryFn: async () => {
      // TODO: these two operations could probably be merged into one.
      const [lastAnsweredTopic, remainingTopics] = await Promise.all([
        answersApi
          .getAnswersLog({ userId: session!.user.id, page: 0, pageSize: 1 })
          .then((log) => log[0]?.question.topic),
        questionsApi.getRemainingTopics(),
      ]);

      const topic = getTopicToQuery(selectedTopic, lastAnsweredTopic, remainingTopics);
      const question = await questionsApi.getNextQuestion({ topic });
      return { question, remainingTopics };
    },
    refetchOnWindowFocus: false,
  });

  return {
    question: data?.question,
    remainingTopics: data?.remainingTopics ?? [],
    isLoading,
    refetch,
    handleChangeTopic: setSelectedTopic,
  };
}
