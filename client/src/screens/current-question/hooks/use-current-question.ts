import { useContext, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext, SessionContext } from '@/providers';
import { RemainingTopic } from '@/generated/api';

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

  return remainingTopics?.[0]?.topic ?? lastAnsweredTopic;
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

      const topicToQuery = getTopicToQuery(selectedTopic, lastAnsweredTopic, remainingTopics);

      const question = await questionsApi.getNextQuestion({ topic: topicToQuery });
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
